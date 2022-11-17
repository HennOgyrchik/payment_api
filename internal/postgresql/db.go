package postgresql

import (
	"database/sql"
	"encoding/csv"
	_ "github.com/lib/pq"
	"os"
	"strconv"
	"time"
	"turbo-carnival/internal/config"
)

type User struct {
	Id          uint `json:"user_id"`
	Cash, Count uint
	OrderID     uint `json:"order_id"`
	ServiceID   uint `json:"service_id"`
}

func (u *User) checkUser(db *sql.DB) (err error) {
	stmt, err := db.Prepare("select id from users where id=$1")
	if err != nil {
		return
	}
	var temp int
	err = stmt.QueryRow(u.Id).Scan(&temp)
	if err != nil {
		return
	}
	return
}

func DbConnection() (*sql.DB, error) {
	c := config.GetConfig()
	connStr := fmt.Sprintf("user=%s password=%s dbname=postgres sslmode=disable host=%s port=%s", c.PSQLLogin, c.PSQLPass, c.DBHost, c.PSQLPort)
	//fmt.Println(test)

	//connStr := "user=test password=123 dbname=postgres sslmode=disable host=postgres port=5432"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetBalance(user *User) (err error) {
	db, err := DbConnection()
	if err != nil {
		return
	}
	defer db.Close()

	err = user.checkUser(db)
	if err != nil {
		return
	}

	stmt, err := db.Prepare("select cash from users where id=$1")
	if err != nil {
		return
	}

	err = stmt.QueryRow(user.Id).Scan(&user.Cash)
	if err != nil {
		return
	}
	return
}

func Replenish(user *User) (err error) {
	db, err := DbConnection()
	if err != nil {
		return
	}
	defer db.Close()

	err = user.checkUser(db)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			stmt, err := db.Prepare("insert into users values ($1,default)")
			if err != nil {
				return err
			}
			_ = stmt.QueryRow(user.Id)
		} else {
			return err
		}
	}

	stmt, err := db.Prepare("insert into transactions (id, user_id,cost,type,date) values (default,$1,$2,'replenishment',default)returning id")
	if err != nil {
		return
	}
	var val string
	err = stmt.QueryRow(user.Id, user.Count).Scan(&val)
	if err != nil {
		return
	}
	return
}

func WriteTransaction(user *User) (err error) {
	db, err := DbConnection()
	if err != nil {
		return
	}
	defer db.Close()

	//проверка на наличие резервного счета (id= -1). Создание в случае его отсутствия.
	var val string
	err = db.QueryRow("select id from users where id=-1").Scan(&val)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			_ = db.QueryRow("insert into users values (-1,default)")
		} else {
			return
		}
	}

	err = user.checkUser(db)
	if err != nil {
		return
	}

	stmt, err := db.Prepare("insert into transactions values (default,$1,$2,$3,$4, 'buy',default)returning id")
	if err != nil {
		return
	}
	//var val string
	err = stmt.QueryRow(user.Id, user.ServiceID, user.OrderID, user.Count).Scan(&val)
	if err != nil {
		return
	}

	return
}

func RecognizeRevenue(user *User) (err error) {
	db, err := DbConnection()
	if err != nil {
		return
	}
	defer db.Close()

	err = user.checkUser(db)
	if err != nil {
		return
	}

	stmt, err := db.Prepare(" select id from transactions where user_id=$1 and service_id=$2 and order_id=$3 and cost=$4 and type='buy'")
	if err != nil {
		return
	}

	var result int
	err = stmt.QueryRow(user.Id, user.ServiceID, user.OrderID, user.Count).Scan(&result)
	if err != nil {
		return
	}

	stmt, err = db.Prepare("insert into transactions values (default,$1,$2,$3,$4, 'revenue',default)returning id")
	if err != nil {
		return
	}
	var val string
	err = stmt.QueryRow(user.Id, user.ServiceID, user.OrderID, user.Count).Scan(&val)
	if err != nil {
		return
	}

	return
}

func MonthlyReport(t time.Time) (err error) {
	db, err := DbConnection()
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("select service_id, SUM(cost) from transactions where type='revenue' AND (date >= $1 and date < (date($1) + INTERVAL '1 month')) group by service_id;")
	if err != nil {
		return
	}

	rows, err := stmt.Query(t)
	if err != nil {
		return
	}
	defer rows.Close()

	var arr [][]string
	arr = append(arr, []string{"Service_ID", "Sum"})

	for rows.Next() {
		var servID, sum int
		if err := rows.Scan(&servID, &sum); err != nil {
			return err
		}
		arr = append(arr, []string{strconv.Itoa(servID), strconv.Itoa(sum)})

	}
	//////////////

	f, err := os.Create("report.csv")
	defer f.Close()
	if err != nil {
		return
	}

	w := csv.NewWriter(f)
	w.Comma = ';'
	for _, record := range arr {
		if err := w.Write(record); err != nil {
			return err
		}
	}
	w.Flush()

	return err
}
