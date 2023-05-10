package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const connStr = os.Getenv("DB_CONNECT_STRING")

var db *sql.DB

func Connect() {
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
