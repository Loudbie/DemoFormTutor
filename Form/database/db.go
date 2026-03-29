package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

var DB *sql.DB

func Connect() error {
	connStr := viper.GetString("db")
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
