package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const connStr = os.Getenv("DB_CONNECT_STRING")

func Connect() (*sql.DB, error) {
	// connecting to Supabase Postgre
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}

	fmt.Println("Successfully connected to database!")

	return db, nil
}
