package db

import (
	"fmt"
	"main/models"
	"main/services"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

func GetRecommendations(client *http.Client, period string) []models.DisplayRelease {

	fmt.Println("Fetching recommendations...")
	artists, tracks, err := services.SpotifyUserTops(client)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	artists, err = services.SpotifyRecommendations(artists, tracks)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	artists, err = services.SpotifyRelatedArtists(artists)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	genres := services.TopGenres(artists)

	artistIds := GetMatchingArtists(artists, genres)
	releases, _ := GetMatchingReleases(artistIds, genres, period)

	return releases
}

func GetMatchingArtists(artists []models.SpotifyArtist, genres []string) []int {
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

func GetMatchingReleases(ids []int, genres []string, period string) ([]models.DisplayRelease, error) {
	var releases []models.DisplayRelease

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
		var release models.DisplayRelease
		var artists, featurings, genres, producers []string

		if err := rows.Scan(&release.ID, pq.Array(&artists), pq.Array(&featurings), &release.Title, &release.Date, &release.Cover, pq.Array(&genres), pq.Array(&producers), pq.Array(&release.Tracklist), &release.Type, &release.AOTYID, &release.TrendingScore); err != nil {
			return nil, err
		}

		release.Artists = artists
		release.Featurings = featurings
		release.Genres = genres
		release.Producers = producers

		releases = append(releases, release)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return releases, nil
}
