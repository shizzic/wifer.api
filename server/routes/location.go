package routes

import (
	"net/http"
	"wifer/server/crud/get"

	"github.com/go-chi/chi/v5"
	decoder "github.com/jesse0michael/go-request"
)

func location(props *Props) {
	props.R.Group(func(r chi.Router) {
		r.Get("/country", func(w http.ResponseWriter, r *http.Request) {
			locale := r.URL.Query().Get("locale")

			countries := get.Countries(props, locale)
			render.JSON(w, http.StatusOK, countries)
		})

		r.Get("/city", func(w http.ResponseWriter, r *http.Request) {
			var data User
			decoder.Decode(r, &data)
			locale := r.URL.Query().Get("locale")

			cities := get.Cities(props, data.Country, locale)
			render.JSON(w, http.StatusOK, cities)
		})
	})
}
