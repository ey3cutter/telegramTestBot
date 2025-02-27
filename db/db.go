package db

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

var Db *sql.DB

func InitDb(connStr string) {
	var err error
	Db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Panic(err)
	}
	// Удалено defer Db.Close()

	err = Db.Ping()
	if err != nil {
		log.Panic(err)
	}
	fmt.Println("Подключились к БД")
}
