package get

import (
	"wifer/server/structs"

	"go.mongodb.org/mongo-driver/bson"
)

// Получаю кол-во новых уведомлений для каждого типа действия
func Notifications(props *structs.Props, id int) map[string]int64 {
	data := make(map[string]int64)

	if iLikes, err := props.DB["likes"].CountDocuments(props.Ctx, bson.M{"target": id, "viewed": false}); err == nil {
		data["likes"] = iLikes
	}

	if iViews, err := props.DB["views"].CountDocuments(props.Ctx, bson.M{"target": id, "viewed": false}); err == nil {
		data["views"] = iViews
	}

	if iPrivates, err := props.DB["private"].CountDocuments(props.Ctx, bson.M{"target": id, "viewed": false}); err == nil {
		data["privates"] = iPrivates
	}

	if iAccesses, err := props.DB["access"].CountDocuments(props.Ctx, bson.M{"target": id, "viewed": false}); err == nil {
		data["accesses"] = iAccesses
	}

	// Получить список пользователей, которые отправили залогиненому юзеру сообщение хотя бы 1 раз (не прочитанное)
	if messagedUsers, err := props.DB["messages"].Distinct(props.Ctx, "user", bson.M{"target": id, "viewed": false}); err == nil {
		data["messages"] = int64(len(messagedUsers))
	}

	return data
}

// Узнать кол-во всех пользователей
func CountAll(props *structs.Props) int64 {
	count, err := props.DB["users"].CountDocuments(props.Ctx, bson.M{"status": true})

	if err != nil {
		return 0
	} else {
		return count
	}
}
