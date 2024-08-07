package oauth2

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"wifer/server/structs"
)

// Получаю токен с помощью кода
func IsVK(props *structs.Props, data *structs.Signin) (string, error) {
	params := url.Values{}
	params.Add("grant_type", `authorization_code`)
	params.Add("code_verifier", props.Conf.VK_SECRET)
	params.Add("client_id", props.Conf.VK_ID)
	params.Add("state", data.State)
	params.Add("redirect_uri", data.Redirect)
	params.Add("device_id", data.Device)
	params.Add("code", data.Token)
	body := strings.NewReader(params.Encode())
	req, err := http.NewRequest("POST", "https://id.vk.com/oauth2/auth", body)
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
	var ready map[string]any
	json.Unmarshal(result, &ready)
	token := ready["access_token"].(string) // получил токен для финального запроса

	email, err := get_vk_email(props, token)
	if err != nil {
		return "", errors.New("wrong_api_token")
	}
	return email, nil
}

// через полученные токен получаю информацию о юзере (доступную мне)
func get_vk_email(props *structs.Props, token string) (string, error) {
	params := url.Values{}
	params.Add("client_id", props.Conf.VK_ID)
	params.Add("access_token", token)
	body := strings.NewReader(params.Encode())
	req, err := http.NewRequest("POST", "https://id.vk.com/oauth2/user_info", body)
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
	email := ready["user"].(map[string]string)["email"]
	return email, nil
}
