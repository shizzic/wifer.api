package middlewares

import (
	"net/http"
	"wifer/server/auth"
	"wifer/server/structs"

	"github.com/go-chi/cors"
)

type Props = structs.Props
type Config = structs.Config

// check if user loged in
func Auth(props *Props) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if _, e := r.Cookie("auth"); e == nil {
				next.ServeHTTP(w, r)
				return
			} else {
				token, err := r.Cookie("token")
				username, e := r.Cookie("username")

				if err == nil && e == nil && auth.DecryptToken(props, token.Value, w) == username.Value {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		})
	}
}

func SetCORS(conf *structs.Config) func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   []string{conf.FRONT_END_LINK},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
}

// http to https всегда. Только для Linux
func Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r,
		"https://"+r.Host+r.URL.String(),
		http.StatusMovedPermanently)
}
