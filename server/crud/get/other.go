package get

import (
	"go.mongodb.org/mongo-driver/bson"
)

// Получаю кол-во новых уведомлений для каждого типа действия
func Notifications(props *Props, id int) map[string]int64 {
	data := make(map[string]int64)

	iLikes, err := props.DB["likes"].CountDocuments(props.Ctx, bson.M{"target": id, "viewed": false})
	if err == nil {
		data["likes"] = iLikes
	}

	iViews, err := props.DB["views"].CountDocuments(props.Ctx, bson.M{"target": id, "viewed": false})
	if err == nil {
		data["views"] = iViews
	}

	iPrivates, err := props.DB["private"].CountDocuments(props.Ctx, bson.M{"target": id, "viewed": false})
	if err == nil {
		data["privates"] = iPrivates
	}

	iAccesses, err := props.DB["access"].CountDocuments(props.Ctx, bson.M{"target": id, "viewed": false})
	if err == nil {
		data["accesses"] = iAccesses
	}

	return data
}

// Узнать кол-во всех пользователей
func CountAll(props *Props) int64 {
	count, err := props.DB["users"].CountDocuments(props.Ctx, bson.M{"status": true})

	if err != nil {
		return 0
	} else {
		return count
	}
}
