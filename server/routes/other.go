package routes

import (
	"net/http"
	"wifer/server/crud/get"
	"wifer/server/crud/update"
)

func other(props Props) {
	props.R.Get("/count", func(w http.ResponseWriter, r *http.Request) {
		quantity := get.CountAll(props.DB, props.Ctx)
		render.JSON(w, http.StatusOK, quantity)
	})

	props.R.Post("/visit", func(w http.ResponseWriter, r *http.Request) {
		update.Visit(props.DB, props.Ctx)
	})
}
