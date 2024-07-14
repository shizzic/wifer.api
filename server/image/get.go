package image

import (
	"errors"
	"net/http"
	"strconv"
	"wifer/server/structs"

	"go.mongodb.org/mongo-driver/bson"
)

// Геттер для всех типов файлов
func GetFilePath(props *structs.Props, r *http.Request, data *Images) (string, error) {
	path := props.Conf.PATH + "/" + data.What + "/" + data.Target + "/" + data.Into + "/" + data.Filename

	if data.Into == "private" {
		if id, err := r.Cookie("id"); err == nil {
			if id.Value == data.Target {
				return path, nil
			}

			idInt, _ := strconv.Atoi(id.Value)
			user_id, _ := strconv.Atoi(data.Target)
			found_accesses, err := props.DB["private"].CountDocuments(props.Ctx, bson.M{"user": user_id, "target": idInt})

			if err == nil && found_accesses != 0 {
				return path, nil
			} else {
				return path, errors.New("401")
			}
		}
	}

	return path, nil
}
