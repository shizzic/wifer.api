package routes

import (
	"net/http"
	"wifer/server/crud/create"
	"wifer/server/crud/get"
	"wifer/server/middlewares"
	"wifer/server/structs"

	"github.com/go-chi/chi/v5"
	decoder "github.com/jesse0michael/go-request"
)

func template(props *Props) {
	props.R.Group(func(r chi.Router) {
		r.Use(middlewares.Auth(props))

		props.R.Get("/templates", func(w http.ResponseWriter, r *http.Request) {
			id := get.UserID(w, r, props)
			result := get.Templates(props, id)
			render.JSON(w, http.StatusOK, result)
		})

		props.R.Post("/templates", func(w http.ResponseWriter, r *http.Request) {
			var data structs.Template
			decoder.Decode(r, &data)
			id := get.UserID(w, r, props)

			create.Template(props, &data, id)
		})
	})
}
