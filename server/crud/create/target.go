package create

import (
	"strings"
	"time"
	"wifer/server/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Добавить просмотр профиля, если профиль не пренадлежит пользователю
func ProfileView(props *structs.Props, id, target int) {
	if id > 0 && id != target && target > 0 {
		props.DB["views"].DeleteOne(props.Ctx, bson.M{"user": id, "target": target})

		date := time.Now().Unix()
		props.DB["views"].InsertOne(props.Ctx, bson.D{
			{Key: "user", Value: id},
			{Key: "target", Value: target},
			{Key: "viewed", Value: false},
			{Key: "created_at", Value: date},
		})
	}
}

// Юзер лайкает другово пользователя
func TargetLike(props *structs.Props, data *structs.Target, id int) {
	if id > 0 && id != data.Target && data.Target > 0 {
		date := time.Now().Unix()
		viewed := false
		var like bson.M
		text := strings.TrimSpace(data.Text)

		// Если уже лайкнул, то просто обновляем текст и дату, иначе создаю
		opts := options.FindOne().SetProjection(bson.M{"_id": 0, "viewed": 1})
		if err := props.DB["likes"].FindOne(props.Ctx, bson.M{"user": id, "target": data.Target}, opts).Decode(&like); err == nil {
			viewed = like["viewed"].(bool)

			props.DB["likes"].UpdateOne(props.Ctx, bson.M{"user": id, "target": data.Target}, bson.D{
				{Key: "$set", Value: bson.D{{Key: "text", Value: text}}},
				{Key: "$set", Value: bson.D{{Key: "viewed", Value: viewed}}},
				{Key: "$set", Value: bson.D{{Key: "created_at", Value: date}}},
			})
		} else {
			props.DB["likes"].InsertOne(props.Ctx, bson.D{
				{Key: "user", Value: id},
				{Key: "target", Value: data.Target},
				{Key: "text", Value: text},
				{Key: "viewed", Value: viewed},
				{Key: "created_at", Value: date},
			})
		}
	}
}

// Юзер дает доступ на просмотр своих приваток
func TargetPrivate(props *structs.Props, data *structs.Target, id int) {
	if id > 0 && id != data.Target && data.Target > 0 {
		// удаляю старый
		props.DB["private"].DeleteOne(props.Ctx, bson.M{"user": id, "target": data.Target})

		// пересоздаю
		date := time.Now().Unix()
		props.DB["private"].InsertOne(props.Ctx, bson.D{
			{Key: "user", Value: id},
			{Key: "target", Value: data.Target},
			{Key: "viewed", Value: false},
			{Key: "created_at", Value: date},
		})
	}
}

// Юзер дает доступ на переписку с ним
func TargetAccess(props *structs.Props, data *structs.Target, id int) {
	if id > 0 && id != data.Target && data.Target > 0 {
		// удаляю старый
		props.DB["access"].DeleteOne(props.Ctx, bson.M{"user": id, "target": data.Target})

		// пересоздаю
		date := time.Now().Unix()
		props.DB["access"].InsertOne(props.Ctx, bson.D{
			{Key: "user", Value: id},
			{Key: "target", Value: data.Target},
			{Key: "viewed", Value: false},
			{Key: "created_at", Value: date},
		})
	}
}
