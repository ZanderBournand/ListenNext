package services

import (
	"fmt"
	"main/db"
	"main/models"
	"net/http"
)

func GetRecommendations(client *http.Client, period string) []models.DisplayRelease {

	fmt.Println("Fetching recommendations...")
	artists, tracks, err := SpotifyUserTops(client)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	artists, err = SpotifyRecommendations(artists, tracks)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	artists, err = SpotifyRelatedArtists(artists)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	genres := TopGenres(artists)

	artistIds := db.GetMatchingArtists(artists, genres)
	releases, _ := db.GetMatchingReleases(artistIds, genres, period)

	return releases
}
