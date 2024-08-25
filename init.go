package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"wifer/server/middlewares"
	"wifer/server/structs"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

var (
	ctx    = context.TODO()
	conf   = get_config()
	router = chi.NewRouter()
	props  = structs.Props{
		Conf: conf,
		Ctx:  ctx,
		R:    router,
	}
)

// init срабатывает перед main()
func init() {
	// os := runtime.GOOS
	// if os != "windows" {
	// 	cron.Start(&props)
	// }

	// connect_to_db()
	// update.ResetOnlineForUsers(&props)
	// setup_middlewares()
}

func get_env() {
	if os := runtime.GOOS; os == "windows" {
		if err := godotenv.Load("dev.env"); err != nil {
			log.Fatal(err)
		}
	}
}

func get_config() *structs.Config {
	get_env()
	port, _ := strconv.Atoi(os.Getenv("EMAIL_PORT"))
	path, _ := filepath.Abs("./")

	return &structs.Config{
		PATH:                    path,
		MONGO_CONNECTION_STRING: os.Getenv("MONGO_CONNECTION_STRING"),

		FRONT_END_LINK: os.Getenv("FRONT_END_LINK"),
		SELF_DOMAIN:    os.Getenv("SELF_DOMAIN"),
		SSL_FULL_CHAIN: os.Getenv("SSL_FULL_CHAIN"),
		SSL_PRIV_KEY:   os.Getenv("SSL_PRIV_KEY"),

		ADMIN_EMAIL: os.Getenv("ADMIN_EMAIL"),
		EMAIL: structs.Email{
			HOST:     os.Getenv("EMAIL_HOST"),
			USERNAME: os.Getenv("EMAIL_USERNAME"),
			PASSWORD: os.Getenv("EMAIL_PASSWORD"),
			PORT:     port,
		},

		BACKBLAZE_ID:  os.Getenv("BACKBLAZE_ID"),
		BACKBLAZE_KEY: os.Getenv("BACKBLAZE_KEY"),
		PRODUCT_NAME:  os.Getenv("PRODUCT_NAME"),

		GOOGLE_ID:     os.Getenv("GOOGLE_ID"),
		GOOGLE_SECRET: os.Getenv("GOOGLE_SECRET"),

		YANDEX_ID:     os.Getenv("YANDEX_ID"),
		YANDEX_SECRET: os.Getenv("YANDEX_SECRET"),

		MAIL_ID:     os.Getenv("MAIL_ID"),
		MAIL_SECRET: os.Getenv("MAIL_SECRET"),

		TWITCH_ID:     os.Getenv("TWITCH_ID"),
		TWITCH_SECRET: os.Getenv("TWITCH_SECRET"),

		VK_ID:     os.Getenv("VK_ID"),
		VK_SECRET: os.Getenv("VK_SECRET"),

		DISCORD_ID:     os.Getenv("DISCORD_ID"),
		DISCORD_SECRET: os.Getenv("DISCORD_SECRET"),
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
		http.ListenAndServe("127.0.0.1:80", router)
	default:
		// go http.ListenAndServe(":80", http.HandlerFunc(middlewares.Redirect))
		http.ListenAndServeTLS(":8443", conf.SSL_FULL_CHAIN, conf.SSL_PRIV_KEY, router)
	}
}
