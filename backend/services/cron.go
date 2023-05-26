package services

import (
	"time"

	"github.com/go-co-op/gocron"
)

func StartJobs() {
	s := gocron.NewScheduler(time.UTC)

	s.Every(59).Minutes().Do(SpotifyGeneralToken)
	s.Every(59).Minutes().Do(SpotifyScrapeTokens)

	s.Every(6).Hours().Do(ScrapeReleases)

	s.StartAsync()
}
