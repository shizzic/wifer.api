package routes

import (
	"net/http"
	"wifer/server/chat"
	"wifer/server/crud/get"
	"wifer/server/middlewares"
	"wifer/server/structs"

	"github.com/go-chi/chi/v5"
	decoder "github.com/jesse0michael/go-request"
)

func web_chat(props *Props) {
	props.R.Group(func(r chi.Router) {
		r.Use(middlewares.Auth(props))

		// Меняю http на ws/wss
		r.Get("/chat", func(w http.ResponseWriter, r *http.Request) {
			id := get.UserID(w, r, props)
			chat.Connect(w, r, id)
		})

		r.Post("/getRooms", func(w http.ResponseWriter, r *http.Request) {
			var data structs.Rooms
			decoder.Decode(r, &data)
			id := get.UserID(w, r, props)

			res, ids := get.ChatRooms(props, &data, id)
			render.JSON(w, http.StatusOK, map[string]interface{}{
				"rooms": res["rooms"],
				"users": res["users"],
				"ids":   ids,
			})
		})
	})
}
