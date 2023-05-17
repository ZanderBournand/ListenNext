package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/config"
	"main/db"
	"main/graph/model"
	"main/types"
	"math/rand"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"
)

var (
	scrapingClientIDs = []string{
		os.Getenv("SPOTIFY_API_SCRAPING_CLIENT_IDS"),
		os.Getenv("SPOTIFY_API_SCRAPING_CLIENT_IDS"),
		os.Getenv("SPOTIFY_API_SCRAPING_CLIENT_IDS"),
	}
	scrapingClientSecrets = []string{
		os.Getenv("SPOTIFY_API_GENERAL_CLIENT_SECRET"),
		os.Getenv("SPOTIFY_API_GENERAL_CLIENT_SECRET"),
		os.Getenv("SPOTIFY_API_GENERAL_CLIENT_SECRET"),
	}
	generalClientID     = os.Getenv("SPOTIFY_API_GENERAL_CLIENT_ID")
	generalClientSecret = os.Getenv("SPOTIFY_API_GENERAL_CLIENT_SECRET")
)

var (
	scrapingTokens []string
	generalToken   string
)

func SpotifyScrapeTokens() {
	tokens := make([]string, 3)

	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader("grant_type=client_credentials"))
		if err != nil {
			panic(err)
		}

		req.SetBasicAuth(scrapingClientIDs[i], scrapingClientSecrets[i])
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()

		var data struct {
			AccessToken string `json:"access_token"`
		}
		if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
			panic(err)
		}

		tokens[i] = data.AccessToken
	}

	scrapingTokens = tokens
}

func SpotifyGeneralToken() {
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		panic(err)
	}

	req.SetBasicAuth(generalClientID, generalClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var data struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&data); err != nil {
		panic(err)
	}

	generalToken = data.AccessToken
}

func SpotifyUserTops(accessToken string) ([]types.SpotifyArtist, []types.SpotifyTrack, error) {
	var spotifyArtists []types.SpotifyArtist
	var spotifyTracks []types.SpotifyTrack

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := SpotifyUserTopArtists(accessToken, &spotifyArtists)
		if err != nil {
			log.Println("Error in SpotifyUserTopArtists:", err)
		}
	}()

	go func() {
		defer wg.Done()
		err := SpotifyUserTopTracks(accessToken, &spotifyArtists, &spotifyTracks)
		if err != nil {
			log.Println("Error in SpotifyUserTopTracks:", err)
		}
	}()

	wg.Wait()

	return spotifyArtists, spotifyTracks, nil
}

func SpotifyUserTopArtists(accessToken string, spotifyArtists *[]types.SpotifyArtist) error {
	limitQuery := "10"

	endpoint := fmt.Sprintf("https://api.spotify.com/v1/me/top/artists?limit=%s&time_range=medium_term", limitQuery)

	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("error making user tops request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making user tops request: %v", err)
	}
	defer res.Body.Close()

	var artistsData struct {
		Artists []types.SpotifyArtist `json:"items"`
	}

	err = json.NewDecoder(res.Body).Decode(&artistsData)
	if err != nil {
		return fmt.Errorf("error decoding artists response: %v", err)
	}

	*spotifyArtists = append(*spotifyArtists, artistsData.Artists...)

	return nil
}

