package structs

import (
	"context"

	"github.com/go-chi/chi/v5"
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

type auth struct {
	ID   int    `form:"id"`
	Code string `form:"code"`
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
