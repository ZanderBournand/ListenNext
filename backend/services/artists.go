package services

import (
	"context"
	"fmt"
	"main/db"
	"main/graph/model"
	"sync"
)

func GetArtist(spotifyID string) *model.Artist {
	artist, _ := SpotifyArtist(spotifyID)

	var wg sync.WaitGroup
	wg.Add(5)

	go func(artist *model.Artist) {
		defer wg.Done()

		err := db.FindArtistRecentReleases(artist)
		if err != nil {
			fmt.Println(err)
		}
	}(artist)

	go func(artist *model.Artist) {
		defer wg.Done()

		err := db.FindArtistUpcomingReleases(artist)
		if err != nil {
			fmt.Println(err)
		}
	}(artist)

	go func(artist *model.Artist) {
		defer wg.Done()

		err := SpotifyArtistTopTracks(artist)
		if err != nil {
			fmt.Println("Error fetching artist top tracks")
		}
	}(artist)

	go func(artist *model.Artist) {
		defer wg.Done()

		err := SpotifyArtistSingles(artist)
		if err != nil {
			fmt.Println("Error fetching artist singles")
		}
	}(artist)

	go func(artist *model.Artist) {
		defer wg.Done()

		err := SpotifyArtistAlbums(artist)
		if err != nil {
			fmt.Println("Error fetching artist albums")
		}
	}(artist)

	wg.Wait()

	return artist
}

func SearchArtists(ctx context.Context, query string) *model.SearchArtists {
	searchResults, _ := SpotifySearchArtists(query)

	var wg sync.WaitGroup
	wg.Add(len(searchResults.Results) + len(searchResults.RelatedArtists))

	for _, artist := range searchResults.Results {
		go func(artist *model.Artist) {
			defer wg.Done()

			err := db.FindArtistReleaseCount(artist)
			if err != nil {
				fmt.Println("Error fetching artist recent/upcoming release count")
			}
		}(artist)
	}

	for _, artist := range searchResults.RelatedArtists {
		go func(artist *model.Artist) {
			defer wg.Done()

			err := db.FindArtistReleaseCount(artist)
			if err != nil {
				fmt.Println("Error fetching artist recent/upcoming release count")
			}
		}(artist)
	}

	wg.Wait()

	return searchResults
}
