package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const connStr = os.Getenv("DB_CONNECT_STRING")

func main() {
	// connecting to Supabase Postgre
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Supabase PostgreSQL database")
	fmt.Println("----------------------------")
	fmt.Println("----------------------------")

	// getting spotify authorization token
	token, err := SpotifyToken()
	if err != nil {
		fmt.Println("Error getting Spotify access token:", err)
		return
	}

	// fetching new releases batch
	releases := Releases()

	// updating databsae with releases
	Upload(releases, token, db)
}
