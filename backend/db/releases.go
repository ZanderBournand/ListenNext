package db

import (
	"database/sql"
	"fmt"
	"main/graph/model"
	"main/types"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

func GetMatchingReleases(ids []int, genres []string, period string) ([]*model.Release, error) {
	var releases []*model.Release

	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	startDate := now
	endDate := now

	if period == "past" {
		startDate = startDate.AddDate(0, 0, -14)
		daysUntilFriday := (5 - int(startDate.Weekday()) + 7) % 7
		startDate = startDate.AddDate(0, 0, daysUntilFriday+1)
		endDate = endDate.AddDate(0, 0, -7)
		daysUntilFriday = (5 - int(endDate.Weekday()) + 7) % 7
		endDate = endDate.AddDate(0, 0, daysUntilFriday)
	} else if period == "week" {
		startDate = startDate.AddDate(0, 0, -7)
		daysUntilFriday := (5 - int(startDate.Weekday()) + 7) % 7
		startDate = startDate.AddDate(0, 0, daysUntilFriday+1)
		daysUntilFriday = (5 - int(endDate.Weekday()) + 7) % 7
		endDate = endDate.AddDate(0, 0, daysUntilFriday)
	} else if period == "month" {
		startDate = startDate.AddDate(0, 0, -7)
		daysUntilFriday := (5 - int(startDate.Weekday()) + 7) % 7
		startDate = startDate.AddDate(0, 0, daysUntilFriday+1)
		endDate = endDate.AddDate(0, 0, 28)
		daysUntilFriday = (5 - int(endDate.Weekday()) + 7) % 7
		endDate = endDate.AddDate(0, 0, daysUntilFriday)
	} else if period == "extended" {
		startDate = startDate.AddDate(0, 0, -7)
		daysUntilFriday := (5 - int(startDate.Weekday()) + 7) % 7
		startDate = startDate.AddDate(0, 0, daysUntilFriday+1)
		endDate = endDate.AddDate(0, 0, 84)
		daysUntilFriday = (5 - int(endDate.Weekday()) + 7) % 7
		endDate = endDate.AddDate(0, 0, daysUntilFriday)
	}

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
		SELECT r.id, array_agg(DISTINCT a.name) AS artists, array_agg(DISTINCT COALESCE(ap.name, '')) AS featurings, r.title, r.date, r.cover, array_agg(DISTINCT COALESCE(g.type, '')) AS genres, array_agg(DISTINCT COALESCE(p.name, '')) AS producers, r.tracklist, r.type, r.aoty_id, r.trending_score
		FROM Releases r
		LEFT JOIN Releases_Artists ra ON r.id = ra.release_id AND ra.relationship = 'main'
		LEFT JOIN Artists a ON ra.artist_id = a.id
		LEFT JOIN Releases_Artists raf ON r.id = raf.release_id AND raf.relationship = 'feature'
		LEFT JOIN Artists ap ON raf.artist_id = ap.id
		LEFT JOIN Releases_Producers rp ON r.id = rp.release_id
		LEFT JOIN Producers p ON rp.producer_id = p.id
		LEFT JOIN Releases_Genres rg ON r.id = rg.release_id
		LEFT JOIN Genres g ON rg.genre_id = g.id
		WHERE r.trending_score IS NOT NULL AND (a.id IN (%s) OR g.compare_type IN (%s))
	`, strings.Join(idStrings, ","), strings.Join(allGenres, ","))

	query += fmt.Sprintf(" AND r.date >= '%s' AND r.date <= '%s'", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	query += ` GROUP BY r.id, r.title, r.date, r.cover, r.tracklist, r.type, r.aoty_id, r.trending_score
	ORDER BY r.trending_score DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var release model.Release
		var artists, featurings, genres, producers []string

		if err := rows.Scan(&release.ID, pq.Array(&artists), pq.Array(&featurings), &release.Title, &release.Date, &release.Cover, pq.Array(&genres), pq.Array(&producers), pq.Array(&release.Tracklist), &release.Type, &release.AotyID, &release.TrendingScore); err != nil {
			return nil, err
		}

		release.Artists = artists
		release.Featurings = featurings
		release.Genres = genres
		release.Producers = producers

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

	now := time.Now()
	now = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	startDate := now
	endDate := now

	if period == "past" {
		startDate = startDate.AddDate(0, 0, -14)
		daysUntilFriday := (5 - int(startDate.Weekday()) + 7) % 7
		startDate = startDate.AddDate(0, 0, daysUntilFriday+1)
		endDate = endDate.AddDate(0, 0, -7)
		daysUntilFriday = (5 - int(endDate.Weekday()) + 7) % 7
		endDate = endDate.AddDate(0, 0, daysUntilFriday)
	} else if period == "week" {
		startDate = startDate.AddDate(0, 0, -7)
		daysUntilFriday := (5 - int(startDate.Weekday()) + 7) % 7
		startDate = startDate.AddDate(0, 0, daysUntilFriday+1)
		daysUntilFriday = (5 - int(endDate.Weekday()) + 7) % 7
		endDate = endDate.AddDate(0, 0, daysUntilFriday)
	} else if period == "month" {
		startDate = startDate.AddDate(0, 0, -7)
		daysUntilFriday := (5 - int(startDate.Weekday()) + 7) % 7
		startDate = startDate.AddDate(0, 0, daysUntilFriday+1)
		endDate = endDate.AddDate(0, 0, 28)
		daysUntilFriday = (5 - int(endDate.Weekday()) + 7) % 7
		endDate = endDate.AddDate(0, 0, daysUntilFriday)
	} else if period == "extended" {
		startDate = startDate.AddDate(0, 0, -7)
		daysUntilFriday := (5 - int(startDate.Weekday()) + 7) % 7
		startDate = startDate.AddDate(0, 0, daysUntilFriday+1)
		endDate = endDate.AddDate(0, 0, 84)
		daysUntilFriday = (5 - int(endDate.Weekday()) + 7) % 7
		endDate = endDate.AddDate(0, 0, daysUntilFriday)
	}

	if direction == "previous" {
		offset = reference - limit
		if offset < 0 {
			offset = 0
		}
	}

	count := releasesCount(releaseType, startDate, endDate)

	query := `SELECT r.id, array_agg(DISTINCT a.name) AS artists, array_agg(DISTINCT COALESCE(ap.name, '')) AS featurings, r.title, r.date, r.cover, array_agg(DISTINCT COALESCE(g.type, '')) AS genres, array_agg(DISTINCT COALESCE(p.name, '')) AS producers, r.tracklist, r.type, r.aoty_id, r.trending_score
              FROM Releases r
              LEFT JOIN Releases_Artists ra ON r.id = ra.release_id AND ra.relationship = 'main'
              LEFT JOIN Artists a ON ra.artist_id = a.id
              LEFT JOIN Releases_Artists raf ON r.id = raf.release_id AND raf.relationship = 'feature'
              LEFT JOIN Artists ap ON raf.artist_id = ap.id
              LEFT JOIN Releases_Producers rp ON r.id = rp.release_id
              LEFT JOIN Producers p ON rp.producer_id = p.id
              LEFT JOIN Releases_Genres rg ON r.id = rg.release_id
              LEFT JOIN Genres g ON rg.genre_id = g.id
			  WHERE r.trending_score IS NOT NULL`

	switch releaseType {
	case "single":
		query += " AND r.type = 'single'"
	case "album":
		query += " AND r.type != 'single'"
	}

	query += fmt.Sprintf(" AND r.date >= '%s' AND r.date <= '%s'", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	query += " GROUP BY r.id ORDER BY r.trending_score DESC"
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

		var release model.Release
		var artists, featurings, genres, producers []string

		err := rows.Scan(&release.ID, pq.Array(&artists), pq.Array(&featurings), &release.Title, &release.Date, &release.Cover, pq.Array(&genres), pq.Array(&producers), pq.Array(&release.Tracklist), &release.Type, &release.AotyID, &release.TrendingScore)
		if err != nil {
			panic(err)
		}

		release.Artists = artists
		release.Featurings = featurings
		release.Genres = genres
		release.Producers = producers

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
	query := `SELECT r.id, array_agg(DISTINCT a.name) AS artists, array_agg(DISTINCT COALESCE(ap.name, '')) AS featurings, r.title, r.date, r.cover, array_agg(DISTINCT COALESCE(g.type, '')) AS genres, array_agg(DISTINCT COALESCE(p.name, '')) AS producers, r.tracklist, r.type, r.aoty_id, r.trending_score
              FROM Releases r
              LEFT JOIN Releases_Artists ra ON r.id = ra.release_id AND ra.relationship = 'main'
              LEFT JOIN Artists a ON ra.artist_id = a.id
              LEFT JOIN Releases_Artists raf ON r.id = raf.release_id AND raf.relationship = 'feature'
              LEFT JOIN Artists ap ON raf.artist_id = ap.id
              LEFT JOIN Releases_Producers rp ON r.id = rp.release_id
              LEFT JOIN Producers p ON rp.producer_id = p.id
              LEFT JOIN Releases_Genres rg ON r.id = rg.release_id
              LEFT JOIN Genres g ON rg.genre_id = g.id
			  WHERE r.trending_score IS NOT NULL`

	query += fmt.Sprintf(" AND r.id = %d", id)
	query += " GROUP BY r.id"

	row := db.QueryRow(query)

	var release model.Release
	var artists, featurings, genres, producers []string

	err := row.Scan(&release.ID, pq.Array(&artists), pq.Array(&featurings), &release.Title, &release.Date, &release.Cover, pq.Array(&genres), pq.Array(&producers), pq.Array(&release.Tracklist), &release.Type, &release.AotyID, &release.TrendingScore)
	if err != nil {
		panic(err)
	}

	release.Artists = artists
	release.Featurings = featurings
	release.Genres = genres
	release.Producers = producers

	return &release
}

func releasesCount(releaseType string, startDate time.Time, endDate time.Time) int {

	var count int

	query := `SELECT COUNT(*) FROM (
	SELECT r.id, array_agg(DISTINCT a.name) AS artists, array_agg(DISTINCT COALESCE(ap.name, '')) AS featurings, r.title, r.date, r.cover, array_agg(DISTINCT COALESCE(g.type, '')) AS genres, array_agg(DISTINCT COALESCE(p.name, '')) AS producers, r.tracklist, r.type, r.aoty_id, r.trending_score
	FROM Releases r
	LEFT JOIN Releases_Artists ra ON r.id = ra.release_id AND ra.relationship = 'main'
	LEFT JOIN Artists a ON ra.artist_id = a.id
	LEFT JOIN Releases_Artists raf ON r.id = raf.release_id AND raf.relationship = 'feature'
	LEFT JOIN Artists ap ON raf.artist_id = ap.id
	LEFT JOIN Releases_Producers rp ON r.id = rp.release_id
	LEFT JOIN Producers p ON rp.producer_id = p.id
	LEFT JOIN Releases_Genres rg ON r.id = rg.release_id
	LEFT JOIN Genres g ON rg.genre_id = g.id
	WHERE r.trending_score IS NOT NULL`

	switch releaseType {
	case "single":
		query += " AND r.type = 'single'"
	case "album":
		query += " AND r.type != 'single'"
	}

	query += fmt.Sprintf(" AND r.date >= '%s' AND r.date <= '%s'", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	query += " GROUP BY r.id ORDER BY r.trending_score DESC) as all_releases"

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
		fmt.Println("error updating release")
		return -1, fmt.Errorf("error updating release")
	} else if err == sql.ErrNoRows {
		var newID int64
		err = db.QueryRow("INSERT INTO Releases(title, date, cover, tracklist, type, updated, aoty_id) VALUES($1, $2, $3, $4, $5, $6, $7) RETURNING id",
			release.Title, release.Date, release.Cover, pq.Array(release.Tracklist), releaseType, updateTime, release.AOTY_Id).Scan(&newID)
		if err == nil {
			return newID, nil
		}
		fmt.Println("error addding release")
		return -1, fmt.Errorf("error addding release")
	} else {
		fmt.Println(err)
		fmt.Println("error querying release")
		return -1, fmt.Errorf("error querying release")
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
