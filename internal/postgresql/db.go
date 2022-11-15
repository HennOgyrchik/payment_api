package postgresql

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
)

type User struct {
	Id        uint `json:"user_id"`
	Cash      string
	Count     uint
	OrderID   uint `json:"order_id"`
	ServiceID uint `json:"service_id"`
}

func (u *User) checkUser() (ok bool, err error) {
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
	err = stmt.QueryRow(u.Id).Scan(&temp)
	if err != nil {

		return
	}
	return true, nil
}

func DbConnection() (*sql.DB, error) {
	connStr := "user=test password=123 dbname=postgres sslmode=disable host=postgres port=5432"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetBalance(user *User) (err error) {
	ok, err := user.checkUser()
	if (err != nil) || (ok == false) {
		return
	}

	db, err := DbConnection()
	if err != nil {
		return errors.New("Database connection error")
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

func Replenish(user *User) (err error) {
	ok, err := user.checkUser()
	if (err != nil) || (ok == false) {
		return err
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

func WriteTransaction(user *User) error {
	ok, err := user.checkUser()
	if (err != nil) || (ok == false) {
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
	err = stmt.QueryRow(user.Id, user.ServiceID, user.OrderID, user.Count).Scan(&val)
	if err != nil {
		return err
	}

	err = Replenish(user)
	if err != nil {
		_ = db.QueryRow("delete from transactions where id=" + val)
		return err
	}

	return err
}

func RecognizeRevenue(user *User) (err error) {
	ok, err := user.checkUser()
	if (err != nil) || (ok == false) {
		return
	}
	//списывает из резерва деньги, добавляет данные в отчет для бухгалтерии. Принимает id пользователя, ИД услуги, ИД заказа, сумму

	return
}
