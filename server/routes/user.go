package routes

import (
	"net/http"
	"wifer/server/auth"
	"wifer/server/crud/get"
	"wifer/server/crud/update"
	"wifer/server/middlewares"

	"github.com/go-chi/chi/v5"
	decoder "github.com/jesse0michael/go-request"
)

func user(props Props) {
	props.R.Group(func(r chi.Router) {
		props.R.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			var data User
			decoder.Decode(r, &data)

			target := get.GetTarget(data.ID, w, r, props)
			if user, err := get.Profile(data.ID, props); err != nil {
				render.JSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			} else {
				render.JSON(w, http.StatusOK, map[string]any{"user": user, "target": target})
			}
		})

		props.R.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
			// var data signin
			// c.Bind(&data)

			// var err error
			// var id int

			// if data.Api {
			// 	id, err = CheckApi(data, c)
			// } else {
			// 	id, err = Signin(data.Email, c, false)
			// }

			// if err != nil {
			// 	c.JSON(400, gin.H{"error": err.Error()})
			// } else {
			// 	c.JSON(200, gin.H{"id": id})
			// }
		})

		props.R.Post("/checkCode", func(w http.ResponseWriter, r *http.Request) {
			data := &Auth{}

			if err := auth.CheckCode(props, data.ID, data.Code, w); err != nil {
				render.JSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
			} else {
				render.JSON(w, http.StatusOK, map[string]int{"id": data.ID})
			}
		})
	})

	props.R.Group(func(r chi.Router) {
		r.Use(middlewares.Auth(props))

		props.R.Put("/logout", func(w http.ResponseWriter, r *http.Request) {
			update.Logout(w, r, props)
		})
	})
}
