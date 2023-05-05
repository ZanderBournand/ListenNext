package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const connStr = "user=postgres password=4KQtqlpX5FhwJRpR host=db.owjswwnsfwjaatozenwo.supabase.co port=5432 dbname=postgres sslmode=require"

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
