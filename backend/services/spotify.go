package services

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"main/config"
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
	limitQuery := "10"

	var spotifyArtists []types.SpotifyArtist
	var spotifyTracks []types.SpotifyTrack

	endpoint := fmt.Sprintf("https://api.spotify.com/v1/me/top/artists?limit=%s&time_range=short_term", limitQuery)

	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error making user tops request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err := client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error making user tops request: %v", err)
	}
	defer res.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding response: %v", err)
	}

	artists := data["items"].([]interface{})

	for _, artist := range artists {
		var spotifyArtist types.SpotifyArtist

		spotifyArtist.Name = artist.(map[string]interface{})["name"].(string)
		spotifyArtist.ID = artist.(map[string]interface{})["id"].(string)
		genres := make([]string, 0)
		for _, genre := range artist.(map[string]interface{})["genres"].([]interface{}) {
			genres = append(genres, genre.(string))
		}
		spotifyArtist.Genres = genres

		spotifyArtists = append(spotifyArtists, spotifyArtist)
	}

	endpoint = fmt.Sprintf("https://api.spotify.com/v1/me/top/tracks?limit=%s&time_range=short_term", limitQuery)

	req, err = http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error making user tops request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	res, err = client.Do(req)
	if err != nil {
		return nil, nil, fmt.Errorf("error making user tops request: %v", err)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return nil, nil, fmt.Errorf("error decoding response: %v", err)
	}

	items := data["items"].([]interface{})

	var wg sync.WaitGroup

	for _, item := range items {
		wg.Add(1)
		go func(item interface{}) {
			defer wg.Done()

			var spotifyTrack types.SpotifyTrack
			spotifyTrack.ID = item.(map[string]interface{})["id"].(string)

			artists := item.(map[string]interface{})["artists"].([]interface{})

			for _, artist := range artists {

				var spotifyArtist types.SpotifyArtist

				spotifyArtist.Name = artist.(map[string]interface{})["name"].(string)
				spotifyArtist.ID = artist.(map[string]interface{})["id"].(string)

				endpoint := fmt.Sprintf("https://api.spotify.com/v1/artists/%s", spotifyArtist.ID)

				rand.Seed(time.Now().UnixNano())
				tokenIndex := rand.Intn(len(scrapingTokens))

				req, err := http.NewRequest("GET", endpoint, nil)
				if err != nil {
					log.Printf("error making artist request: %v", err)
					return
				}
				req.Header.Set("Authorization", "Bearer "+scrapingTokens[tokenIndex])

				client := &http.Client{}
				res, err := client.Do(req)
				if err != nil || res.StatusCode != 200 {
					log.Printf("error making artist request: %v", err)
					return
				}
				defer res.Body.Close()

				var artistData map[string]interface{}
				err = json.NewDecoder(res.Body).Decode(&artistData)
				if err != nil {
					log.Printf("error decoding artist response: %v", err)
					return
				}

				genres := artistData["genres"].([]interface{})
				for _, genre := range genres {
					spotifyArtist.Genres = append(spotifyArtist.Genres, genre.(string))
				}

				spotifyArtists = append(spotifyArtists, spotifyArtist)
			}

			spotifyTracks = append(spotifyTracks, spotifyTrack)
		}(item)
	}

	wg.Wait()

	return spotifyArtists, spotifyTracks, nil
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

			for _, artist := range recommendations {
				found := false
				for _, existingArtist := range artists {
					if existingArtist.ID == artist.ID {
						found = true
						break
					}
				}
				if !found {
					artists = append(artists, artist)
				}
			}
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

	var releaseArtists []types.SpotifyArtist
	var wg sync.WaitGroup

	for _, artistID := range artistIDs {
		wg.Add(1)
		go func(artistID string) {
			defer wg.Done()

			var spotifyArtist types.SpotifyArtist
			spotifyArtist.ID = artistID

			endpoint := fmt.Sprintf("https://api.spotify.com/v1/artists/%s", spotifyArtist.ID)

			rand.Seed(time.Now().UnixNano())
			tokenIndex := rand.Intn(len(scrapingTokens))

			req, err := http.NewRequest("GET", endpoint, nil)
			if err != nil {
				panic(err)
			}
			req.Header.Set("Authorization", "Bearer "+scrapingTokens[tokenIndex])

			client := &http.Client{}
			res, err := client.Do(req)
			if err != nil {
				panic(err)
			}
			defer res.Body.Close()

			var artistData map[string]interface{}
			err = json.NewDecoder(res.Body).Decode(&artistData)
			if err != nil {
				panic(err)
			}

			spotifyArtist.Name = artistData["name"].(string)

			genres := artistData["genres"].([]interface{})
			for _, genre := range genres {
				spotifyArtist.Genres = append(spotifyArtist.Genres, genre.(string))
			}

			releaseArtists = append(releaseArtists, spotifyArtist)
		}(artistID)
	}

	wg.Wait()

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
