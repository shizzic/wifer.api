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

func IsGoogle(props *structs.Props, data *structs.Signin) (string, error) {
	params := url.Values{}
	params.Add("code", data.Token)
	params.Add("client_id", props.Conf.GOOGLE_ID)
	params.Add("client_secret", props.Conf.GOOGLE_SECRET)
	params.Add("redirect_uri", data.Redirect)
	params.Add("grant_type", `authorization_code`)
	body := strings.NewReader(params.Encode())
	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", body)
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

	email, err := get_google_email(token)
	if err != nil {
		return "", errors.New("wrong_api_token")
	}
	return email, nil
}

func get_google_email(token string) (string, error) {
	response, err := http.Get("https://www.googleapis.com/oauth2/v1/userinfo?access_token=" + token)
	if err != nil {
		return "", errors.New("wrong_api_token")
	}
	defer response.Body.Close()

	result, _ := io.ReadAll(response.Body)
	var ready map[string]interface{}
	json.Unmarshal(result, &ready)
	email := ready["email"]
	verified := ready["verified_email"]

	// Валидирую почту
	if verified != nil && verified.(bool) && email != nil && auth.IsEmailValid(email.(string)) {
		return email.(string), nil
	} else {
		return "", errors.New("wrong_api_token")
	}
}