func SpotifyUserTopTracks(accessToken string, spotifyArtists *[]types.SpotifyArtist, spotifyTracks *[]types.SpotifyTrack) error {
	limitQuery := "10"

	endpoint := fmt.Sprintf("https://api.spotify.com/v1/me/top/tracks?limit=%s&time_range=medium_term", limitQuery)

	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return fmt.Errorf("error making user tops request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making user tops request: %v", err)
	}
	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return fmt.Errorf("error decoding response: %v", err)
	}

	items := data["items"].([]interface{})

	var artistIDs []string

	for _, item := range items {
		var spotifyTrack types.SpotifyTrack
		spotifyTrack.ID = item.(map[string]interface{})["id"].(string)
		*spotifyTracks = append(*spotifyTracks, spotifyTrack)

		artists := item.(map[string]interface{})["artists"].([]interface{})
		for _, artist := range artists {
			artistID := artist.(map[string]interface{})["id"].(string)
			artistIDs = append(artistIDs, artistID)
		}
	}

	artistsEndpoint := fmt.Sprintf("https://api.spotify.com/v1/artists?ids=%s", strings.Join(artistIDs, ","))

	req, err = http.NewRequest("GET", artistsEndpoint, nil)
	if err != nil {
		return fmt.Errorf("error making artists request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("error making artists request: %v", err)
	}
	defer res.Body.Close()

	var artistsData struct {
		Artists []types.SpotifyArtist `json:"artists"`
	}

	err = json.NewDecoder(res.Body).Decode(&artistsData)
	if err != nil {
		return fmt.Errorf("error decoding artists response: %v", err)
	}

	*spotifyArtists = append(*spotifyArtists, artistsData.Artists...)

	return nil
}

func SpotifyRelatedArtists(artists []types.SpotifyArtist) ([]types.SpotifyArtist, error) {
	var wg sync.WaitGroup

	for _, artist := range artists {
		wg.Add(1)
		go func(artist types.SpotifyArtist) {
			defer wg.Done()
			endpoint := fmt.Sprintf("https://api.spotify.com/v1/artists/%s/related-artists", artist.ID)

			rand.Seed(time.Now().UnixNano())
			tokenIndex := rand.Intn(len(scrapingTokens))

			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				panic(err)
			}
			req.Header.Set("Authorization", "Bearer "+scrapingTokens[tokenIndex])

			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil || res.StatusCode != 200 {
				panic(err)
			}
			defer res.Body.Close()

			var data map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&data)
			if err != nil {
				panic(err)
			}

			relatedArtistsData := data["artists"].([]interface{})

			for _, relatedArtistData := range relatedArtistsData {
				relatedArtist := types.SpotifyArtist{
					Name: relatedArtistData.(map[string]interface{})["name"].(string),
					ID:   relatedArtistData.(map[string]interface{})["id"].(string),
				}
				genres := make([]string, 0)
				for _, genre := range relatedArtistData.(map[string]interface{})["genres"].([]interface{}) {
					genres = append(genres, genre.(string))
				}
				relatedArtist.Genres = genres

				artists = append(artists, relatedArtist)
			}
		}(artist)
	}

	wg.Wait()

	return artists, nil
}

func SpotifyRecommendations(artists []types.SpotifyArtist, tracks []types.SpotifyTrack) ([]types.SpotifyArtist, error) {
	posArtists := 0
	posTracks := 0

	artistsCount := len(artists)
	tracksCount := len(tracks)

	var wg sync.WaitGroup

	for posArtists < artistsCount || posTracks < tracksCount {

		var artistIDs []string
		if posArtists < artistsCount {
			// numArtists := rand.Intn(5) + 1
			// if posTracks == tracksCount {
			// 	numArtists = 5
			// }
			numArtists := 5
			for i := 0; i < numArtists && posArtists < artistsCount; i++ {
				artistIDs = append(artistIDs, artists[posArtists].ID)
				posArtists++
			}
		}

		var trackIDs []string
		if posTracks < tracksCount {
			numTracks := 5 - len(artistIDs)
			for i := 0; i < numTracks && posTracks < tracksCount; i++ {
				trackIDs = append(trackIDs, tracks[posTracks].ID)
				posTracks++
			}
		}

		wg.Add(1)

		go func() {
			defer wg.Done()

			recommendations, err := Recommendations(artistIDs, nil, trackIDs)
			if err != nil {
				log.Printf("error getting recommendations: %v", err)
				return
			}
			artists = append(artists, recommendations...)
		}()
	}

	wg.Wait()

	return artists, nil
}

