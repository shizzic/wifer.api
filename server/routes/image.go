package routes

import (
	"errors"
	"net/http"
	"os"
	im "wifer/server/image"

	decoder "github.com/jesse0michael/go-request"
)

func image(props *Props) {
	props.R.Get("/file", func(w http.ResponseWriter, r *http.Request) {
		var data Images
		decoder.Decode(r, &data)
		path, err := im.GetFilePath(props, r, &data)
		if err != nil {
			render.JSON(w, http.StatusUnauthorized, map[string]string{})
			return
		}
		fileBytes, err := os.ReadFile(path)
		if err != nil {
			render.JSON(w, http.StatusUnauthorized, map[string]string{})
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Write(fileBytes)
	})

	props.R.Post("/upload-image", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(20 << 20); err != nil {
			render.JSON(w, http.StatusBadRequest, map[string]string{"error": errors.New("max_size").Error()})
			return
		}

		var data Images
		decoder.Decode(r, &data)
		im.FillStrcut(props, r, &data)
		overcount, err := im.Upload(props, r, &data)
		if err != nil {
			render.JSON(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error(), "overcount": overcount})
			return
		}

		render.JSON(w, http.StatusOK, map[string]string{})
	})
}
