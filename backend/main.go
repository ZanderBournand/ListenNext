package main

import (
	"main/api"
	"main/db"
	"main/services"
)

func main() {
	db.Connect()

	services.SpotifyGeneralToken()
	services.SpotifyScrapeTokens()

	services.ScrapeReleases()

	api.Setup()
}
