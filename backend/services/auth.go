package services

import (
	"context"
	"database/sql"
	"main/config"
	"main/db"
	"main/graph/model"
	"main/tools"

	"github.com/vektah/gqlparser/v2/gqlerror"
)

func UserRegister(ctx context.Context, input model.NewUser) (interface{}, error) {
	user, err := db.UserGetByEmail(ctx, input.Email)
	if err == nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	createdUser, err := db.UserCreate(ctx, input)
	if err != nil {
		return nil, err
	}

	jwtToken, err := tools.JwtGenerate(ctx, createdUser.ID, "native")
	if err != nil {
		return nil, err
	}

	userClient := model.UserClient{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
	}

	return map[string]interface{}{
		"user":  userClient,
		"token": jwtToken,
	}, nil
}

func UserLogin(ctx context.Context, email string, password string) (interface{}, error) {
	user, err := db.UserGetByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, &gqlerror.Error{
				Message: "Email not found",
			}
		}
		return nil, err
	}

	if err := tools.ComparePassword(user.Password, password); err != nil {
		return nil, err
	}

	jwtToken, err := tools.JwtGenerate(ctx, user.ID, "native")
	if err != nil {
		return nil, err
	}

	userClient := model.UserClient{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
	}

	return map[string]interface{}{
		"user":  userClient,
		"token": jwtToken,
	}, nil
}

func SpotifyUserLogin(ctx context.Context, code string) (interface{}, error) {

	token, err := config.SpotifyOAuth.Exchange(context.Background(), code)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: "Spotify token exchange error",
		}
	}

	email, name, err := SpotifyUserInfo(token.AccessToken)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: "Spotify user fetch error",
		}
	}

	user, err := db.UserGetByEmail(ctx, email)
	if err == nil {
		if err != sql.ErrNoRows {
			db.UpdateSpotifyUserTokens(user.ID, token.AccessToken, token.RefreshToken, token.Expiry)
		}
	} else {
		user, err = db.SpotifyUserCreate(ctx, email, name, token.AccessToken, token.RefreshToken, token.Expiry)
		if err != nil {
			return nil, err
		}
	}

	jwtToken, err := tools.JwtGenerate(ctx, user.ID, "spotify")
	if err != nil {
		return nil, err
	}

	userClient := model.UserClient{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.DisplayName,
	}

	return map[string]interface{}{
		"user":  userClient,
		"token": jwtToken,
	}, nil
}
