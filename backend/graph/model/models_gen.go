// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type Artist struct {
	SpotifyID             string            `json:"spotify_id"`
	Name                  string            `json:"name"`
	Image                 string            `json:"image"`
	Genres                []*string         `json:"genres"`
	RecentReleasesCount   int               `json:"recent_releases_count"`
	UpcomingReleasesCount int               `json:"upcoming_releases_count"`
	RecentReleases        []*Release        `json:"recent_releases"`
	UpcomingReleases      []*Release        `json:"upcoming_releases"`
	TopTracks             []*SpotifyRelease `json:"top_tracks"`
	Singles               []*SpotifyRelease `json:"singles"`
	Albums                []*SpotifyRelease `json:"albums"`
}

type AuthOps struct {
	Login        interface{} `json:"login"`
	Register     interface{} `json:"register"`
	SpotifyLogin interface{} `json:"spotifyLogin"`
}

type NewUser struct {
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type RecommendationsInput struct {
	Period string `json:"period"`
}

type Release struct {
	ID            int       `json:"_id"`
	Title         string    `json:"title"`
	Artists       []*Artist `json:"artists"`
	Featurings    []*Artist `json:"featurings,omitempty"`
	ReleaseDate   time.Time `json:"release_date"`
	Cover         *string   `json:"cover,omitempty"`
	Genres        []string  `json:"genres,omitempty"`
	Producers     []string  `json:"producers,omitempty"`
	Tracklist     []string  `json:"tracklist,omitempty"`
	Type          string    `json:"type"`
	AotyID        *string   `json:"aoty_id,omitempty"`
	TrendingScore *float64  `json:"trending_score,omitempty"`
	ArtistRole    string    `json:"artist_role"`
}

type ReleasesInput struct {
	Type      string `json:"type"`
	Direction string `json:"direction"`
	Reference int    `json:"reference"`
	Period    string `json:"period"`
}

type ReleasesList struct {
	Releases []*Release `json:"releases"`
	Prev     bool       `json:"prev"`
	Next     bool       `json:"next"`
}

type SpotifyRelease struct {
	SpotifyReleaseID string    `json:"spotify_release_id"`
	Title            string    `json:"title"`
	Cover            string    `json:"cover"`
	Artists          []*Artist `json:"artists"`
	ReleaseDate      time.Time `json:"release_date"`
	Type             string    `json:"type"`
}
