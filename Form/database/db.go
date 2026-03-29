package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	connStr := "user=postgres dbname=Demo sslmode=disable password=1234955 port=9090 host=localhost"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("ошибка подключения к БД: %v", err)
	}
	if err := db.Ping(); err != nil {
		return fmt.Errorf("не удалось подключиться к базе данных: %v", err)
	}
	DB = db
	log.Println("Успешно")
	return nil
}
