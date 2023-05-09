package db

import (
	"fmt"
	"time"

	"github.com/lib/pq"
)

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

func GetTrendings(releaseType string, direction string, reference int, period string) ([]DisplayRelease, bool, bool) {
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

	var releases []DisplayRelease

	batchSize := 0

	for rows.Next() {
		batchSize += 1

		var release DisplayRelease
		var artists, featurings, genres, producers []string

		err := rows.Scan(&release.ID, pq.Array(&artists), pq.Array(&featurings), &release.Title, &release.Date, &release.Cover, pq.Array(&genres), pq.Array(&producers), pq.Array(&release.Tracklist), &release.Type, &release.AOTYID, &release.TrendingScore)
		if err != nil {
			panic(err)
		}

		release.Artists = artists
		release.Featurings = featurings
		release.Genres = genres
		release.Producers = producers

		releases = append(releases, release)
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

	return releases, prev, next
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
