package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"wifer/server/crud/get"
	"wifer/server/middlewares"
	"wifer/server/structs"

	"github.com/go-chi/chi/v5"
)

func sse(props *Props) {
	props.R.Group(func(r chi.Router) {
		r.Use(middlewares.Auth(props))

		r.Get("/sse", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", props.Conf.CLIENT_DOMAIN)
			w.Header().Set("Access-Control-Expose-Headers", "Content-Type")
			w.Header().Set("Content-Type", "text/event-stream")
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set("Connection", "keep-alive")
			id := get.UserID(w, r, props)

			// отправляю ивенты перед циклом, для как можно более быстрой актуализации
			sendNotifications(props, w, id)
			sendPremium(props, w, r)

			// открываю тикеры и безопасно закрываю их
			notifications := time.NewTicker(10 * time.Second)
			defer notifications.Stop()
			prem := time.NewTicker(60 * time.Second)
			defer prem.Stop()

			for {
				select {
				case <-r.Context().Done():
					return // закрываю цикл если юзер ливнул
				case <-notifications.C:
					sendNotifications(props, w, id)
				case <-prem.C:
					sendPremium(props, w, r)
				}

				time.Sleep(time.Second * 1)
			}
		})
	})
}

func sendNotifications(props *structs.Props, w http.ResponseWriter, id int) {
	notifications := get.Notifications(props, id)
	encode, _ := json.Marshal(notifications)
	fmt.Fprintf(w, "event: notifications\ndata: %v\n\n", string(encode))

	// Flush the data immediately instead of buffering it for later
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}

func sendPremium(props *structs.Props, w http.ResponseWriter, r *http.Request) {
	_, premium := get.UserMainInfo(props, w, r)
	fmt.Fprintf(w, "event: premium\ndata: %v\n\n", premium)

	// Flush the data immediately instead of buffering it for later
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}
