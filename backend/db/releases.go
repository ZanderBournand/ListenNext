package db

import (
	"database/sql"
	"fmt"
	"main/graph/model"
	"main/tools"
	"main/types"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

func GetMatchingReleases(ids []int, genres []string, period string) ([]*model.Release, error) {
	var releases []*model.Release

	startDate, endDate := tools.GetReleaseDates(period)

	idStrings := make([]string, 0, len(ids))
	for _, id := range ids {
		idStrings = append(idStrings, strconv.Itoa(id))
	}

	var allGenres []string
	for _, genre := range genres {
		compareType := strings.ToLower(strings.ReplaceAll(genre, " ", ""))
		allGenres = append(allGenres, "'"+compareType+"'")
	}

	query := fmt.Sprintf(`
		SELECT r.id, r.title, array_agg(DISTINCT main.name) AS artists, array_agg(DISTINCT feature.name) AS featurings, array_agg(DISTINCT main.spotify_id) AS artists_spotify_ids, array_agg(DISTINCT feature.spotify_id) AS featurings_spotify_ids, r.date, r.cover, array_agg(DISTINCT g.type) AS genres, array_agg(DISTINCT p.name) AS producers, r.tracklist, r.type, r.aoty_id, r.trending_score
		FROM releases AS r
		JOIN releases_artists AS ra ON r.id = ra.release_id AND ra.relationship = 'main'
		JOIN artists AS main ON ra.artist_id = main.id
		LEFT JOIN releases_artists AS raf ON r.id = raf.release_id AND raf.relationship = 'feature'
		LEFT JOIN artists AS feature ON raf.artist_id = feature.id
		LEFT JOIN releases_producers AS rp ON r.id = rp.release_id
		LEFT JOIN producers AS p ON rp.producer_id = p.id
		LEFT JOIN releases_genres AS rg ON r.id = rg.release_id
		LEFT JOIN genres AS g ON rg.genre_id = g.id
		WHERE r.trending_score IS NOT NULL AND (main.id IN (%s) OR g.compare_type IN (%s))
	`, strings.Join(idStrings, ","), strings.Join(allGenres, ","))

	query += fmt.Sprintf(" AND r.date >= '%s' AND r.date <= '%s'", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	query += ` GROUP BY r.id ORDER BY r.trending_score DESC, r.id DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var scanRelease types.ScanRelease

		err := rows.Scan(
			&scanRelease.ID,
			&scanRelease.Title,
			pq.Array(&scanRelease.Artists),
			pq.Array(&scanRelease.Featurings),
			pq.Array(&scanRelease.ArtistsIds),
			pq.Array(&scanRelease.FeaturingsIds),
			&scanRelease.ReleaseDate,
			&scanRelease.Cover,
			pq.Array(&scanRelease.Genres),
			pq.Array(&scanRelease.Producers),
			pq.Array(&scanRelease.Tracklist),
			&scanRelease.Type,
			&scanRelease.AotyID,
			&scanRelease.TrendingScore,
		)
		if err != nil {
			return nil, err
		}

		release := types.ScanToRelease(scanRelease)
		releases = append(releases, &release)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return releases, nil
}

func GetTrendingReleases(releaseType string, direction string, reference int, period string) *model.ReleasesList {
	prev := false
	next := false

	limit := 30
	offset := reference

	startDate, endDate := tools.GetReleaseDates(period)

	if direction == "previous" {
		offset = reference - limit
		if offset < 0 {
			offset = 0
		}
	}

	count := releasesCount(releaseType, startDate, endDate)

	query := `
		SELECT r.id, r.title, array_agg(DISTINCT main.name) AS artists, array_agg(DISTINCT feature.name) AS featurings, array_agg(DISTINCT main.spotify_id) AS artists_spotify_ids, array_agg(DISTINCT feature.spotify_id) AS featurings_spotify_ids, r.date, r.cover, array_agg(DISTINCT g.type) AS genres, array_agg(DISTINCT p.name) AS producers, r.tracklist, r.type, r.aoty_id, r.trending_score
		FROM releases AS r
		JOIN releases_artists AS ra ON r.id = ra.release_id AND ra.relationship = 'main'
		JOIN artists AS main ON ra.artist_id = main.id
		LEFT JOIN releases_artists AS raf ON r.id = raf.release_id AND raf.relationship = 'feature'
		LEFT JOIN artists AS feature ON raf.artist_id = feature.id
		LEFT JOIN releases_producers AS rp ON r.id = rp.release_id
		LEFT JOIN producers AS p ON rp.producer_id = p.id
		LEFT JOIN releases_genres AS rg ON r.id = rg.release_id
		LEFT JOIN genres AS g ON rg.genre_id = g.id
		WHERE r.trending_score IS NOT NULL
	`

	switch releaseType {
	case "single":
		query += " AND r.type = 'single'"
	case "album":
		query += " AND r.type != 'single'"
	}

	query += fmt.Sprintf(" AND r.date >= '%s' AND r.date <= '%s'", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	query += " GROUP BY r.id ORDER BY r.trending_score DESC, r.id DESC"
	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var releases []*model.Release

	batchSize := 0

	for rows.Next() {
		batchSize += 1

		var scanRelease types.ScanRelease

		err := rows.Scan(
			&scanRelease.ID,
			&scanRelease.Title,
			pq.Array(&scanRelease.Artists),
			pq.Array(&scanRelease.Featurings),
			pq.Array(&scanRelease.ArtistsIds),
			pq.Array(&scanRelease.FeaturingsIds),
			&scanRelease.ReleaseDate,
			&scanRelease.Cover,
			pq.Array(&scanRelease.Genres),
			pq.Array(&scanRelease.Producers),
			pq.Array(&scanRelease.Tracklist),
			&scanRelease.Type,
			&scanRelease.AotyID,
			&scanRelease.TrendingScore,
		)
		if err != nil {
			panic(err)
		}

		release := types.ScanToRelease(scanRelease)
		releases = append(releases, &release)
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	if direction == "next" {
		if count > reference+batchSize {
			next = true
		}
		if reference > 0 {
			prev = true
		}
	} else {
		next = true
		if offset > 0 {
			prev = true
		}
	}

	return &model.ReleasesList{
		Releases: releases,
		Prev:     prev,
		Next:     next,
	}
}

func GetRelease(id int) *model.Release {
	query := `
		SELECT r.id, r.title, array_agg(DISTINCT main.name) AS artists, array_agg(DISTINCT feature.name) AS featurings, array_agg(DISTINCT main.spotify_id) AS artists_spotify_ids, array_agg(DISTINCT feature.spotify_id) AS featurings_spotify_ids, r.date, r.cover, array_agg(DISTINCT g.type) AS genres, array_agg(DISTINCT p.name) AS producers, r.tracklist, r.type, r.aoty_id, r.trending_score
		FROM releases AS r
		JOIN releases_artists AS ra ON r.id = ra.release_id AND ra.relationship = 'main'
		JOIN artists AS main ON ra.artist_id = main.id
		LEFT JOIN releases_artists AS raf ON r.id = raf.release_id AND raf.relationship = 'feature'
		LEFT JOIN artists AS feature ON raf.artist_id = feature.id
		LEFT JOIN releases_producers AS rp ON r.id = rp.release_id
		LEFT JOIN producers AS p ON rp.producer_id = p.id
		LEFT JOIN releases_genres AS rg ON r.id = rg.release_id
		LEFT JOIN genres AS g ON rg.genre_id = g.id
		WHERE r.trending_score IS NOT NULL
	`

	query += fmt.Sprintf(" AND r.id = %d", id)
	query += " GROUP BY r.id"

	row := db.QueryRow(query)

	var scanRelease types.ScanRelease

	err := row.Scan(
		&scanRelease.ID,
		&scanRelease.Title,
		pq.Array(&scanRelease.Artists),
		pq.Array(&scanRelease.Featurings),
		pq.Array(&scanRelease.ArtistsIds),
		pq.Array(&scanRelease.FeaturingsIds),
		&scanRelease.ReleaseDate,
		&scanRelease.Cover,
		pq.Array(&scanRelease.Genres),
		pq.Array(&scanRelease.Producers),
		pq.Array(&scanRelease.Tracklist),
		&scanRelease.Type,
		&scanRelease.AotyID,
		&scanRelease.TrendingScore,
	)
	if err != nil {
		panic(err)
	}

	release := types.ScanToRelease(scanRelease)

	return &release
}

func GetAllReleasesCount() *model.AllReleasesCount {
	allReleasesCount := &model.AllReleasesCount{}

	pastStartDate, pastEndDate := tools.GetReleaseDates("past")
	weekStartDate, weekEndDate := tools.GetReleaseDates("week")
	monthStartDate, monthEndDate := tools.GetReleaseDates("month")
	extendedStartDate, extendedEndDate := tools.GetReleaseDates("extended")

	query := `
		SELECT
			COUNT(*) FILTER (WHERE r.date >= $1 AND r.date <= $2) AS past_releases,
			COUNT(*) FILTER (WHERE r.date >= $1 AND r.date <= $2 AND r.type != 'single') AS past_albums,
			COUNT(*) FILTER (WHERE r.date >= $1 AND r.date <= $2 AND r.type = 'single') AS past_singles,
			COUNT(*) FILTER (WHERE r.date >= $3 AND r.date <= $4) AS week_releases,
			COUNT(*) FILTER (WHERE r.date >= $3 AND r.date <= $4 AND r.type != 'single') AS week_albums,
			COUNT(*) FILTER (WHERE r.date >= $3 AND r.date <= $4 AND r.type = 'single') AS week_singles,
			COUNT(*) FILTER (WHERE r.date >= $5 AND r.date <= $6) AS month_releases,
			COUNT(*) FILTER (WHERE r.date >= $5 AND r.date <= $6 AND r.type != 'single') AS month_albums,
			COUNT(*) FILTER (WHERE r.date >= $5 AND r.date <= $6 AND r.type = 'single') AS month_singles,
			COUNT(*) FILTER (WHERE r.date >= $7 AND r.date <= $8) AS extended_releases,
			COUNT(*) FILTER (WHERE r.date >= $7 AND r.date <= $8 AND r.type != 'single') AS extended_albums,
			COUNT(*) FILTER (WHERE r.date >= $7 AND r.date <= $8 AND r.type = 'single') AS extended_singles
		FROM releases AS r
		WHERE r.trending_score IS NOT NULL
	`

	rows, err := db.Query(query, pastStartDate, pastEndDate, weekStartDate, weekEndDate, monthStartDate, monthEndDate, extendedStartDate, extendedEndDate)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	if rows.Next() {
		var pastReleases, pastAlbums, pastSingles, weekReleases, weekAlbums, weekSingles, monthReleases, monthAlbums, monthSingles, extendedReleases, extendedAlbums, extendedSingles int

		err := rows.Scan(&pastReleases, &pastAlbums, &pastSingles, &weekReleases, &weekAlbums, &weekSingles, &monthReleases, &monthAlbums, &monthSingles, &extendedReleases, &extendedAlbums, &extendedSingles)
		if err != nil {
			panic(err)
		}

		allReleasesCount.Past = &model.ReleasesCount{
			All:     pastReleases,
			Albums:  pastAlbums,
			Singles: pastSingles,
		}

		allReleasesCount.Week = &model.ReleasesCount{
			All:     weekReleases,
			Albums:  weekAlbums,
			Singles: weekSingles,
		}

		allReleasesCount.Month = &model.ReleasesCount{
			All:     monthReleases,
			Albums:  monthAlbums,
			Singles: monthSingles,
		}

		allReleasesCount.Extended = &model.ReleasesCount{
			All:     extendedReleases,
			Albums:  extendedAlbums,
			Singles: extendedSingles,
		}
	}

	return allReleasesCount
}

func releasesCount(releaseType string, startDate time.Time, endDate time.Time) int {
	var count int

	query := `
		SELECT COUNT(*) FROM (
		SELECT r.id
		FROM Releases r
		LEFT JOIN Releases_Artists ra ON r.id = ra.release_id AND ra.relationship = 'main'
		LEFT JOIN Artists a ON ra.artist_id = a.id
		LEFT JOIN Releases_Artists raf ON r.id = raf.release_id AND raf.relationship = 'feature'
		LEFT JOIN Artists ap ON raf.artist_id = ap.id
		LEFT JOIN Releases_Producers rp ON r.id = rp.release_id
		LEFT JOIN Producers p ON rp.producer_id = p.id
		LEFT JOIN Releases_Genres rg ON r.id = rg.release_id
		LEFT JOIN Genres g ON rg.genre_id = g.id
		WHERE r.trending_score IS NOT NULL
	`

	switch releaseType {
	case "single":
		query += " AND r.type = 'single'"
	case "album":
		query += " AND r.type != 'single'"
	}

	query += fmt.Sprintf(" AND r.date >= '%s' AND r.date <= '%s'", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	query += " GROUP BY r.id ORDER BY r.trending_score DESC, r.id DESC) as all_releases"

	err := db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0
	}

	return count
}

func AddOrUpdateRelease(releaseType string, release types.Release, updateTime time.Time) (int64, error) {
	var id int64
	err := db.QueryRow("SELECT id FROM Releases WHERE aoty_id=$1", release.AOTY_Id).Scan(&id)
	if err == nil {
		_, err = db.Exec("UPDATE Releases SET title=$1, date=$2, cover=$3, tracklist=$4, type=$5, updated=$6 WHERE id=$7",
			release.Title, release.Date, release.Cover, pq.Array(release.Tracklist), releaseType, updateTime, id)
		if err == nil {
			return id, nil
		}
		return -1, err
	} else if err == sql.ErrNoRows {
		var newID int64
		err = db.QueryRow("INSERT INTO Releases(title, date, cover, tracklist, type, updated, aoty_id) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
			release.Title, release.Date, release.Cover, pq.Array(release.Tracklist), releaseType, updateTime, release.AOTY_Id).Scan(&newID)
		if err == nil {
			return newID, nil
		}
		return -1, err
	} else {
		return -1, err
	}

}

func UploadReleaseArtists(releaseId int64, artistId int64, relationship string) {
	_, err := db.Exec("INSERT INTO releases_artists (release_id, artist_id, relationship) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING", releaseId, artistId, relationship)
	if err != nil {
		fmt.Println("Error inserting release main artist join")
	}
}

func UpdateReleaseTrendingScore(releaseId int64, trendingScore float64) {
	_, err := db.Exec("UPDATE Releases SET trending_score=$1 WHERE id=$2",
		trendingScore, releaseId)
	if err != nil {
		fmt.Println("Error adding trending score to release")
	}
}

func PurgeReleases(updateTime time.Time) {
	_, err := db.Exec("DELETE FROM Releases WHERE updated!=$1", updateTime)
	if err == nil {
		fmt.Println("Old releases purged!")
	}
}
