package image

import (
	"net/http"
	"os"
	"strconv"
	"wifer/server/structs"

	"go.mongodb.org/mongo-driver/bson"
)

// Заполняю данные в структуру для дальнейшего использования
func FillStrcut(props *structs.Props, r *http.Request, data *Images) {
	id, _ := r.Cookie("id")
	data.ID, _ = strconv.Atoi(id.Value)
	data.StrId = id.Value
	data.Path = props.Conf.PATH + "/images/" + data.StrId
	data.Public = data.Path + "/public"
	data.Private = data.Path + "/private"
	data.Avatar = data.Path + "/avatar.webp"
	if data.Into != "" {
		if data.Into == "public" {
			data.FullPath = data.Public + "/"
		} else {
			data.FullPath = data.Private + "/"
		}
	}
	count(data)

	// создаю дириктории, если их нет
	os.MkdirAll(data.Public, os.ModePerm)
	os.MkdirAll(data.Private, os.ModePerm)
}

// Считаю кол-во фоток у пользователя
func count(data *Images) {
	public, _ := os.ReadDir(data.Public)
	private, _ := os.ReadDir(data.Private)
	public_count := int8(len(public))
	private_count := int8(len(private))
	data.CountPublic = public_count
	data.CountPrivate = private_count
	data.Count = data.CountPublic + data.CountPrivate // кол-во фоток у юзера без учета аватарки

	_, err := os.Stat(data.Avatar)
	if err == nil {
		data.IsAvatar = true
		data.Count++
	}
}

// Обновить кол-во фоток в бд для пользователя
func updateCounts(props *structs.Props, data *Images) {
	props.DB["users"].UpdateOne(props.Ctx, bson.M{"_id": data.ID}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "avatar", Value: data.IsAvatar}}},
		{Key: "$set", Value: bson.D{{Key: "public", Value: data.CountPublic}}},
		{Key: "$set", Value: bson.D{{Key: "private", Value: data.CountPrivate}}},
		{Key: "$set", Value: bson.D{{Key: "images", Value: data.Count}}},
	})
}
