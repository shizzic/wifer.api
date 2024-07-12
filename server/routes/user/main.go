package user

import (
	"net/http"
	"wifer/server/middlewares"
	"wifer/server/structs"

	"github.com/go-chi/chi/v5"
)

type Props = structs.Props

func Declare(props Props) {
	props.R.Group(func(r chi.Router) {
		props.R.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			// var data user
			// c.Bind(&data)

			// target := GetTarget(data.ID, c)
			// if user, err := GetProfile(data.ID); err != nil {
			// 	c.JSON(404, gin.H{"error": err.Error()})
			// } else {
			// 	c.JSON(200, gin.H{"user": user, "target": target})
			// }
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
			// var data auth
			// c.Bind(&data)

			// if err := CheckCode(data.ID, data.Code, c); err != nil {
			// 	c.JSON(401, gin.H{"error": err.Error()})
			// } else {
			// 	c.JSON(200, gin.H{"id": data.ID})
			// }
		})
	})

	props.R.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(props))
		props.R.Put("/logout", func(w http.ResponseWriter, r *http.Request) {
			// Logout(c)
		})
	})
}
