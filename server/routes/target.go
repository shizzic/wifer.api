package routes

import (
	"net/http"
	"wifer/server/crud/create"
	"wifer/server/crud/delete"
	"wifer/server/crud/get"
	"wifer/server/middlewares"
	"wifer/server/structs"

	"github.com/go-chi/chi/v5"
	decoder "github.com/jesse0michael/go-request"
)

func target(props *Props) {
	props.R.Group(func(r chi.Router) {
		r.Use(middlewares.Auth(props))

		r.Post("/targets", func(w http.ResponseWriter, r *http.Request) {
			var data structs.Target
			decoder.Decode(r, &data)
			id := get.UserID(w, r, props)

			count, result := get.TargetList(props, &data, id)
			render.JSON(w, http.StatusOK, map[string]interface{}{
				"data":  result,
				"count": count,
			})
		})

		r.Post("/like", func(w http.ResponseWriter, r *http.Request) {
			var data structs.Target
			decoder.Decode(r, &data)
			id := get.UserID(w, r, props)

			create.TargetLike(props, &data, id)
		})

		r.Post("/private", func(w http.ResponseWriter, r *http.Request) {
			var data structs.Target
			decoder.Decode(r, &data)
			id := get.UserID(w, r, props)

			create.TargetPrivate(props, &data, id)
		})

		r.Post("/access", func(w http.ResponseWriter, r *http.Request) {
			var data structs.Target
			decoder.Decode(r, &data)
			id := get.UserID(w, r, props)

			create.TargetAccess(props, &data, id)
		})

		r.Delete("/like", func(w http.ResponseWriter, r *http.Request) {
			var data structs.Target
			decoder.Decode(r, &data)
			id := get.UserID(w, r, props)

			delete.TargetLike(props, &data, id)
		})

		r.Delete("/private", func(w http.ResponseWriter, r *http.Request) {
			var data structs.Target
			decoder.Decode(r, &data)
			id := get.UserID(w, r, props)

			delete.TargetPrivate(props, &data, id)
		})

		r.Delete("/access", func(w http.ResponseWriter, r *http.Request) {
			var data structs.Target
			decoder.Decode(r, &data)
			id := get.UserID(w, r, props)

			delete.TargetAccess(props, &data, id)
		})
	})
}
