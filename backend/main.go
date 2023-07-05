package main

import (
	"main/api"
	"main/config"
	"main/db"
	"main/services"
	"main/tools"
)

func main() {
	// Load .env variables
	tools.LoadEnvVariables()

	// Initialize config variables
	config.InitConfig()

	// Connecting to postgres database
	db.Connect()

	// Schedule all cron jobs
	services.StartJobs()

	// Open API port for incoming requests
	api.Setup()
}
