package db

import (
	"context"
	"fmt"
	"main/graph/model"
	"main/tools"
	"main/types"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

func UserCreate(ctx context.Context, input model.NewUser) (*model.User, error) {
	input.Password = tools.HashPassword(input.Password)

	user := model.User{
		ID:          uuid.New().String(),
		DisplayName: input.DisplayName,
		Email:       strings.ToLower(input.Email),
		Password:    input.Password,
	}

	query := `INSERT INTO users (id, display_name, email, password) VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(query, user.ID, user.DisplayName, user.Email, user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UserGetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User

	err := db.QueryRow("SELECT id, display_name, email, password FROM users WHERE id = $1", id).Scan(&user.ID, &user.DisplayName, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UserGetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	err := db.QueryRow("SELECT id, display_name, email, password FROM users WHERE email = $1", email).Scan(&user.ID, &user.DisplayName, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func UpdateSpotifyUserTokens(id string, accessToken string, refreshToken string, tokenExpiration time.Time) error {
	query := `UPDATE users SET access_token=$1, refresh_token=$2, token_expiration=$3 WHERE id=$4`

	_, err := db.Exec(query, accessToken, refreshToken, tokenExpiration, id)
	if err != nil {
		return err
	}

	return nil
}

func SpotifyUserCreate(ctx context.Context, email string, name string, accessToken string, refreshToken string, tokenExpiration time.Time) (*model.User, error) {
	user := model.User{
		ID:          uuid.New().String(),
		DisplayName: name,
		Email:       strings.ToLower(email),
	}

	query := `INSERT INTO users (id, display_name, email, access_token, refresh_token, token_expiration) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := db.Exec(query, user.ID, user.DisplayName, user.Email, accessToken, refreshToken, tokenExpiration)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetSpotifyUserTokens(id string) (string, string, time.Time) {
	var accessToken string
	var refreshToken string
	var tokenExpiration time.Time

	err := db.QueryRow("SELECT access_token, refresh_token, token_expiration FROM users WHERE id = $1", id).Scan(&accessToken, &refreshToken, &tokenExpiration)
	if err != nil {
		fmt.Println(err)
		return "", "", time.Time{}
	}

	return accessToken, refreshToken, tokenExpiration
}

func UploadUserRecommendations(UserId string, releases []*model.Release) {
	stmt, err := db.Prepare("INSERT INTO users_recommendations (user_id, release_id) VALUES ($1, $2)")
	if err != nil {
		return
	}
	defer stmt.Close()

	for _, release := range releases {
		_, err := stmt.Exec(UserId, release.ID)
		if err != nil {
			continue
		}
	}
}

func UpdateUserLastRecommendationsTimestamp(UserId string) {
	query := `UPDATE users SET last_recommendations=$1 WHERE id=$2`

	utcNow := time.Now().UTC()

	db.Exec(query, utcNow, UserId)
}

func GetUserLastRecommendationsTimestamp(UserId string) time.Time {
	var timestamp time.Time

	err := db.QueryRow("SELECT last_recommendations FROM users WHERE id = $1", UserId).Scan(&timestamp)
	if err != nil {
		fmt.Println(err)
		return time.Time{}
	}

	return timestamp
}

func GetUserRecommendations(UserId string, period string) ([]*model.Release, error) {
	startDate, endDate := tools.GetReleaseDates(period)

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
		LEFT JOIN users_recommendations as ur ON r.id = ur.release_id
		WHERE r.trending_score IS NOT NULL
	`

	query += fmt.Sprintf(" AND r.date >= '%s' AND r.date <= '%s'", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))
	query += fmt.Sprintf(" AND ur.user_id = '%s'", UserId)
	query += " GROUP BY r.id ORDER BY r.trending_score DESC, r.id DESC"

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var releases []*model.Release

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
