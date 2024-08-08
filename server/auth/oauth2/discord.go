package oauth2

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"wifer/server/auth"
	"wifer/server/structs"
)

func IsDiscord(props *structs.Props, data *structs.Signin) (string, error) {
	params := url.Values{}
	params.Add("grant_type", `client_credentials`)
	params.Add("redirect_uri", data.Redirect)
	params.Add("client_id", props.Conf.DISCORD_ID)
	params.Add("client_secret", props.Conf.DISCORD_SECRET)
	params.Add("scope", `identify email`)
	params.Add("code", data.Token)
	body := strings.NewReader(params.Encode())
	req, err := http.NewRequest("POST", "https://discord.com/api/oauth2/token", body)
	if err != nil {
		return "", errors.New("wrong_api_token")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.New("wrong_api_token")
	}
	defer response.Body.Close()

	result, _ := io.ReadAll(response.Body)
	var ready map[string]interface{}
	json.Unmarshal(result, &ready)
	token := ready["access_token"].(string) // получил токен для финального запроса

	email, err := get_discord_email(token)
	if err != nil {
		return "", err
	}
	return email, nil
}

// через полученные токен получаю информацию о юзере (доступную мне)
func get_discord_email(token string) (string, error) {
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		return "", errors.New("wrong_api_token")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.New("wrong_api_token")
	}
	defer response.Body.Close()

	result, _ := io.ReadAll(response.Body)
	var ready map[string]interface{}
	json.Unmarshal(result, &ready)
	email := ready["email"]
	verified := ready["verified"]

	// Валидирую почту
	if verified != nil && verified.(bool) && email != nil && auth.IsEmailValid(email.(string)) {
		return email.(string), nil
	} else {
		return "", errors.New("email_not_verified")
	}
}
