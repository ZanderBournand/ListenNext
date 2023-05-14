package config

import "golang.org/x/oauth2"

var SpotifyOAuth = &oauth2.Config{
	ClientID:     os.Getenv("SPOTIFY_API_GENERAL_CLIENT_ID"),
	ClientSecret: os.Getenv("SPOTIFY_API_GENERAL_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8000/callback",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://accounts.spotify.com/authorize",
		TokenURL: "https://accounts.spotify.com/api/token",
	},
	Scopes: []string{"user-top-read", "user-read-private", "user-read-email"},
}
