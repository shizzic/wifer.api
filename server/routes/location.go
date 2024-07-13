package routes

import (
	"net/http"
	"wifer/server/crud/get"

	decoder "github.com/jesse0michael/go-request"
)

func location(props *Props) {
	props.R.Get("/country", func(w http.ResponseWriter, r *http.Request) {
		countries := get.Countries(props)
		render.JSON(w, http.StatusOK, countries)
	})

	props.R.Get("/city", func(w http.ResponseWriter, r *http.Request) {
		var data User
		decoder.Decode(r, &data)

		cities := get.Cities(props, data.Country)
		render.JSON(w, http.StatusOK, cities)
	})
}
