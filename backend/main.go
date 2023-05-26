package main

import (
	"main/api"
	"main/db"
	"main/services"
)

func main() {
	// Connecting to postgres database
	db.Connect()

	// Schedule all cron jobs
	services.StartJobs()

	// Open API port for incoming requests
	api.Setup()
}
