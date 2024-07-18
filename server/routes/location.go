package routes

import (
	"net/http"
	"wifer/server/crud/get"

	decoder "github.com/jesse0michael/go-request"
)

func location(props *Props) {
	props.R.Get("/country", func(w http.ResponseWriter, r *http.Request) {
		locale := r.URL.Query().Get("locale")

		countries := get.Countries(props, locale)
		render.JSON(w, http.StatusOK, countries)
	})

	props.R.Get("/city", func(w http.ResponseWriter, r *http.Request) {
		var data User
		decoder.Decode(r, &data)
		locale := r.URL.Query().Get("locale")

		cities := get.Cities(props, data.Country, locale)
		render.JSON(w, http.StatusOK, cities)
	})
}
