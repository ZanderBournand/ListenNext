package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
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
)

var (
	generalClientID     = os.Getenv("SPOTIFY_API_GENERAL_CLIENT_ID")
	generalClientSecret = os.Getenv("SPOTIFY_API_GENERAL_CLIENT_SECRET")
)

type SpotifyArtist struct {
	Name   string   `json:"name"`
	ID     string   `json:"id"`
	Genres []string `json:"genres"`
}

type SpotifyTrack struct {
	ID string `json:"id"`
}

func SpotifyScrapeTokens() ([]string, error) {
	tokens := make([]string, 3)

	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader("grant_type=client_credentials"))
		if err != nil {
			return nil, fmt.Errorf("error creating auth request: %v", err)
		}

		req.SetBasicAuth(scrapingClientIDs[i], scrapingClientSecrets[i])
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making auth request: %v", err)
		}
		defer resp.Body.Close()

		var data struct {
			AccessToken string `json:"access_token"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return nil, fmt.Errorf("error decoding auth response: %v", err)
		}

		tokens[i] = data.AccessToken
	}

	return tokens, nil
}

func SpotifyGeneralToken() (string, error) {

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return "", fmt.Errorf("error creating auth request: %v", err)
	}

	req.SetBasicAuth(generalClientID, generalClientSecret)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making auth request: %v", err)
	}
	defer resp.Body.Close()

	var data struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("error decoding auth response: %v", err)
	}

	token := data.AccessToken

	return token, nil

}

func SpotifyUserArtistsTops(client *http.Client) ([]SpotifyArtist, error) {
	endpoint := "https://api.spotify.com/v1/me/top/artists?limit=10&time_range=short_term"
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("error making user tops request: %v", err)
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	var spotifyArtists []SpotifyArtist

	artists := data["items"].([]interface{})

	for _, artist := range artists {
		var spotifyArtist SpotifyArtist

		spotifyArtist.Name = artist.(map[string]interface{})["name"].(string)
		spotifyArtist.ID = artist.(map[string]interface{})["id"].(string)
		genres := make([]string, 0)
		for _, genre := range artist.(map[string]interface{})["genres"].([]interface{}) {
			genres = append(genres, genre.(string))
		}
		spotifyArtist.Genres = genres

		spotifyArtists = append(spotifyArtists, spotifyArtist)
	}

	return spotifyArtists, nil
}

// func SpotifyUserTracksTops(client *http.Client) []SpotifyTrack {

// }
