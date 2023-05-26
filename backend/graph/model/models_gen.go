// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"time"
)

type AllReleasesCount struct {
	Past     *ReleasesCount `json:"past"`
	Week     *ReleasesCount `json:"week"`
	Month    *ReleasesCount `json:"month"`
	Extended *ReleasesCount `json:"extended"`
}

type AllReleasesList struct {
	Past     []*Release `json:"past"`
	Week     []*Release `json:"week"`
	Month    []*Release `json:"month"`
	Extended []*Release `json:"extended"`
}

type Artist struct {
	SpotifyID             *string    `json:"spotify_id,omitempty"`
	Name                  string     `json:"name"`
	Image                 *string    `json:"image,omitempty"`
	Genres                []string   `json:"genres,omitempty"`
	Popularity            *int       `json:"popularity,omitempty"`
	RecentReleasesCount   *int       `json:"recent_releases_count,omitempty"`
	UpcomingReleasesCount *int       `json:"upcoming_releases_count,omitempty"`
	RecentReleases        []*Release `json:"recent_releases,omitempty"`
	UpcomingReleases      []*Release `json:"upcoming_releases,omitempty"`
	TopTracks             []*Release `json:"top_tracks,omitempty"`
	Singles               []*Release `json:"singles,omitempty"`
	Albums                []*Release `json:"albums,omitempty"`
}

type AuthOps struct {
	Login        interface{} `json:"login"`
	RefreshLogin interface{} `json:"refreshLogin"`
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
	ID            *int      `json:"_id,omitempty"`
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
	SpotifyID     *string   `json:"spotify_id,omitempty"`
	TrendingScore *float64  `json:"trending_score,omitempty"`
	ArtistRole    *string   `json:"artist_role,omitempty"`
}

type ReleasesCount struct {
	All     int `json:"all"`
	Albums  int `json:"albums"`
	Singles int `json:"singles"`
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

type SearchArtists struct {
	Results        []*Artist `json:"results"`
	RelatedArtists []*Artist `json:"related_artists"`
}