func Recommendations(artistIds []string, genres []string, trackIds []string) ([]types.SpotifyArtist, error) {
	limitQuery := "10"

	params := url.Values{}
	if len(artistIds) > 0 {
		params.Set("seed_artists", strings.Join(artistIds, ","))
	}
	if len(genres) > 0 {
		params.Set("seed_genres", strings.Join(genres, ","))
	}
	if len(trackIds) > 0 {
		params.Set("seed_tracks", strings.Join(trackIds, ","))
	}

	endpoint := "https://api.spotify.com/v1/recommendations?limit=" + limitQuery + "&" + params.Encode()

	rand.Seed(time.Now().UnixNano())
	tokenIndex := rand.Intn(len(scrapingTokens))

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+scrapingTokens[tokenIndex])

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var artistIDs []string
	var recommendations struct {
		Tracks []struct {
			Artists []struct {
				ID string `json:"id"`
			} `json:"artists"`
		} `json:"tracks"`
	}

	err = json.NewDecoder(res.Body).Decode(&recommendations)
	if err != nil {
		return nil, err
	}

	for _, track := range recommendations.Tracks {
		for _, artist := range track.Artists {
			artistIDs = append(artistIDs, artist.ID)
		}
	}

	endpoint = fmt.Sprintf("https://api.spotify.com/v1/artists?ids=%s", strings.Join(artistIDs, ","))

	req, err = http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+scrapingTokens[tokenIndex])

	res, err = client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var artistsData struct {
		Artists []types.SpotifyArtist `json:"artists"`
	}

	err = json.NewDecoder(res.Body).Decode(&artistsData)
	if err != nil {
		return nil, err
	}

	var releaseArtists []types.SpotifyArtist
	releaseArtists = append(releaseArtists, artistsData.Artists...)

	return releaseArtists, nil
}

func TopGenres(artists []types.SpotifyArtist) []string {
	totalGenres := 20

	genreCount := make(map[string]int)
	for _, artist := range artists {
		for _, genre := range artist.Genres {
			genreCount[genre]++
		}
	}

	sortedGenres := make([]string, 0, len(genreCount))
	for genre := range genreCount {
		sortedGenres = append(sortedGenres, genre)
	}
	sort.Slice(sortedGenres, func(i, j int) bool {
		return genreCount[sortedGenres[i]] > genreCount[sortedGenres[j]]
	})

	topGenres := make([]string, 0, totalGenres)
	for i := 0; i < totalGenres && i < len(sortedGenres); i++ {
		topGenres = append(topGenres, sortedGenres[i])
	}

	return topGenres
}

func SpotifySearch(artist string) (*types.SpotifyArtist, error) {
	compareName := strings.ToLower(strings.ReplaceAll(artist, " ", ""))

	spotifyURL := fmt.Sprintf("https://api.spotify.com/v1/search?type=artist&q=%s", url.QueryEscape(artist))

	rand.Seed(time.Now().UnixNano())
	tokenIndex := rand.Intn(len(scrapingTokens))

	req, err := http.NewRequest("GET", spotifyURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+scrapingTokens[tokenIndex])

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != 200 {
		fmt.Println("SPOTIFY ERROR!!!!")
		return nil, err
	}
	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	artists := data["artists"].(map[string]interface{})["items"].([]interface{})
	for _, artist := range artists {
		name := strings.ToLower(strings.ReplaceAll(artist.(map[string]interface{})["name"].(string), " ", ""))
		if name == compareName {
			genres := make([]string, 0)
			for _, genre := range artist.(map[string]interface{})["genres"].([]interface{}) {
				genres = append(genres, genre.(string))
			}

			spotifyAritst := types.SpotifyArtist{
				Name:       artist.(map[string]interface{})["name"].(string),
				ID:         artist.(map[string]interface{})["id"].(string),
				Genres:     genres,
				Popularity: int(artist.(map[string]interface{})["popularity"].(float64)),
			}

			return &spotifyAritst, nil
		}
	}

	return nil, fmt.Errorf("no matching artist found")
}

