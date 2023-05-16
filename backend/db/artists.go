package db

import (
	"database/sql"
	"fmt"
	"main/graph/model"
	"main/types"
	"strings"
	"sync"
	"time"

	"github.com/lib/pq"
)

func GetMatchingArtists(artists []types.SpotifyArtist, genres []string) []int {
	var ids []int

	var allIDs []string
	var allGenres []string

	for _, artist := range artists {
		allIDs = append(allIDs, "'"+artist.ID+"'")
	}

	for _, genre := range genres {
		compareType := strings.ToLower(strings.ReplaceAll(genre, " ", ""))
		allGenres = append(allGenres, "'"+compareType+"'")
	}

	query := `
	    SELECT id
	    FROM Artists
	    WHERE spotify_id IN (` + strings.Join(allIDs, ",") + `)
	        OR EXISTS (
	            SELECT 1
	            FROM Artists_Genres
	            JOIN Genres ON Artists_Genres.genre_id = Genres.id
	            WHERE Artists_Genres.artist_id = Artists.id
	                AND Genres.compare_type IN (` + strings.Join(allGenres, ",") + `)
	        )
	`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			fmt.Println(err)
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		fmt.Println(err)
	}

	return ids
}

func UploadArtist(artist string, spotifyArtist *types.SpotifyArtist) (int64, int, error) {
	compareName := strings.ToLower(strings.ReplaceAll(artist, " ", ""))

	if spotifyArtist != nil {
		var id int64
		err := db.QueryRow("SELECT id FROM Artists WHERE name=$1 AND spotify_id=$2", spotifyArtist.Name, spotifyArtist.ID).Scan(&id)
		if err == nil {
			_, err = db.Exec("UPDATE Artists SET popularity=$1 WHERE name=$2 AND spotify_id=$3",
				spotifyArtist.Popularity, spotifyArtist.Name, spotifyArtist.ID)
			if err == nil {
				AddOrUpdateArtistGenres(id, spotifyArtist.Genres)
				return id, spotifyArtist.Popularity, nil
			}
			return -1, -1, fmt.Errorf("error updating artist w/ spotify")
		} else if err == sql.ErrNoRows {
			var newID int64
			err = db.QueryRow("INSERT INTO Artists(name, spotify_id, popularity, compare_name) VALUES($1, $2, $3, $4) RETURNING id",
				spotifyArtist.Name, spotifyArtist.ID, spotifyArtist.Popularity, compareName).Scan(&newID)
			if err == nil {
				AddOrUpdateArtistGenres(newID, spotifyArtist.Genres)
				return newID, spotifyArtist.Popularity, nil
			}
			return -1, -1, fmt.Errorf("error inserting artist w/ spotify")
		} else {
			return -1, -1, fmt.Errorf("error querying artist w/ spotify")
		}
	} else {
		compareName := strings.ToLower(strings.ReplaceAll(artist, " ", ""))

		var id int64
		err := db.QueryRow("SELECT id FROM Artists WHERE compare_name=$1", compareName).Scan(&id)
		if err == nil {
			return id, -1, nil
		} else if err == sql.ErrNoRows {
			var newID int64
			err = db.QueryRow("INSERT INTO Artists(name, compare_name) VALUES($1, $2) RETURNING id",
				artist, compareName).Scan(&newID)
			if err == nil {
				return newID, -1, nil
			}
			return -1, -1, fmt.Errorf("error inserting artist")
		} else {
			return -1, -1, fmt.Errorf("error querying artist")
		}
	}
}

func AddOrUpdateArtistGenres(artistId int64, genres []string) {
	for _, genre := range genres {
		compareType := strings.ToLower(strings.ReplaceAll(genre, " ", ""))

		var id int64
		err := db.QueryRow("SELECT id FROM Genres WHERE compare_type=$1", compareType).Scan(&id)
		if err == sql.ErrNoRows {
			_ = db.QueryRow("INSERT INTO Genres(type, compare_type) VALUES($1, $2) RETURNING id",
				genre, compareType).Scan(&id)
		}

		if id != 0 {
			_, err := db.Exec("INSERT INTO artists_genres (artist_id, genre_id) VALUES ($1, $2) ON CONFLICT DO NOTHING", artistId, id)
			if err != nil {
				fmt.Println(artistId)
				fmt.Println("Error inserting artist genre join")
			}
		}
	}
}

