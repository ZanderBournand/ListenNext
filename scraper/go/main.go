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
	tokens, err := SpotifyTokens()
	if err != nil {
		fmt.Println("Error getting Spotify access tokens:", err)
		return
	}

	// fetching new releases batch
	releases := Releases()

	// updating databsae with releases
	Upload(releases, "spotify", tokens, db)

	// fetching releases
	displayReleases, prev, next := getTrendings("album", "next", 0, "past", db)

	fmt.Println("PREV:", prev, "/ NEXT:", next)
	fmt.Println("----------------------------")
	fmt.Println("----------------------------")

	for _, r := range displayReleases {
		// fmt.Printf("ID: %d\n", r.ID)
		fmt.Printf("Title: %s\n", r.Title)
		fmt.Printf("Artists: %v\n", r.Artists)
		// fmt.Printf("Featurings: %v\n", r.Featurings)
		// fmt.Printf("Date: %s\n", r.Date)
		// fmt.Printf("Cover: %s\n", r.Cover)
		// fmt.Printf("Genres: %v\n", r.Genres)
		// fmt.Printf("Producers: %v\n", r.Producers)
		// fmt.Printf("Tracklist: %v\n", r.Tracklist)
		// fmt.Printf("Type: %s\n", r.Type)
		// fmt.Printf("AOTY ID: %v\n", r.AOTYID)
		fmt.Printf("Trending Score: %v\n", r.TrendingScore)
		fmt.Println("-----------------------------------")
	}
}
