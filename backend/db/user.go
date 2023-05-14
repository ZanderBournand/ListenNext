package db

import (
	"context"
	"fmt"
	"main/graph/model"
	"main/tools"
	"strings"
	"time"

	"github.com/google/uuid"
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
		fmt.Println("DADAD")
		return "", "", time.Time{}
	}

	return accessToken, refreshToken, tokenExpiration
}
