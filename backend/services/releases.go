package services

import (
	"context"
	"fmt"
	"main/db"
	"main/graph/model"
	"main/middlewares"
	"sync"
)

func GetAllTrendingReleases(ctx context.Context, releaseType string) *model.AllReleasesList {
	var wg sync.WaitGroup
	wg.Add(4)

	var allReleases model.AllReleasesList

	go func() {
		defer wg.Done()

		pastReleases := db.GetTrendingReleases(releaseType, "desc", 0, "past")
		allReleases.Past = pastReleases.Releases
	}()
	go func() {
		defer wg.Done()

		weekReleases := db.GetTrendingReleases(releaseType, "desc", 0, "week")
		allReleases.Week = weekReleases.Releases
	}()
	go func() {
		defer wg.Done()

		monthReleases := db.GetTrendingReleases(releaseType, "desc", 0, "month")
		allReleases.Month = monthReleases.Releases
	}()
	go func() {
		defer wg.Done()

		extendedReleases := db.GetTrendingReleases(releaseType, "desc", 0, "extended")
		allReleases.Extended = extendedReleases.Releases
	}()

	wg.Wait()

	return &allReleases
}

func GetAllRecommendations(ctx context.Context) *model.AllRecommendations {
	userID := middlewares.CtxUserID(ctx)
	accessToken := SpotifyUserToken(userID)

	var wg sync.WaitGroup
	storedRecommendations := make(map[string][]*model.Release)
	periods := []string{"past", "week", "month", "extended"}

	for _, period := range periods {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()

			recommendations, _ := db.GetUserRecommendations(userID, p)
			storedRecommendations[p] = recommendations
		}(period)
	}

	fmt.Println("Fetching all recommendations...")

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

	newRecommendations := make(map[string][]*model.Release)
	for _, period := range periods {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()

			releases, _ := db.GetMatchingReleases(artistIds, genres, p)
			newRecommendations[p] = releases
		}(period)
	}

	for _, period := range periods {
		go func(p string) {
			db.UploadUserRecommendations(userID, newRecommendations[p])
		}(period)
	}

	wg.Wait()

	var allRecommendations model.AllRecommendations
	wg.Add(4)

	go func(newRecs, storedRecs []*model.Release) {
		defer wg.Done()
		allRecommendations.Past = CombineReleases(newRecs, storedRecs)
	}(newRecommendations["past"], storedRecommendations["past"])

	go func(newRecs, storedRecs []*model.Release) {
		defer wg.Done()
		allRecommendations.Week = CombineReleases(newRecs, storedRecs)
	}(newRecommendations["week"], storedRecommendations["week"])

	go func(newRecs, storedRecs []*model.Release) {
		defer wg.Done()
		allRecommendations.Month = CombineReleases(newRecs, storedRecs)
	}(newRecommendations["month"], storedRecommendations["month"])

	go func(newRecs, storedRecs []*model.Release) {
		defer wg.Done()
		allRecommendations.Extended = CombineReleases(newRecs, storedRecs)
	}(newRecommendations["extended"], storedRecommendations["extended"])

	wg.Wait()

	addRecommendationsTopArtists(&allRecommendations)

	return &allRecommendations
}

func addRecommendationsTopArtists(releases *model.AllRecommendations) {
	var topArtistsNames []string

	counter := 0
	for _, release := range releases.Past {
		for _, artist := range release.Artists {
			topArtistsNames = append(topArtistsNames, artist.Name)
		}

		counter++
		if counter == 5 {
			break
		}
	}

	for _, release := range releases.Week {
		for _, artist := range release.Artists {
			topArtistsNames = append(topArtistsNames, artist.Name)
		}

		counter++
		if counter == 10 {
			break
		}
	}

	for _, release := range releases.Month {
		for _, artist := range release.Artists {
			topArtistsNames = append(topArtistsNames, artist.Name)
		}

		counter++
		if counter == 15 {
			break
		}
	}

	for _, release := range releases.Extended {
		for _, artist := range release.Artists {
			topArtistsNames = append(topArtistsNames, artist.Name)
		}

		counter++
		if counter == 20 {
			break
		}
	}

	releases.Artists, _ = db.FindMatchingArtistNames(topArtistsNames)
}

func GetRecommendations(ctx context.Context, input model.RecommendationsInput) []*model.Release {
	userID := middlewares.CtxUserID(ctx)
	accessToken := SpotifyUserToken(userID)

	var wg sync.WaitGroup
	wg.Add(1)

	var storedRecommendations []*model.Release

	go func() {
		defer wg.Done()
		storedRecommendations, _ = db.GetUserRecommendations(userID, input.Period)
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
