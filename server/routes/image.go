package routes

import (
	"errors"
	"net/http"
	"os"
	im "wifer/server/image"
	"wifer/server/middlewares"

	"github.com/go-chi/chi/v5"
	decoder "github.com/jesse0michael/go-request"
)

func image(props *Props) {
	props.R.Group(func(r chi.Router) {
		// получение файлов
		r.Get("/file", func(w http.ResponseWriter, r *http.Request) {
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
	})

	props.R.Group(func(r chi.Router) {
		r.Use(middlewares.Auth(props))

		r.Post("/upload-image", func(w http.ResponseWriter, r *http.Request) {
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

		r.Put("/changeImageDir", func(w http.ResponseWriter, r *http.Request) {
			var data Images
			im.FillStrcut(props, r, &data)
			decoder.Decode(r, &data)

			if err := im.ChangeDir(props, &data); err != nil {
				render.JSON(w, http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
				return
			}

			render.JSON(w, http.StatusOK, map[string]interface{}{"message": "changed"})
		})

		r.Put("/replaceAvatar", func(w http.ResponseWriter, r *http.Request) {
			var data Images
			im.FillStrcut(props, r, &data)
			decoder.Decode(r, &data)

			im.ReplaceAvatar(props, &data)
			render.JSON(w, http.StatusOK, map[string]interface{}{"message": "replaced"})
		})

		r.Delete("/deleteImage", func(w http.ResponseWriter, r *http.Request) {
			var data Images
			im.FillStrcut(props, r, &data)
			decoder.Decode(r, &data)
			im.DeleteImage(props, &data)
			render.JSON(w, http.StatusBadRequest, map[string]interface{}{"message": "deleted"})
		})
	})
}
