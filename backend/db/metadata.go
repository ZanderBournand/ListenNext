package db

import (
	"time"
)

func UpdateLastScrapeTime(updateTime time.Time) error {
	updateTimeString := updateTime.Format("2006-01-02 15:04:05")

	query := `UPDATE metadata SET value=$1 WHERE property=$2`

	_, err := db.Exec(query, updateTimeString, "last_scrape_time")
	if err != nil {
		return err
	}

	return nil
}

func GetLastScrapeTime() *time.Time {
	var lastScrapeTime string
	backupTime := time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)

	err := db.QueryRow("SELECT value FROM metadata WHERE property = 'last_scrape_time'").Scan(&lastScrapeTime)
	if err != nil {
		return &backupTime
	}

	timestamp, err := time.Parse("2006-01-02 15:04:05", lastScrapeTime)
	if err != nil {
		return &backupTime
	}
	timestamp = timestamp.UTC()

	return &timestamp
}
