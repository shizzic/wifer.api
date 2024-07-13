package routes

import (
	"net/http"
	"wifer/server/crud/get"
	"wifer/server/crud/update"
	"wifer/server/middlewares"

	"github.com/go-chi/chi/v5"
)

func other(props Props) {
	props.R.Get("/count", func(w http.ResponseWriter, r *http.Request) {
		quantity := get.CountAll(props.DB, props.Ctx)
		render.JSON(w, http.StatusOK, quantity)
	})

	props.R.Post("/visit", func(w http.ResponseWriter, r *http.Request) {
		update.Visit(props.DB, props.Ctx)
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