func SpotifyUserInfo(accessToken string) (string, string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me", nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}

	var response struct {
		Email       string `json:"email"`
		DisplayName string `json:"display_name"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", "", err
	}

	return response.Email, response.DisplayName, nil
}

func SpotifyUserToken(userID string) string {
	accessToken, refreshToken, tokenExpiration := db.GetSpotifyUserTokens(userID)

	utcNow := time.Now().UTC()
	utcExpirationTime := tokenExpiration.UTC()

	if utcNow.After(utcExpirationTime) {
		newAccessToken, ExpirationIn, err := SpotifyRefreshToken(refreshToken)
		if err != nil {
			fmt.Println(err)
			return ""
		}

		location, err := time.LoadLocation("UTC")
		if err != nil {
			fmt.Println(err)
			return ""
		}

		expirationTime := time.Now().Add(time.Second * time.Duration(ExpirationIn)).In(location)
		db.UpdateSpotifyUserTokens(userID, newAccessToken, refreshToken, expirationTime)

		accessToken = newAccessToken
	}

	return accessToken
}

func SpotifyRefreshToken(refreshToken string) (string, int, error) {
	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", config.SpotifyOAuth.ClientID, config.SpotifyOAuth.ClientSecret)))

	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", -1, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", authHeader))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", -1, err
	}
	defer res.Body.Close()

	var token struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		ExpirationIn int    `json:"expires_in"`
		Scope        string `json:"scope"`
	}

	if err := json.NewDecoder(res.Body).Decode(&token); err != nil {
		return "", -1, err
	}

	return token.AccessToken, token.ExpirationIn, nil
}

func SpotifySearchArtists(query string) ([]*model.Artist, error) {
	queryLimit := 10

	spotifyURL := fmt.Sprintf("https://api.spotify.com/v1/search?type=artist&q=%s&limit=%d", url.QueryEscape(query), queryLimit)

	req, err := http.NewRequest("GET", spotifyURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+generalToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != 200 {
		fmt.Println("SPOTIFY ERROR!!!!")
		return nil, err
	}
	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	var spotifyArtists []*model.Artist

	artists := data["artists"].(map[string]interface{})["items"].([]interface{})
	for _, artist := range artists {
		var spotifyArtist model.Artist

		spotifyArtist.Name = artist.(map[string]interface{})["name"].(string)

		spotifyId := artist.(map[string]interface{})["id"].(string)
		spotifyArtist.SpotifyID = &spotifyId

		images := artist.(map[string]interface{})["images"].([]interface{})
		if len(images) > 0 {
			url := images[0].(map[string]interface{})["url"].(string)
			spotifyArtist.Image = &url
		}

		spotifyArtists = append(spotifyArtists, &spotifyArtist)
	}

	return spotifyArtists, nil
}

func SpotifyArtist(spotifyId string) (*model.Artist, error) {
	spotifyURL := fmt.Sprintf("https://api.spotify.com/v1/artists/%s", spotifyId)

	req, err := http.NewRequest("GET", spotifyURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+generalToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != 200 {
		fmt.Println("SPOTIFY ERROR!!!!")
		return nil, err
	}
	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	var artist model.Artist

	artist.SpotifyID = &spotifyId

	artist.Name = data["name"].(string)

	images := data["images"].([]interface{})
	if len(images) > 0 {
		imageURL := images[0].(map[string]interface{})["url"].(string)
		artist.Image = &imageURL
	}

	genres := data["genres"].([]interface{})
	for _, genre := range genres {
		genreStr := genre.(string)
		artist.Genres = append(artist.Genres, genreStr)
	}

	return &artist, nil
}

func SpotifyArtistTopTracks(artist *model.Artist) error {
	endpoint := fmt.Sprintf("https://api.spotify.com/v1/artists/%s/top-tracks?country=US", *artist.SpotifyID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+generalToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != 200 {
		fmt.Println("SPOTIFY ERROR!!!!")
		return err
	}
	defer res.Body.Close()

	var data struct {
		Tracks []struct {
			Name  string `json:"name"`
			ID    string `json:"id"`
			Album struct {
				ReleaseDate string `json:"release_date"`
				Images      []struct {
					URL string `json:"url"`
				} `json:"images"`
			} `json:"album"`
			Artists []struct {
				Name      string `json:"name"`
				SpotifyId string `json:"id"`
			} `json:"artists"`
			Type string `json:"type"`
		} `json:"tracks"`
	}

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return err
	}

	var topTracks []*model.Release

	for _, track := range data.Tracks {
		releaseDate, err := time.Parse("2006-01-02", track.Album.ReleaseDate)
		if err != nil {
			return err
		}

		spotifyRelease := &model.Release{
			Title:       track.Name,
			SpotifyID:   &track.ID,
			ReleaseDate: releaseDate,
			Type:        track.Type,
			Cover:       &track.Album.Images[0].URL,
			Artists:     make([]*model.Artist, len(track.Artists)),
		}

		for i, artist := range track.Artists {
			spotifyId := artist.SpotifyId
			spotifyRelease.Artists[i] = &model.Artist{
				Name:      artist.Name,
				SpotifyID: &spotifyId,
			}
		}

		topTracks = append(topTracks, spotifyRelease)
	}

	artist.TopTracks = topTracks

	return nil
}

func SpotifyArtistSingles(artist *model.Artist) error {
	endpoint := fmt.Sprintf("https://api.spotify.com/v1/artists/%s/albums?album_type=single&limit=10", *artist.SpotifyID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+generalToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != 200 {
		fmt.Println("SPOTIFY ERROR!!!!")
		return err
	}
	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return err
	}

	var singles []*model.Release

	items := data["items"].([]interface{})
	for _, album := range items {
		albumData := album.(map[string]interface{})
		var spotifyRelease model.Release

		spotifyRelease.Title = albumData["name"].(string)

		spotifyID := albumData["id"].(string)
		spotifyRelease.SpotifyID = &spotifyID

		cover := albumData["images"].([]interface{})[0].(map[string]interface{})["url"].(string)
		spotifyRelease.Cover = &cover

		artistsData := albumData["artists"].([]interface{})
		artists := make([]*model.Artist, len(artistsData))
		for i, artistData := range artistsData {
			spotifyID := artistData.(map[string]interface{})["id"].(string)
			artist := &model.Artist{
				SpotifyID: &spotifyID,
				Name:      artistData.(map[string]interface{})["name"].(string),
			}
			artists[i] = artist
		}
		spotifyRelease.Artists = artists

		releaseDate, err := time.Parse("2006-01-02", albumData["release_date"].(string))
		if err != nil {
			return err
		}
		spotifyRelease.ReleaseDate = releaseDate

		spotifyRelease.Type = albumData["album_type"].(string)

		singles = append(singles, &spotifyRelease)
	}

	artist.Singles = singles

	return nil
}

func SpotifyArtistAlbums(artist *model.Artist) error {
	endpoint := fmt.Sprintf("https://api.spotify.com/v1/artists/%s/albums?album_type=album&limit=10", *artist.SpotifyID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+generalToken)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil || res.StatusCode != 200 {
		fmt.Println("SPOTIFY ERROR!!!!")
		return err
	}
	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return err
	}

	var albums []*model.Release

	items := data["items"].([]interface{})
	for _, album := range items {
		albumData := album.(map[string]interface{})
		var spotifyRelease model.Release

		spotifyRelease.Title = albumData["name"].(string)

		spotifyID := albumData["id"].(string)
		spotifyRelease.SpotifyID = &spotifyID

		cover := albumData["images"].([]interface{})[0].(map[string]interface{})["url"].(string)
		spotifyRelease.Cover = &cover

		artistsData := albumData["artists"].([]interface{})
		artists := make([]*model.Artist, len(artistsData))
		for i, artistData := range artistsData {
			spotifyID := artistData.(map[string]interface{})["id"].(string)
			artist := &model.Artist{
				SpotifyID: &spotifyID,
				Name:      artistData.(map[string]interface{})["name"].(string),
			}
			artists[i] = artist
		}
		spotifyRelease.Artists = artists

		releaseDate, err := time.Parse("2006-01-02", albumData["release_date"].(string))
		if err != nil {
			return err
		}
		spotifyRelease.ReleaseDate = releaseDate

		spotifyRelease.Type = albumData["album_type"].(string)

		albums = append(albums, &spotifyRelease)
	}

	artist.Albums = albums

	return nil
}
