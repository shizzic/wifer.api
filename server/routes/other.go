package routes

import (
	"net/http"
	"wifer/server/crud/get"
	"wifer/server/crud/update"
	"wifer/server/lang"
	"wifer/server/middlewares"
	"wifer/server/structs"

	"github.com/go-chi/chi/v5"
	decoder "github.com/jesse0michael/go-request"
)

type Translate = structs.Translate

func other(props *Props) {
	props.R.Get("/count", func(w http.ResponseWriter, r *http.Request) {
		quantity := get.CountAll(props)
		render.JSON(w, http.StatusOK, quantity)
	})

	props.R.Post("/visit", func(w http.ResponseWriter, r *http.Request) {
		update.Visit(props)
	})

	props.R.Put("/translate", func(w http.ResponseWriter, r *http.Request) {
		var data Translate
		decoder.Decode(r, &data)
		text, err := lang.TranslateText(&data)

		if err != nil {
			render.JSON(w, http.StatusBadRequest, map[string]string{"error": "0"})
		} else {
			render.JSON(w, http.StatusOK, map[string]string{"text": text})
		}
	})

	props.R.Group(func(r chi.Router) {
		r.Use(middlewares.Auth(props))

		r.Get("/notifications", func(w http.ResponseWriter, r *http.Request) {
			id := get.UserID(w, r, props)
			res := get.Notifications(props, id)
			render.JSON(w, http.StatusOK, res)
		})
	})
}
