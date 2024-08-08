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

// Получаю токен с помощью кода (работает только на https)
func IsTwitch(props *structs.Props, data *structs.Signin) (string, error) {
	params := url.Values{}
	params.Add("client_id", props.Conf.TWITCH_ID)
	params.Add("client_secret", props.Conf.TWITCH_SECRET)
	params.Add("code", data.Token)
	params.Add("grant_type", `authorization_code`)
	params.Add("redirect_uri", data.Redirect)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://id.twitch.tv/oauth2/token", body)
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
	token := ready["access_token"].(string)

	email, err := validate_twitch_token(props, token)
	if err != nil {
		return "", errors.New("wrong_api_token")
	}
	return email, nil
}

// Валидирую полученный токен для его проверки + получаю user_id (чтобы знать почту какого юзера я вообще получаю)
// Такой вот кек от Твича
func validate_twitch_token(props *structs.Props, token string) (string, error) {
	req, err := http.NewRequest("GET", "https://id.twitch.tv/oauth2/validate", nil)
	if err != nil {
		return "", errors.New("wrong_api_token")
	}
	req.Header.Set("Authorization", "OAuth "+token)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.New("wrong_api_token")
	}
	defer response.Body.Close()

	result, _ := io.ReadAll(response.Body)
	var ready map[string]interface{}
	json.Unmarshal(result, &ready)

	email, err := get_twitch_email(props, ready["user_id"].(string), token)
	if err != nil {
		return "", errors.New("wrong_api_token")
	}

	return email, nil
}

func get_twitch_email(props *structs.Props, user_id, token string) (string, error) {
	req, err := http.NewRequest("GET", "https://api.twitch.tv/helix/users?id="+user_id, nil)
	if err != nil {
		return "", errors.New("wrong_api_token")
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Client-Id", props.Conf.TWITCH_ID)
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", errors.New("wrong_api_token")
	}
	defer response.Body.Close()

	result, _ := io.ReadAll(response.Body)
	var ready map[string]interface{}
	json.Unmarshal(result, &ready)
	var data []interface{} = ready["data"].([]interface{})

	// Валидирую почту
	if len(data) == 1 {
		data := data[0]
		email := data.(map[string]interface{})["email"]

		if email != nil && auth.IsEmailValid(email.(string)) {
			return email.(string), nil
		} else {
			return "", errors.New("wrong_api_token")
		}
	}

	return "", errors.New("wrong_api_token")
}
