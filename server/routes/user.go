package routes

import (
	"fmt"
	"net/http"
	"strings"
	"wifer/server/auth"
	"wifer/server/crud/create"
	"wifer/server/crud/get"
	"wifer/server/crud/update"
	"wifer/server/middlewares"
	"wifer/server/structs"

	"github.com/go-chi/chi/v5"
	decoder "github.com/jesse0michael/go-request"
)

func user(props *Props) {
	props.R.Group(func(r chi.Router) {
		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			var data User
			decoder.Decode(r, &data)

			target := get.TargetProfileActions(data.ID, w, r, props)
			if user, err := get.Profile(data.ID, props); err != nil {
				render.JSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
			} else {
				render.JSON(w, http.StatusOK, map[string]any{"user": user, "target": target})
			}
		})

		r.Post("/getUsers", func(w http.ResponseWriter, r *http.Request) {
			var data structs.Template
			decoder.Decode(r, &data)
			filter := get.PrepareFilter(&data)

			if data.Count {
				render.JSON(w, http.StatusOK, map[string]interface{}{
					"users": get.UsersByFilter(props, &data, &filter),
					"count": get.CountUsersByFilter(props, &filter),
				})
			} else {
				render.JSON(w, http.StatusOK, map[string]interface{}{"users": get.CountUsersByFilter(props, &filter)})
			}
		})

		r.Get("/checkUsername", func(w http.ResponseWriter, r *http.Request) {
			username := strings.TrimSpace(r.URL.Query().Get("username"))
			fmt.Print(username)
			result := get.CheckUsernameAvailability(props, username)
			render.JSON(w, http.StatusOK, result)
		})

		r.Post("/signin", func(w http.ResponseWriter, r *http.Request) {
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

		r.Post("/checkCode", func(w http.ResponseWriter, r *http.Request) {
			var data Auth
			decoder.Decode(r, &data)

			fmt.Print(data)

			if err := auth.CheckCode(props, data.ID, data.Code, w); err != nil {
				render.JSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
			} else {
				render.JSON(w, http.StatusOK, map[string]int{"id": data.ID})
			}
		})
	})

	props.R.Group(func(r chi.Router) {
		r.Use(middlewares.Auth(props))

		r.Get("/online", func(w http.ResponseWriter, r *http.Request) {
			var data User
			decoder.Decode(r, &data)
			id := get.UserID(w, r, props)

			update.ChangeLastOnline(props, data.Online, id)
		})

		r.Get("/getParamsAfterLogin", func(w http.ResponseWriter, r *http.Request) {
			user, messages := get.UserAndMessagedHimIds(props, w, r)
			render.JSON(w, http.StatusOK, map[string]interface{}{
				"user":     user,
				"messages": messages,
			})
		})

		r.Put("/change", func(w http.ResponseWriter, r *http.Request) {
			var data User
			decoder.Decode(r, &data)
			id := get.UserID(w, r, props)

			if err := update.Change(props, r, w, &data, id); err != nil {
				render.JSON(w, http.StatusBadRequest, map[string]string{"err": err.Error()})
			} else {
				render.JSON(w, http.StatusOK, map[string]string{})
			}
		})

		r.Put("/logout", func(w http.ResponseWriter, r *http.Request) {
			id := get.UserID(w, r, props)
			update.Logout(w, r, props, id)
		})

		r.Put("/deactivate", func(w http.ResponseWriter, r *http.Request) {
			id := get.UserID(w, r, props)
			update.DeactivateAccount(w, r, props, id)
		})
	})
}
