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
	ID   int    `json:"id" query:"id"`     // id пользователя
	Code string `json:"code" query:"code"` // код валидации
}

type Email struct {
	HOST     string
	USERNAME string
	PASSWORD string
	PORT     int
}

type Props struct {
	Conf *Config
	Ctx  context.Context
	DB   map[string]*mongo.Collection
	R    *chi.Mux
}

// Поля пользователя
type User struct {
	ID        int    `json:"id" query:"id"`
	Username  string `json:"username" query:"username"`
	Email     string `json:"email" query:"email"`
	Title     string `json:"title" query:"title"`
	About     string `json:"about" query:"about"`
	Country   int    `json:"country_id" query:"country_id"` // country_id
	City      int    `json:"city_id" query:"city_id"`       // city_id
	Sex       int    `json:"sex" query:"sex"`
	Age       int    `json:"age" query:"age"`
	Height    int    `json:"height" query:"height"`
	Weight    int    `json:"weight" query:"weight"`
	Body      int    `json:"body" query:"body"`
	Smokes    int    `json:"smokes" query:"smokes"`
	Drinks    int    `json:"drinks" query:"drinks"`
	Ethnicity int    `json:"ethnicity" query:"ethnicity"`
	Search    []int  `json:"search" query:"search"`
	Prefer    int    `json:"prefer" query:"prefer"`
	Income    int    `json:"income" query:"income"`
	Children  int    `json:"children" query:"children"`
	Industry  int    `json:"industry" query:"industry"`
	Online    bool   `json:"online" query:"online"`
	Premium   int64  `json:"premium" query:"premium"`
}

// Набор действий открытого профиля в отношении открывающего
type Target struct {
	Like    bson.M   // Лайкнул или нет
	Private []bson.M // Дал доступ к приватным фото или нет
	Access  []bson.M // Дал доступ к переписке
}

// Структура для регистрации
type Signin struct {
	ID     string `json:"id" query:"id"`
	Token  string `json:"token" query:"token"`
	Email  string `json:"email" query:"email"`
	Method string `json:"method" query:"method"`
	Api    bool   `json:"api" query:"api"`
}

// Набор данных для перевода
type Translate struct {
	Text string `json:"text" query:"text"`
	Lang string `json:"lang" query:"lang"`
}
