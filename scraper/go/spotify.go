package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	clientID     = "47f23432a2a6440096f223baaca92173"
	clientSecret = "e852ef0f03754e68a99dbf20a2735e61"
)

func SpotifyToken() (string, error) {
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader("grant_type=client_credentials"))
	if err != nil {
		return "", fmt.Errorf("error creating auth request: %v", err)
	}

	req.SetBasicAuth(clientID, clientSecret)
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

	return data.AccessToken, nil
}
