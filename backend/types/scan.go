package types

import (
	"main/graph/model"

	"gopkg.in/guregu/null.v4"
)

type ScanRelease struct {
	ID            null.Int
	Title         null.String
	Artists       []null.String
	Featurings    []null.String
	ArtistsIds    []null.String
	FeaturingsIds []null.String
	ReleaseDate   null.Time
	Cover         null.String
	Genres        []null.String
	Producers     []null.String
	Tracklist     []null.String
	Type          null.String
	AotyID        null.String
	SpotifyID     null.String
	TrendingScore null.Float
	ArtistRole    null.String
}

func ScanToRelease(scanRelease ScanRelease) model.Release {
	var release model.Release

	if scanRelease.ID.Valid {
		id := int(scanRelease.ID.Int64)
		release.ID = &id
	}

	if scanRelease.Title.Valid {
		release.Title = scanRelease.Title.String
	}

	if scanRelease.ReleaseDate.Valid {
		release.ReleaseDate = scanRelease.ReleaseDate.Time
	}

	if scanRelease.Cover.Valid {
		release.Cover = &scanRelease.Cover.String
	}

	if scanRelease.Type.Valid {
		release.Type = scanRelease.Type.String
	}

	if scanRelease.AotyID.Valid {
		release.AotyID = &scanRelease.AotyID.String
	}

	if scanRelease.SpotifyID.Valid {
		release.SpotifyID = &scanRelease.SpotifyID.String
	}

	if scanRelease.TrendingScore.Valid {
		trendingScore := scanRelease.TrendingScore.Float64
		release.TrendingScore = &trendingScore
	}

	if scanRelease.ArtistRole.Valid {
		release.ArtistRole = &scanRelease.ArtistRole.String
	}

	var artists []*model.Artist
	for _, artist := range scanRelease.Artists {
		if artist.Valid {
			var newArtist = &model.Artist{
				Name: artist.String,
			}
			artists = append(artists, newArtist)
		}
	}
	release.Artists = artists

	for i, artistId := range scanRelease.ArtistsIds {
		if artistId.Valid && i < len(release.Artists) {
			release.Artists[i].SpotifyID = &artistId.String
		}
	}

	var featurings []*model.Artist
	for _, featuring := range scanRelease.Featurings {
		if featuring.Valid {
			var newFeaturing = &model.Artist{
				Name: featuring.String,
			}
			featurings = append(featurings, newFeaturing)
		}
	}
	release.Featurings = featurings

	for i, featuringId := range scanRelease.FeaturingsIds {
		if featuringId.Valid && i < len(release.Featurings) {
			release.Featurings[i].SpotifyID = &featuringId.String
		}
	}

	var genres []string
	for _, genre := range scanRelease.Genres {
		if genre.Valid {
			genres = append(genres, genre.String)
		}
	}
	release.Genres = genres

	var producers []string
	for _, producer := range scanRelease.Producers {
		if producer.Valid {
			producers = append(producers, producer.String)
		}
	}
	release.Producers = producers

	var tracklist []string
	for _, track := range scanRelease.Tracklist {
		if track.Valid {
			tracklist = append(tracklist, track.String)
		}
	}
	release.Tracklist = tracklist

	return release
}
