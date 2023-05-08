package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var (
	clientIDs = []string{
		os.Getenv("SPOTIFY_API_SCRAPING_CLIENT_IDS"),
		os.Getenv("SPOTIFY_API_SCRAPING_CLIENT_IDS"),
		os.Getenv("SPOTIFY_API_SCRAPING_CLIENT_IDS"),
	}
	clientSecrets = []string{
		os.Getenv("SPOTIFY_API_GENERAL_CLIENT_SECRET"),
		os.Getenv("SPOTIFY_API_GENERAL_CLIENT_SECRET"),
		os.Getenv("SPOTIFY_API_GENERAL_CLIENT_SECRET"),
	}
)

func SpotifyTokens() ([]string, error) {
	tokens := make([]string, 3)

	for i := 0; i < 3; i++ {
		req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader("grant_type=client_credentials"))
		if err != nil {
			return nil, fmt.Errorf("error creating auth request: %v", err)
		}

		req.SetBasicAuth(clientIDs[i], clientSecrets[i])
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
