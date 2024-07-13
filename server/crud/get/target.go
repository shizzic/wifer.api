package get

import (
	"net/http"
	"wifer/server/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Props = structs.Props
type Target = structs.Target

// Получить все действия открываемого профиля, относящиеся к открывающему пользователю
func TargetProfileActions(target int, w http.ResponseWriter, r *http.Request, props Props) Target {
	id := UserID(w, r, props)
	var data Target

	if id > 0 && id != target && target > 0 {
		// AddView(idInt, target, c)

		if err, like := likes(id, target, props); !err {
			data.Like = like
		}

		if err, priv := accessesForImages(id, target, props); !err {
			data.Private = priv
		}

		if err, access := accessesForTexting(id, target, props); !err {
			data.Access = access
		}

		return data
	}

	return data
}

// Узнать, лайкнул ли человек, открываемый профиль
func likes(id, target int, props Props) (bool, bson.M) {
	var like bson.M
	opts := options.FindOne().SetProjection(bson.M{"_id": 0, "text": 1})

	if err := props.DB["likes"].FindOne(props.Ctx, bson.M{"user": id, "target": target}, opts).Decode(&like); err == nil {
		return false, like
	} else {
		return true, like
	}
}

// Узнать, есть ли у пользователя доступ к приватным фотографиям для профиля, который он открывает
func accessesForImages(id, target int, props Props) (bool, []bson.M) {
	arr := [2]int{}
	arr[0] = id
	arr[1] = target
	var data []bson.M
	opts := options.Find().SetProjection(bson.M{"_id": 0, "user": 1})

	if cursor, err := props.DB["private"].Find(props.Ctx, bson.M{"user": bson.M{"$in": arr}, "target": bson.M{"$in": arr}}, opts); err == nil {
		if e := cursor.All(props.Ctx, &data); e == nil {
			return false, data
		} else {
			return true, data
		}
	}

	return true, data
}

// Узнать, есть ли у пользователя доступ к написанию сообщений профилю, который он открывает
func accessesForTexting(id, target int, props Props) (bool, []bson.M) {
	arr := [2]int{}
	arr[0] = id
	arr[1] = target
	var data []bson.M
	opts := options.Find().SetProjection(bson.M{"_id": 0, "user": 1})

	if cursor, err := props.DB["access"].Find(props.Ctx, bson.M{"user": bson.M{"$in": arr}, "target": bson.M{"$in": arr}}, opts); err == nil {
		if e := cursor.All(props.Ctx, &data); e == nil {
			return false, data
		} else {
			return true, data
		}
	}

	return true, data
}
