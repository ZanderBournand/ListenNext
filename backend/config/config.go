package config

import (
	"os"

	"golang.org/x/oauth2"
)

var SpotifyOAuth = &oauth2.Config{
	ClientID:     "",
	ClientSecret: "",
	RedirectURL:  "http://localhost:3000/login",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://accounts.spotify.com/authorize",
		TokenURL: "https://accounts.spotify.com/api/token",
	},
	Scopes: []string{"user-top-read", "user-read-private", "user-read-email"},
}

func InitConfig() {
	// Initialize necessary config variables
	SpotifyOAuth.ClientID = os.Getenv("SPOTIFY_API_GENERAL_CLIENT_ID")
	SpotifyOAuth.ClientSecret = os.Getenv("SPOTIFY_API_GENERAL_CLIENT_SECRET")
}
