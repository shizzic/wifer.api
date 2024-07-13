package auth

import (
	"math/rand"
	"net/http"
	"time"
	"wifer/server/structs"
)

const nums = "1234567890"
const letters = "1234567890_-abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Props = structs.Props

// Дешифровать токен из куки, 30 ms speed average
func DecryptToken(props *Props, token string, w http.ResponseWriter) (username string) {
	key := 0
	minus := 0

	for i, char := range token {
		if key == i {
			if char%2 == 0 {
				username += string(char - 1)
			} else {
				username += string(char + 1)
			}

			key += minus + 1
			minus += 1
		}
	}

	cookie := http.Cookie{
		Name:     "auth",
		Value:    "auth",
		Path:     "/",
		Domain:   "." + props.Conf.SELF_DOMAIN_NAME,
		MaxAge:   1800,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &cookie)
	return
}

// Зашифровать токен для куки
func EncryptToken(username string) (token string) {
	for i, char := range username {
		if char%2 == 0 {
			token += string(char - 1)
		} else {
			token += string(char + 1)
		}

		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		b := make([]byte, i)
		for i := range b {
			b[i] = letters[r.Int63()%int64(len(letters))]
		}

		token += string(b)
	}

	return
}

// Создать или очистеть куки для аутентификации
func MakeCookies(props *Props, id string, username string, time int, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    EncryptToken(username),
		Path:     "/",
		Domain:   "." + props.Conf.SELF_DOMAIN_NAME,
		MaxAge:   time,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    username,
		Path:     "/",
		Domain:   "." + props.Conf.SELF_DOMAIN_NAME,
		MaxAge:   time,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "id",
		Value:    id,
		Path:     "/",
		Domain:   "." + props.Conf.SELF_DOMAIN_NAME,
		MaxAge:   time,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
}
