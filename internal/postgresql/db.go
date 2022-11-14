package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type User struct {
	Id   int
	Cash string
}

func DbConnection() *sql.DB {
	connStr := "user=test password=123 dbname=postgres sslmode=disable port=5431"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		fmt.Printf("Ошибка подключения к БД:%s\n", err)
	}
	return db
}

func GetBalance(user *User) (err error) {
	db := DbConnection()
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
	db := DbConnection()
	defer db.Close()

	stmt, err := db.Prepare("select id from users where id=$1")
	if err != nil {
		return
	}

	err = stmt.QueryRow(user.Id).Scan(&user.Id)
	if err != nil {
		return
	}

	stmt, err = db.Prepare("update users set cash=cash+$1 where id=$2")
	if err != nil {
		return
	}

	_ = stmt.QueryRow(user.Count, user.Id)
	return err
}
