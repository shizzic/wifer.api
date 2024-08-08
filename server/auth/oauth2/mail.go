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
func IsMail(props *structs.Props, data *structs.Signin) (string, error) {
	params := url.Values{}
	params.Add("grant_type", `authorization_code`)
	params.Add("redirect_uri", data.Redirect)
	params.Add("code", data.Token)
	body := strings.NewReader(params.Encode())
	req, err := http.NewRequest("POST", "https://oauth.mail.ru/token", body)
	if err != nil {
		return "", errors.New("wrong_api_token")
	}
	req.SetBasicAuth(props.Conf.MAIL_ID, props.Conf.MAIL_SECRET)
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

	email, err := get_mail_email(token)
	if err != nil {
		return "", err
	}
	return email, nil
}

// через полученные токен получаю информацию о юзере (доступную мне)
func get_mail_email(token string) (string, error) {
	response, err := http.Get("https://oauth.mail.ru/userinfo?access_token=" + token)
	if err != nil {
		return "", errors.New("wrong_api_token")
	}
	defer response.Body.Close()

	result, _ := io.ReadAll(response.Body)
	var ready map[string]interface{}
	json.Unmarshal(result, &ready)
	email := ready["email"]

	// Валидирую почту
	if email != nil && auth.IsEmailValid(email.(string)) {
		return email.(string), nil
	} else {
		return "", errors.New("email_not_verified")
	}
}
