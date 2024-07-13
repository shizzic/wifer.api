package routes

import (
	"net/http"
	"wifer/server/auth"
	"wifer/server/crud/create"
	"wifer/server/crud/get"
	"wifer/server/crud/update"
	"wifer/server/middlewares"

	"github.com/go-chi/chi/v5"
	decoder "github.com/jesse0michael/go-request"
)

func user(props *Props) {
	props.R.Group(func(r chi.Router) {
		props.R.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			var data User
			decoder.Decode(r, &data)

			target := get.TargetProfileActions(data.ID, w, r, props)
			if user, err := get.Profile(data.ID, props); err != nil {
				render.JSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			} else {
				render.JSON(w, http.StatusOK, map[string]any{"user": user, "target": target})
			}
		})

		props.R.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
			var data Signin
			decoder.Decode(r, &data)

			// если через api, тогда получаю почту и только затим впускаю/регаю
			if data.Api {
				email, err := get.UserEmailByApi(data)

				if err != nil {
					render.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				}

				data.Email = email
			}

			id, err := create.Signin(props, w, data.Email, data.Api)

			if err != nil {
				render.JSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			} else {
				render.JSON(w, http.StatusOK, map[string]int{"id": id})
			}
		})

		props.R.Post("/checkCode", func(w http.ResponseWriter, r *http.Request) {
			var data Auth
			decoder.Decode(r, &data)

			if err := auth.CheckCode(props, data.ID, data.Code, w); err != nil {
				render.JSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
			} else {
				render.JSON(w, http.StatusOK, map[string]int{"id": data.ID})
			}
		})
	})

	props.R.Group(func(r chi.Router) {
		r.Use(middlewares.Auth(props))

		props.R.Get("/online", func(w http.ResponseWriter, r *http.Request) {
			var data User
			decoder.Decode(r, &data)
			id := get.UserID(w, r, props)

			update.ChangeLastOnline(props, data.Online, id)
		})

		props.R.Get("/getParamsAfterLogin", func(w http.ResponseWriter, r *http.Request) {
			user, messages := get.UserAndMessagedHimIds(props, w, r)
			render.JSON(w, http.StatusOK, map[string]interface{}{
				"user":     user,
				"messages": messages,
			})
		})

		props.R.Put("/logout", func(w http.ResponseWriter, r *http.Request) {
			id := get.UserID(w, r, props)
			update.Logout(w, r, props, id)
		})
	})
}
