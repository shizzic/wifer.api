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

// Поля поиска по шаблону
type Template struct {
	Limit       int64  `json:"limit"`
	Skip        int64  `json:"skip"`
	Sort        string `json:"sort"`
	AgeMin      int    `json:"ageMin"`
	AgeMax      int    `json:"ageMax"`
	ImagesMin   int    `json:"imagesMin"`
	ImagesMax   int    `json:"imagesMax"`
	HeightMin   int    `json:"heightMin"`
	HeightMax   int    `json:"heightMax"`
	WeightMin   int    `json:"weightMin"`
	WeightMax   int    `json:"weightMax"`
	ChildrenMin int    `json:"childrenMin"`
	ChildrenMax int    `json:"childrenMax"`
	Body        []int  `json:"body"`
	Sex         []int  `json:"sex"`
	Smokes      []int  `json:"smokes"`
	Drinks      []int  `json:"drinks"`
	Ethnicity   []int  `json:"ethnicity"`
	Search      []int  `json:"search"`
	Income      []int  `json:"income"`
	Industry    []int  `json:"industry"`
	Premium     []int  `json:"premium"`
	Prefer      []int  `json:"prefer"`
	Country     []int  `json:"country"`
	City        []int  `json:"city"`
	Text        string `json:"text"`
	IsAbout     bool   `json:"is_about"`
	Avatar      bool   `json:"avatar"`
	Count       bool   `json:"count"`
}

// Набор доступов-действий открытого профиля в отношении открывающего
type Actions struct {
	Target  int      // id таргета
	Like    bson.M   // Лайкнул или нет
	Private []bson.M // Дал доступ к приватным фото или нет
	Access  []bson.M // Дал доступ к переписке или нет
}

// Параметры для поиска всех кто сделал хотя бы 1 действие в отношении юзера и наоборот
type Target struct {
	Target int    `json:"target" query:"target"`
	Which  int    `json:"which" query:"which"`
	Skip   int64  `json:"skip" query:"skip"`
	Limit  int64  `json:"limit" query:"limit"`
	Mode   bool   `json:"mode" query:"mode"`
	Count  bool   `json:"count" query:"count"`
	Text   string `json:"text" query:"text"` // Это поле нужно для добавления заметок на лайк
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

// Данные для добавления фотографий
type Images struct {
	ID           int
	StrId        string
	Path         string // путь до входной папки
	Avatar       string
	IsAvatar     bool `json:"isAvatar" query:"isAvatar"`
	Public       string
	Private      string
	FullPath     string // путь без названия
	Output       string // путь с названием
	Count        int8
	CountPublic  int8
	CountPrivate int8
	Into         string `json:"dir" query:"dir"`
	NewDir       string `json:"newDir" query:"newDir"`
	Filename     string `json:"filename" query:"filename"`
	What         string `json:"what" query:"what"`
	Target       string `json:"target_id" query:"target_id"`
}

type Messages struct {
	Target int   `json:"target" query:"target"`
	Skip   int64 `json:"skip" query:"skip"`
	Access bool  `json:"access" query:"access"`
}

type Rooms struct {
	Nin        []int  `json:"nin" query:"nin"`
	Username   string `json:"username" query:"username"`
	ByUsername bool   `json:"byUsername" query:"byUsername"`
}

type Message struct {
	Api        string `json:"api"`
	Text       string `json:"text"`
	Username   string `json:"username"`
	User       int    `json:"user"`
	Target     int    `json:"target"`
	Access     bool   `json:"access"`
	Typing     bool   `json:"typing"`
	Avatar     bool   `json:"avatar"`
	Created_at int64  `json:"created_at"`
}
