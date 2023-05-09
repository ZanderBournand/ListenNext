package db

import (
	"fmt"
	"main/services"
	"strconv"
	"strings"

	"github.com/lib/pq"
)

func GetRecommendations(artists []services.SpotifyArtist) []DisplayRelease {

	artistIds := GetMatchingArtists(artists)

	releases, _ := GetMatchingReleases(artistIds, artists)

	return releases

}

func GetMatchingArtists(artists []services.SpotifyArtist) []int {
	var ids []int

	var allIDs []string
	var allGenres []string
	for _, artist := range artists {
		allIDs = append(allIDs, "'"+artist.ID+"'")
		for _, genre := range artist.Genres {
			compareType := strings.ToLower(strings.ReplaceAll(genre, " ", ""))
			allGenres = append(allGenres, "'"+compareType+"'")
		}
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

func GetMatchingReleases(ids []int, artists []services.SpotifyArtist) ([]DisplayRelease, error) {
	var releases []DisplayRelease

	idStrings := make([]string, 0, len(ids))
	for _, id := range ids {
		idStrings = append(idStrings, strconv.Itoa(id))
	}

	var allGenres []string
	for _, artist := range artists {
		for _, genre := range artist.Genres {
			compareType := strings.ToLower(strings.ReplaceAll(genre, " ", ""))
			allGenres = append(allGenres, "'"+compareType+"'")
		}
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
		GROUP BY r.id, r.title, r.date, r.cover, r.tracklist, r.type, r.aoty_id, r.trending_score
		ORDER BY r.trending_score DESC
	`, strings.Join(idStrings, ","), strings.Join(allGenres, ","))

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var release DisplayRelease
		var artists, featurings, genres, producers []string

		if err := rows.Scan(&release.ID, pq.Array(&artists), pq.Array(&featurings), &release.Title, &release.Date, &release.Cover, pq.Array(&genres), pq.Array(&producers), pq.Array(&release.Tracklist), &release.Type, &release.AOTYID, &release.TrendingScore); err != nil {
			return nil, err
		}

		releases = append(releases, release)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return releases, nil
}
