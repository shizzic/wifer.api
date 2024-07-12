package structs

import (
	"context"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	MONGO_CONNECTION_STRING string

	SERVER_IP         string
	CLIENT_DOMAIN     string
	SELF_DOMAIN_NAME  string
	ENCRYPT_CERT_FILE string
	ENCRYPT_KEY_FILE  string

	ADMIN_EMAIL string
	EMAIL       Email
	PATH        string

	BACKBLAZE_ID  string
	BACKBLAZE_KEY string
	PRODUCT_NAME  string
}

type Auth struct {
	ID   int    `query:"id"`   // id пользователя
	Code string `query:"code"` // код валидации
}

type Email struct {
	HOST     string
	USERNAME string
	PASSWORD string
	PORT     int
}

type Props struct {
	Conf Config
	Ctx  context.Context
	DB   map[string]*mongo.Collection
	R    *chi.Mux
}

// Поля пользователя
type User struct {
	ID         int    `query:"id"`
	Username   string `query:"username"`
	Email      string `query:"email"`
	Title      string `query:"title"`
	About      string `query:"about"`
	Country_id int    `json:"country_id"`
	City_id    int    `json:"city_id"`
	Country    int    `query:"country_id"`
	City       int    `query:"city_id"`
	Sex        int    `query:"sex"`
	Age        int    `query:"age"`
	Height     int    `query:"height"`
	Weight     int    `query:"weight"`
	Body       int    `query:"body"`
	Smokes     int    `query:"smokes"`
	Drinks     int    `query:"drinks"`
	Ethnicity  int    `query:"ethnicity"`
	Search     []int  `query:"search"`
	Prefer     int    `query:"prefer"`
	Income     int    `query:"income"`
	Children   int    `query:"children"`
	Industry   int    `query:"industry"`
	Online     bool   `query:"online"`
	Premium    int64  `json:"premium"`
}

// Набор действий открытого профиля в отношении открывающего
type Target struct {
	Like    bson.M   // Лайкнул или нет
	Private []bson.M // Дал доступ к приватным фото или нет
	Access  []bson.M // Дал доступ к переписке
}
