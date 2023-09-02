package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB(chosenDB string) error {
	var (
		DB_NAME     = os.Getenv("DB_NAME")
		DB_USER     = os.Getenv("DB_USER")
		DB_PASSWORD = os.Getenv("DB_PASSWORD")
		DB_HOST     = os.Getenv("DB_HOST")
		DB_PORT     = os.Getenv("DB_PORT")
	)

	if chosenDB == "PostgreSQL" {
		// connect to db
		var err error
		connStr := fmt.Sprintf(
			"postgres://%v:%v@%v:%v/%v?sslmode=disable",
			DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME,
		)

		DB, err = sql.Open("postgres", connStr)
		if err != nil {
			return err
		}

		if err = DB.Ping(); err != nil {
			return err
		}
	} else {
		log.Printf("%v is not registered as a database.", chosenDB)
	}

	return nil
}
