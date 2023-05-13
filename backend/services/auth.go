package services

import (
	"context"
	"database/sql"
	"main/db"
	"main/graph/model"
	"main/tools"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

func UserRegister(ctx context.Context, input model.NewUser) (interface{}, error) {
	// Check Email
	_, err := db.UserGetByEmail(ctx, input.Email)
	if err == nil {
		// if err != record not found
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	createdUser, err := db.UserCreate(ctx, input)
	if err != nil {
		return nil, err
	}

	token, err := tools.JwtGenerate(ctx, createdUser.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
	}, nil
}

func UserLogin(ctx context.Context, email string, password string) (interface{}, error) {
	getUser, err := db.UserGetByEmail(ctx, email)
	if err != nil {
		// if user not found
		if err == sql.ErrNoRows {
			return nil, &gqlerror.Error{
				Message: "Email not found",
			}
		}
		return nil, err
	}

	if err := tools.ComparePassword(getUser.Password, password); err != nil {
		return nil, err
	}

	token, err := tools.JwtGenerate(ctx, getUser.ID)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"token": token,
	}, nil
}
