package auth

import (
	"context"
	"errors"

	fb "github.com/huandu/facebook/v2"
	"google.golang.org/api/idtoken"
)

func IsGoogle(id, token string) (string, error) {
	data, err := idtoken.Validate(context.Background(), token, id)

	if err != nil {
		return "", errors.New("wrong_api_token")
	}

	return data.Claims["email"].(string), nil
}

func IsFacebook(id, token string) (string, error) {
	data, err := fb.Get("/"+id, fb.Params{
		"fields":       "email",
		"access_token": token,
	})

	if err != nil {
		return "", errors.New("wrong_api_token")
	}

	return data["email"].(string), nil
}
