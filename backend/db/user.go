package db

import (
	"context"
	"main/graph/model"
	"main/tools"
	"strings"

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
