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

// Получаю токен с помощью кода
func IsYandex(props *structs.Props, code string) (string, error) {
	params := url.Values{}
	params.Add("grant_type", `authorization_code`)
	params.Add("format", `json`)
	params.Add("code", code)
	params.Add("client_id", props.Conf.YANDEX_ID)
	params.Add("client_secret", props.Conf.YANDEX_SECRET)
	body := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", "https://oauth.yandex.ru/token", body)
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

	email, err := get_yandex_email(token)
	if err != nil {
		return "", err
	}
	return email, nil
}

// через полученные токен получаю информацию о юзере (доступную мне)
func get_yandex_email(token string) (string, error) {
	req, err := http.NewRequest("GET", "https://login.yandex.ru/info?format=json", nil)
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
	email := ready["default_email"]

	// Валидирую почту
	if email != nil && auth.IsEmailValid(email.(string)) {
		return email.(string), nil
	} else {
		return "", errors.New("email_not_verified")
	}
}
