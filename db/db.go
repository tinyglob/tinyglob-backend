package db

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitDB() *sql.DB {
	connectionString := os.Getenv("DB_CONNECTION_URL")
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	return db
}
