package main

import (
	"main/api"
	"main/db"
	"main/services"
)

func main() {
	db.Connect()

	services.SpotifyScrapeTokens()

	releases := services.Releases()
	db.Upload(releases, "spotify")

	api.Setup()

	// displayReleases, prev, next := db.GetTrendings("album", "next", 0, "past")
}
