package services

import (
	"context"
	"fmt"
	"main/db"
	"main/graph/model"
	"main/middlewares"
	"sync"
)

func GetRecommendations(ctx context.Context, input model.RecommendationsInput) []*model.Release {
	userID := middlewares.CtxUserID(ctx)
	accessToken := SpotifyUserToken(userID)

	var wg sync.WaitGroup
	wg.Add(1)

	var storedRecommendations []*model.Release

	go func() {
		defer wg.Done()
		storedRecommendations, _ = db.GetUserRecommendations(userID)
	}()

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

	go func(releases []*model.Release) {
		db.UploadUserRecommendations(userID, releases)
	}(releases)

	wg.Wait()

	return CombineReleases(releases, storedRecommendations)
}

func CombineReleases(releases1, releases2 []*model.Release) []*model.Release {
	combinedReleases := make([]*model.Release, 0, len(releases1)+len(releases2))
	i, j := 0, 0

	for i < len(releases1) && j < len(releases2) {
		release1 := releases1[i]
		release2 := releases2[j]

		if *release1.ID == *release2.ID {
			combinedReleases = append(combinedReleases, release1)
			i++
			j++
		} else if int(*release1.TrendingScore) == int(*release2.TrendingScore) && *release1.ID > *release2.ID {
			combinedReleases = append(combinedReleases, release1)
			i++
		} else if int(*release1.TrendingScore) == int(*release2.TrendingScore) && *release1.ID < *release2.ID {
			combinedReleases = append(combinedReleases, release2)
			j++
		} else if int(*release1.TrendingScore) > int(*release2.TrendingScore) {
			combinedReleases = append(combinedReleases, release1)
			i++
		} else if int(*release1.TrendingScore) < int(*release2.TrendingScore) {
			combinedReleases = append(combinedReleases, release2)
			j++
		}
	}

	for ; i < len(releases1); i++ {
		combinedReleases = append(combinedReleases, releases1[i])
	}

	for ; j < len(releases2); j++ {
		combinedReleases = append(combinedReleases, releases2[j])
	}

	return combinedReleases
}
