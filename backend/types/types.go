package types

import "time"

type DisplayRelease struct {
	ID            int64
	Title         string
	Artists       []string
	Featurings    []string
	Date          time.Time
	Cover         string
	Genres        []string
	Producers     []string
	Tracklist     []string
	Type          string
	AOTYID        string
	TrendingScore float64
}

type Release struct {
	AOTY_Id    string    `json:"aoty_id"`
	Artists    []string  `json:"artists"`
	Featurings []string  `json:"featurings"`
	Title      string    `json:"title"`
	Date       time.Time `json:"date"`
	Cover      string    `json:"cover"`
	Genres     []string  `json:"genres"`
	Producers  []string  `json:"producers"`
	Tracklist  []string  `json:"tracklist"`
}

type SpotifyArtist struct {
	Name       string   `json:"name"`
	ID         string   `json:"id"`
	Genres     []string `json:"genres"`
	Popularity int      `json:"popularity"`
	Image      string   `json:"image"`
}

type SpotifyTrack struct {
	ID string `json:"id"`
}

type PopularityAverage struct {
	count int
	sum   int
}

func (pa *PopularityAverage) AddValue(value int) {
	pa.count++
	pa.sum += value
}

func (pa *PopularityAverage) GetAverage() float64 {
	if pa.count == 0 {
		return 0
	}
	return float64(pa.sum) / float64(pa.count)
}
