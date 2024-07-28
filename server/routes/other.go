package routes

import (
	"net/http"
	"wifer/server/crud/get"
	"wifer/server/crud/update"
	"wifer/server/lang"
	"wifer/server/mail"
	"wifer/server/middlewares"
	"wifer/server/structs"

	"github.com/go-chi/chi/v5"
	decoder "github.com/jesse0michael/go-request"
)

type Translate = structs.Translate

func other(props *Props) {
	props.R.Group(func(r chi.Router) {
		r.Get("/count", func(w http.ResponseWriter, r *http.Request) {
			quantity := get.CountAll(props)
			render.JSON(w, http.StatusOK, quantity)
		})

		r.Post("/visit", func(w http.ResponseWriter, r *http.Request) {
			update.Visit(props)
		})

		r.Put("/translate", func(w http.ResponseWriter, r *http.Request) {
			var data Translate
			decoder.Decode(r, &data)
			text, err := lang.TranslateText(&data)

			if err != nil {
				render.JSON(w, http.StatusBadRequest, map[string]string{"error": "0"})
			} else {
				render.JSON(w, http.StatusOK, map[string]string{"text": text})
			}
		})

		// контактная форма на главной странице
		r.Post("/contact", func(w http.ResponseWriter, r *http.Request) {
			var data structs.EmailMessage
			decoder.Decode(r, &data)

			if err := mail.ContactMe(props, &data); err != nil {
				render.JSON(w, http.StatusBadRequest, err)
			} else {
				render.JSON(w, http.StatusOK, err)
			}
		})

		r.Group(func(r chi.Router) {
			r.Use(middlewares.Auth(props))

			r.Get("/notifications", func(w http.ResponseWriter, r *http.Request) {
				id := get.UserID(w, r, props)
				res := get.Notifications(props, id)
				render.JSON(w, http.StatusOK, res)
			})
		})
	})
}
