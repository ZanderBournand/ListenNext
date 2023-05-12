package main

import (
	"main/api"
	"main/db"
)

func main() {
	db.Connect()

	// services.SpotifyScrapeTokens()

	// releases := services.Releases()
	// db.Upload(releases, "spotify")

	api.Setup()
}