func FindArtistReleaseCount(artist *model.Artist) error {
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	recentStart := now
	recentEnd := now
	upcomingStart := now
	upcomingEnd := now

	recentStart = recentStart.AddDate(0, 0, -14)
	daysUntilFriday := (5 - int(recentStart.Weekday()) + 7) % 7
	recentStart = recentStart.AddDate(0, 0, daysUntilFriday+1)
	recentEnd = recentEnd.AddDate(0, 0, -7)
	daysUntilFriday = (5 - int(recentEnd.Weekday()) + 7) % 7
	recentEnd = recentEnd.AddDate(0, 0, daysUntilFriday)

	upcomingStart = upcomingStart.AddDate(0, 0, -7)
	daysUntilFriday = (5 - int(upcomingStart.Weekday()) + 7) % 7
	upcomingStart = upcomingStart.AddDate(0, 0, daysUntilFriday+1)
	upcomingEnd = upcomingEnd.AddDate(0, 0, 84)
	daysUntilFriday = (5 - int(upcomingEnd.Weekday()) + 7) % 7
	upcomingEnd = upcomingEnd.AddDate(0, 0, daysUntilFriday)

	query := `
		SELECT
			COUNT(*) FILTER (WHERE r.date >= $2 AND r.date <= $3) AS recent_release_count,
			COUNT(*) FILTER (WHERE r.date >= $4 AND r.date <= $5) AS upcoming_release_count
		FROM releases AS r
		JOIN releases_artists AS ra ON r.id = ra.release_id
		JOIN artists AS a ON ra.artist_id = a.id
		WHERE a.spotify_id = $1
	`

	var recentReleaseCount, upcomingReleaseCount int
	err := db.QueryRow(query, artist.SpotifyID, recentStart, recentEnd, upcomingStart, upcomingEnd).Scan(&recentReleaseCount, &upcomingReleaseCount)
	if err != nil {
		return err
	}

	artist.RecentReleasesCount = recentReleaseCount
	artist.UpcomingReleasesCount = upcomingReleaseCount

	return nil
}

