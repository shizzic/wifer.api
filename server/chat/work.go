package chat

import (
	"strings"
	"wifer/server/structs"

	"go.mongodb.org/mongo-driver/bson"
)

// Написать сообщение
func write(props *structs.Props, data *structs.Message) {
	trimmed := strings.TrimSpace(data.Text)
	var roommates []int

	if data.User > data.Target {
		roommates = []int{data.Target, data.User}
	} else {
		roommates = []int{data.User, data.Target}
	}

	props.DB["messages"].InsertOne(props.Ctx, bson.D{
		{Key: "roommates", Value: roommates},
		{Key: "user", Value: data.User},
		{Key: "target", Value: data.Target},
		{Key: "viewed", Value: false},
		{Key: "text", Value: trimmed},
		{Key: "created_at", Value: data.Created_at},
	})
}

// Прочитать все непрочитанные сообщения собеседника
func view(props *structs.Props, data *structs.Message) {
	props.DB["messages"].UpdateMany(props.Ctx, bson.M{"user": data.Target, "target": data.User}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}},
	})
}

// Юзер покидает соединение полностью
func quit(id int) {
	clients.Delete(id)
}
