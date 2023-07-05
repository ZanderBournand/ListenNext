package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

func Connect() {
	connStr := os.Getenv("DB_CONNECT_STRING")
	database, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	err = database.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected to database!")

	db = database
}
