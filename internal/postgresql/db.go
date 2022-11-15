package postgresql

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
)

type User struct {
	Id   uint
	Cash string
}

func checkUser(id uint) (ok bool, err error) {
	db, err := DbConnection()
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("select id from users where id=$1")
	if err != nil {
		return
	}
	var temp int
	err = stmt.QueryRow(id).Scan(&temp)
	if err != nil {
		return
	}
	return true, nil
}

func DbConnection() (*sql.DB, error) {
	connStr := "user=test password=123 dbname=postgres sslmode=disable host=postgres port=5432" //ip = postgres
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetBalance(user *User) (err error) {
	ok, err := checkUser(user.Id)

	if err != nil {
		return err
	}
	if ok == false {
		return errors.New("User not found")
	}

	db, err := DbConnection()
	if err != nil {
		return
	}
	defer db.Close()

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

func Replenish(user *struct{ Id, Count uint }) (err error) {
	ok, err := checkUser(user.Id)

	if err != nil {
		return err
	}
	if ok == false {
		return errors.New("User not found")
	}
	db, err := DbConnection()
	if err != nil {
		return
	}
	defer db.Close()

	stmt, err := db.Prepare("update users set cash=cash+$1 where id=$2")
	if err != nil {
		return
	}

	_ = stmt.QueryRow(user.Count, user.Id)
	return err
}

func WriteTransaction(transaction struct{ UserID, ServiceID, OrderID, Cost uint }) error {
	ok, err := checkUser(transaction.UserID)
	if err != nil {
		return err
	}
	if ok == false {
		return err
	}

	db, err := DbConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("insert into transactions values (default,$1,$2,$3,$4, 'debet')returning id")
	if err != nil {
		return err
	}
	var val string
	err = stmt.QueryRow(transaction.UserID, transaction.ServiceID, transaction.OrderID, transaction.Cost).Scan(&val)
	if err != nil {
		return err
	}

	err = Replenish(&struct{ Id, Count uint }{1, transaction.Cost})
	if err != nil {
		_ = db.QueryRow("delete from transactions where id=" + val)
		return err
	}

	return err
}
