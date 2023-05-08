package main

import (
	"fmt"
	"main/api"
	"main/db"
	"main/services"
)

func main() {
	_, dbErr := db.Connect()
	if dbErr != nil {
		fmt.Println(dbErr)
		return
	}

	_, SpErr := services.SpotifyScrapeTokens()
	if SpErr != nil {
		fmt.Println(SpErr)
		return
	}

	api.Setup()

	// releases := Releases()

	// Upload(releases, "spotify", tokens, db)

	// displayReleases, prev, next := getTrendings("album", "next", 0, "past", db)
}