func FindArtistRecentReleases(artist *model.Artist) error {
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	recentStart := now
	recentEnd := now

	recentStart = recentStart.AddDate(0, 0, -14)
	daysUntilFriday := (5 - int(recentStart.Weekday()) + 7) % 7
	recentStart = recentStart.AddDate(0, 0, daysUntilFriday+1)
	recentEnd = recentEnd.AddDate(0, 0, -7)
	daysUntilFriday = (5 - int(recentEnd.Weekday()) + 7) % 7
	recentEnd = recentEnd.AddDate(0, 0, daysUntilFriday)

	query := `
		SELECT r.id, r.title, array_agg(DISTINCT COALESCE(main.name, '')) AS artists, array_agg(DISTINCT COALESCE(feature.name, '')) AS featurings, array_agg(DISTINCT COALESCE(main.spotify_id, '')) AS artists_spotify_ids, array_agg(DISTINCT COALESCE(feature.spotify_id, '')) AS featurings_spotify_ids, r.date, r.cover, array_agg(DISTINCT COALESCE(g.type, '')) AS genres, array_agg(DISTINCT COALESCE(p.name, '')) AS producers,
		r.tracklist, r.type, r.aoty_id, r.trending_score, $1 as relationship
		FROM releases AS r
		JOIN releases_artists AS ra ON r.id = ra.release_id AND ra.relationship = 'main'
		JOIN artists AS main ON ra.artist_id = main.id
		LEFT JOIN releases_artists AS raf ON r.id = raf.release_id AND raf.relationship = 'feature'
		LEFT JOIN artists AS feature ON raf.artist_id = feature.id
		LEFT JOIN releases_producers AS rp ON r.id = rp.release_id
		LEFT JOIN producers AS p ON rp.producer_id = p.id
		LEFT JOIN releases_genres AS rg ON r.id = rg.release_id
		LEFT JOIN genres AS g ON rg.genre_id = g.id
	`

	artistRoles := []string{"main", "feature"}

	var releases []*model.Release

	errCh := make(chan error, len(artistRoles))
	var wg sync.WaitGroup
	wg.Add(len(artistRoles))

	for _, artistRole := range artistRoles {
		go func(artistRole string) {
			defer wg.Done()

			finalQuery := query

			finalQuery += fmt.Sprintf(" WHERE %s.spotify_id = '%s'", artistRole, artist.SpotifyID)
			finalQuery += fmt.Sprintf(" AND r.date >= '%s' AND r.date <= '%s'", recentStart.Format("2006-01-02"), recentEnd.Format("2006-01-02"))
			finalQuery += "GROUP BY r.id"

			rows, err := db.Query(finalQuery, artistRole)
			if err != nil {
				errCh <- err
				return
			}
			defer rows.Close()

			for rows.Next() {
				var release model.Release
				var artists, featurings, artistsSpotifyIDs, featuringsSpotifyIDs []string
				err := rows.Scan(
					&release.ID,
					&release.Title,
					pq.Array(&artists),
					pq.Array(&featurings),
					pq.Array(&artistsSpotifyIDs),
					pq.Array(&featuringsSpotifyIDs),
					&release.ReleaseDate,
					&release.Cover,
					pq.Array(&release.Genres),
					pq.Array(&release.Producers),
					pq.Array(&release.Tracklist),
					&release.Type,
					&release.AotyID,
					&release.TrendingScore,
					&release.ArtistRole,
				)
				if err != nil {
					errCh <- err
					return
				}

				var releaseArtists []*model.Artist
				var releaseFeaturings []*model.Artist

				for i, artist := range artists {
					var releaseArtist model.Artist

					releaseArtist.Name = artist
					releaseArtist.SpotifyID = artistsSpotifyIDs[i]

					releaseArtists = append(releaseArtists, &releaseArtist)
				}

				for i, featuring := range featurings {
					var releaseFeaturing model.Artist

					releaseFeaturing.Name = featuring
					releaseFeaturing.SpotifyID = featuringsSpotifyIDs[i]

					releaseFeaturings = append(releaseFeaturings, &releaseFeaturing)
				}

				release.Artists = releaseArtists
				release.Featurings = releaseFeaturings

				releases = append(releases, &release)
			}

			if err = rows.Err(); err != nil {
				errCh <- err
				return
			}
		}(artistRole)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	artist.RecentReleases = releases

	return nil
}

func FindArtistUpcomingReleases(artist *model.Artist) error {
	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	upcomingStart := now
	upcomingEnd := now

	upcomingStart = upcomingStart.AddDate(0, 0, -7)
	daysUntilFriday := (5 - int(upcomingStart.Weekday()) + 7) % 7
	upcomingStart = upcomingStart.AddDate(0, 0, daysUntilFriday+1)
	upcomingEnd = upcomingEnd.AddDate(0, 0, 84)
	daysUntilFriday = (5 - int(upcomingEnd.Weekday()) + 7) % 7
	upcomingEnd = upcomingEnd.AddDate(0, 0, daysUntilFriday)

	query := `
		SELECT r.id, r.title, array_agg(DISTINCT main.name) AS artists, array_agg(DISTINCT feature.name) AS featurings, array_agg(DISTINCT main.spotify_id) AS artists_spotify_ids, array_agg(DISTINCT feature.spotify_id) AS featurings_spotify_ids, r.date, r.cover, array_agg(DISTINCT g.type) AS genres, array_agg(DISTINCT p.name) AS producers, r.tracklist, r.type, r.aoty_id, r.trending_score, $1 as relationship
		FROM releases AS r
		JOIN releases_artists AS ra ON r.id = ra.release_id AND ra.relationship = 'main'
		JOIN artists AS main ON ra.artist_id = main.id
		LEFT JOIN releases_artists AS raf ON r.id = raf.release_id AND raf.relationship = 'feature'
		LEFT JOIN artists AS feature ON raf.artist_id = feature.id
		LEFT JOIN releases_producers AS rp ON r.id = rp.release_id
		LEFT JOIN producers AS p ON rp.producer_id = p.id
		LEFT JOIN releases_genres AS rg ON r.id = rg.release_id
		LEFT JOIN genres AS g ON rg.genre_id = g.id
	`

	artistRoles := []string{"main", "feature"}

	var releases []*model.Release

	errCh := make(chan error, len(artistRoles))
	var wg sync.WaitGroup
	wg.Add(len(artistRoles))

	for _, artistRole := range artistRoles {
		go func(artistRole string) {
			defer wg.Done()

			finalQuery := query

			finalQuery += fmt.Sprintf(" WHERE %s.spotify_id = '%s'", artistRole, artist.SpotifyID)
			finalQuery += fmt.Sprintf(" AND r.date >= '%s' AND r.date <= '%s'", upcomingStart.Format("2006-01-02"), upcomingEnd.Format("2006-01-02"))
			finalQuery += "GROUP BY r.id"

			rows, err := db.Query(finalQuery, artistRole)
			if err != nil {
				errCh <- err
				return
			}
			defer rows.Close()

			for rows.Next() {
				var release model.Release
				var cover, aoty_id sql.NullString
				var trending_score sql.NullFloat64
				var artists, featurings, artistsSpotifyIDs, featuringsSpotifyIDs, genres, producers, tracklist []sql.NullString
				err := rows.Scan(
					&release.ID,
					&release.Title,
					pq.Array(&artists),
					pq.Array(&featurings),
					pq.Array(&artistsSpotifyIDs),
					pq.Array(&featuringsSpotifyIDs),
					&release.ReleaseDate,
					&cover,
					pq.Array(&genres),
					pq.Array(&producers),
					pq.Array(&tracklist),
					&release.Type,
					&aoty_id,
					&trending_score,
					&release.ArtistRole,
				)
				if err != nil {
					errCh <- err
					return
				}

				if cover.Valid {
					release.Cover = &cover.String
				}
				if aoty_id.Valid {
					release.AotyID = &aoty_id.String
				}
				if trending_score.Valid {
					release.TrendingScore = &trending_score.Float64
				}

				var releaseArtists []*model.Artist
				var releaseFeaturings []*model.Artist
				var releaseGenres []string
				var releaseProducers []string
				var releaseTracklist []string

				for i, artist := range artists {
					if artist.Valid {
						var releaseArtist model.Artist
						releaseArtist.Name = artist.String
						releaseArtist.SpotifyID = artistsSpotifyIDs[i].String
						releaseArtists = append(releaseArtists, &releaseArtist)
					}
				}
				for i, featuring := range featurings {
					if featuring.Valid {
						var releaseFeaturing model.Artist
						releaseFeaturing.Name = featuring.String
						releaseFeaturing.SpotifyID = featuringsSpotifyIDs[i].String
						releaseFeaturings = append(releaseFeaturings, &releaseFeaturing)
					}
				}
				for _, genre := range genres {
					if genre.Valid {
						releaseGenres = append(releaseGenres, genre.String)
					}
				}
				for _, producer := range producers {
					if producer.Valid {
						releaseProducers = append(releaseProducers, producer.String)
					}
				}
				for _, track := range tracklist {
					if track.Valid {
						releaseTracklist = append(releaseTracklist, track.String)
					}
				}

				release.Artists = releaseArtists
				release.Featurings = releaseFeaturings
				release.Genres = releaseGenres
				release.Producers = releaseProducers
				release.Tracklist = releaseTracklist

				releases = append(releases, &release)
			}

			if err = rows.Err(); err != nil {
				errCh <- err
				return
			}
		}(artistRole)
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	artist.UpcomingReleases = releases

	return nil
}
