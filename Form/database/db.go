package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var localDB *sql.DB

func Connect() *sql.DB {
	if localDB == nil {

		connStr := viper.GetString("db")

		db, err := sql.Open("postgres", connStr)
		if err != nil {
			panic(err)
		}
		if err := db.Ping(); err != nil {
			panic(err)
		}

		localDB = db
	}
	log.Println("Успешно")

	return localDB
}
