package other

import (
	"encoding/json"
	"net/http"
	"wifer/server/get"
	"wifer/server/structs"
)

type Props = structs.Props

func Declare(props Props) {
	props.R.Get("/count", func(w http.ResponseWriter, r *http.Request) {
		quantity := get.CountAll(props.DB, props.Ctx)
		w.Header().Set("Content-Type", "application/json")
		res, _ := json.Marshal(quantity)
		w.Write(res)
	})

	// props.R.Post("/visit", func(w http.ResponseWriter, r *http.Request) {
	// 	get.Visit(DB, ctx)
	// })
}
