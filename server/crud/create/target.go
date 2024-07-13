package create

import (
	"time"
	"wifer/server/structs"

	"go.mongodb.org/mongo-driver/bson"
)

type Props = structs.Props

// Добавить просмотр профиля, если профиль не пренадлежит пользователю
func ProfileView(props *Props, id, target int) {
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
