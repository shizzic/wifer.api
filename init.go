package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"wifer/server/crud/update"
	"wifer/server/middlewares"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

var (
	ctx    = context.TODO()
	conf   = get_config()
	router = chi.NewRouter()
)

// init срабатывает перед main()
func init() {
	update.ResetOnlineForUsers(&props)
	setup_middlewares()
}

func get_env() {
	var err error
	switch os := runtime.GOOS; os {
	case "windows":
		err = godotenv.Load(".env.development")
	default:
		err = godotenv.Load(".env.production")
	}
	if err != nil {
		log.Fatal(err)
	}
}

func get_config() *Config {
	get_env()
	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	path, _ := filepath.Abs("./")

	return &Config{
		PATH:                    path,
		MONGO_CONNECTION_STRING: os.Getenv("MONGO_CONNECTION_STRING"),

		SERVER_IP:         os.Getenv("SERVER_IP"),
		CLIENT_DOMAIN:     os.Getenv("CLIENT_DOMAIN"),
		SELF_DOMAIN_NAME:  os.Getenv("SELF_DOMAIN_NAME"),
		ENCRYPT_CERT_FILE: os.Getenv("ENCRYPT_CERT_FILE"),
		ENCRYPT_KEY_FILE:  os.Getenv("ENCRYPT_KEY_FILE"),

		ADMIN_EMAIL: os.Getenv("ADMIN_EMAIL"),
		EMAIL: Email{
			HOST:     os.Getenv("EMAIL_HOST"),
			USERNAME: os.Getenv("EMAIL_USERNAME"),
			PASSWORD: os.Getenv("EMAIL_PASSWORD"),
			PORT:     port,
		},

		BACKBLAZE_ID:  os.Getenv("BACKBLAZE_ID"),
		BACKBLAZE_KEY: os.Getenv("BACKBLAZE_KEY"),
		PRODUCT_NAME:  os.Getenv("PRODUCT_NAME"),
	}
}

func setup_middlewares() {
	router.Use(
		// middleware.Logger,
		middleware.RedirectSlashes,
		middleware.Recoverer,
		middleware.RealIP,
		middlewares.SetCORS(conf),
	)
}

func run() {
	switch os := runtime.GOOS; os {
	case "windows":
		http.ListenAndServe(conf.SERVER_IP+":80", router)
	default:
		go http.ListenAndServe(":80", http.HandlerFunc(middlewares.Redirect))
		http.ListenAndServeTLS(":443", "cert.pem", "key.pem", router)
	}
}
