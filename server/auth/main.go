package auth

import (
	"errors"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"wifer/server/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const nums = "1234567890"
const letters = "1234567890_-abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

type Props = structs.Props

// 30 ms speed average
func DecryptToken(props Props, token string, w http.ResponseWriter) (username string) {
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

// Check code's fit for ensure
func CheckCode(props Props, id int, code string, w http.ResponseWriter) error {
	if !isCode(code) {
		return errors.New("0")
	}

	var data bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 1})

	if err := props.DB["ensure"].FindOne(props.Ctx, bson.M{"_id": id, "code": code}, opts).Decode(&data); err != nil {
		return errors.New("1")
	}

	// Delete document in ensure collection, if given code was valid
	props.DB["ensure"].DeleteOne(props.Ctx, bson.M{"_id": id, "code": code})

	var user bson.M
	opt := options.FindOne().SetProjection(bson.M{"username": 1})

	if err := props.DB["users"].FindOne(props.Ctx, bson.M{"_id": id}, opt).Decode(&user); err != nil {
		return errors.New("2")
	}

	props.DB["users"].UpdateOne(props.Ctx, bson.M{"_id": id}, bson.D{{Key: "$set", Value: bson.D{
		{Key: "status", Value: true},
		{Key: "active", Value: true},
	}}})

	MakeCookies(props, strconv.Itoa(id), user["username"].(string), 86400*120, w)

	return nil
}

// Cookies for auth
func MakeCookies(props Props, id string, username string, time int, w http.ResponseWriter) {
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

// Make token for auth any email operations or something :)
func MakeCode() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, 6)
	for i := range b {
		b[i] = nums[r.Int63()%int64(len(nums))]
	}

	return string(b)
}
