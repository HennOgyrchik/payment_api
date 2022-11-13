package postgresql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type User struct {
	Id   int
	Cash int
}

func DbConnection() *sql.DB {
	connStr := "user=test password=123 dbname=postgres sslmode=disable port=5431"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		fmt.Printf("Ошибка подключения к БД:%s\n", err)
	}
	return db
}

func GetBalance(userId int) (cash string, err error) {
	db := DbConnection()
	defer db.Close()

	stmt, err := db.Prepare("select cash from users where id=$1")
	if err != nil {
		return
	}
	fmt.Println(cash)
	err = stmt.QueryRow(userId).Scan(&cash)
	if err != nil {
		return
	}
	return
}
