package services

import (
	"context"
	"fmt"
	"main/db"
	"main/graph/model"
	"main/middlewares"
)

func GetRecommendations(ctx context.Context, input model.RecommendationsInput) []*model.Release {
	userID := middlewares.CtxUserID(ctx)
	accessToken := SpotifyUserToken(userID)

	fmt.Println("Fetching recommendations...")

	artists, tracks, err := SpotifyUserTops(accessToken)
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
	releases, _ := db.GetMatchingReleases(artistIds, genres, input.Period)

	return releases
}
