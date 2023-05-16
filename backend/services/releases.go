package services

import (
	"context"
	"fmt"
	"main/db"
	"main/graph/model"
	"main/middlewares"
	"time"
)

func GetRecommendations(ctx context.Context, input model.RecommendationsInput) []*model.Release {
	userID := middlewares.CtxUserID(ctx)
	accessToken, refreshToken, tokenExpiration := db.GetSpotifyUserTokens(userID)

	utcNow := time.Now().UTC()
	utcExpirationTime := tokenExpiration.UTC()

	if utcNow.After(utcExpirationTime) {
		newAccessToken, ExpirationIn, err := SpotifyRefreshToken(refreshToken)
		if err != nil {
			fmt.Println(err)
			return nil
		}

		location, err := time.LoadLocation("UTC")
		if err != nil {
			fmt.Println(err)
			return nil
		}

		expirationTime := time.Now().Add(time.Second * time.Duration(ExpirationIn)).In(location)
		db.UpdateSpotifyUserTokens(userID, newAccessToken, refreshToken, expirationTime)

		accessToken = newAccessToken
	}

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
