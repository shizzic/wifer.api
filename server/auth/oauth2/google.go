package oauth2

import (
	"context"
	"errors"

	"google.golang.org/api/idtoken"
)

func IsGoogle(id, token string) (string, error) {
	data, err := idtoken.Validate(context.Background(), token, id)

	if err != nil {
		return "", errors.New("wrong_api_token")
	}

	return data.Claims["email"].(string), nil
}
