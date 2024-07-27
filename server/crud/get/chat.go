package get

import (
	"net/http"
	"wifer/server/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ChatRooms(props *structs.Props, data *structs.Rooms, r *http.Request, id int) (map[string][]bson.M, []int) {
	var rooms []bson.M
	var users []bson.M
	var ids = []int{}

	groupFilter := bson.D{{Key: "$group",
		Value: bson.D{
			{Key: "_id", Value: "$roommates"},
			{Key: "user", Value: bson.D{{Key: "$first", Value: "$user"}}},
			{Key: "target", Value: bson.D{{Key: "$first", Value: "$target"}}},
			{Key: "text", Value: bson.D{{Key: "$first", Value: "$text"}}},
			{Key: "created_at", Value: bson.D{{Key: "$first", Value: "$created_at"}}},
			{Key: "viewed", Value: bson.D{{Key: "$first", Value: "$viewed"}}},
			{Key: "news", Value: bson.D{{Key: "$sum", Value: bson.D{{
				Key: "$cond", Value: bson.D{
					{Key: "if", Value: bson.D{
						{Key: "$and", Value: []bson.D{
							{
								{Key: "$eq", Value: []interface{}{"$viewed", false}},
							},
							{
								{Key: "$eq", Value: []interface{}{"$target", id}},
							},
						}},
					}},
					{Key: "then", Value: 1},
					{Key: "else", Value: 0},
				},
			}}}}},
		},
	}}

	sortFilter := bson.D{{Key: "$sort", Value: bson.D{{Key: "created_at", Value: -1}}}}
	limitFilter := bson.D{{Key: "$limit", Value: 25}}

	if !data.ByUsername {
		matchFilter := bson.D{{Key: "$match",
			Value: bson.D{{
				Key: "roommates", Value: bson.D{
					{Key: "$in", Value: []int{id}},
					{Key: "$nin", Value: data.Nin},
				},
			}},
		}}

		cursor, _ := props.DB["messages"].Aggregate(props.Ctx, mongo.Pipeline{matchFilter, sortFilter, groupFilter, sortFilter, limitFilter})
		cursor.All(props.Ctx, &rooms)

		var _, no_premium = r.Cookie("premium")
		for _, v := range rooms {
			user := int(v["user"].(int32))
			v["access"] = true

			if user != id {
				ids = append(ids, user)

				// подменяю текст последнего сообщения, если у читателя нет доступа
				if no_premium != nil {
					var access bson.M
					opts := options.FindOne().SetProjection(bson.M{"_id": 1})
					if err := props.DB["access"].FindOne(props.Ctx, bson.M{"target": id, "user": user}, opts).Decode(&access); err != nil {
						v["text"] = ""
						v["access"] = false
					}
				}
				continue
			}

			ids = append(ids, int(v["target"].(int32)))
		}

		if len(ids) > 0 {

			opts := options.Find().SetProjection(bson.M{"username": 1, "avatar": 1, "online": 1})
			cur, _ := props.DB["users"].Find(props.Ctx, bson.M{"_id": bson.M{"$in": ids}, "status": true}, opts)
			cur.All(props.Ctx, &users)
		}
	} else {
		// Ищу всех, кроме забаненных и самого юзера (неактивных тоже)
		nin := data.Nin
		nin = append(nin, id)
		opts := options.Find().SetProjection(bson.M{"username": 1, "avatar": 1, "online": 1})
		cur, _ := props.DB["users"].Find(props.Ctx, bson.M{"username": bson.M{"$regex": data.Username, "$options": "i"}, "_id": bson.M{"$nin": nin}, "status": true}, opts)
		cur.All(props.Ctx, &users)
		for _, v := range users {
			ids = append(ids, int(v["_id"].(int32)))
		}

		if len(ids) > 0 {
			freshIds := ids
			freshIds = append(freshIds, id)

			matchFilter := bson.D{{Key: "$match",
				Value: bson.D{
					{Key: "user", Value: bson.D{{Key: "$in", Value: freshIds}}},
					{Key: "target", Value: bson.D{{Key: "$in", Value: freshIds}}},
				},
			}}

			cursor, _ := props.DB["messages"].Aggregate(props.Ctx, mongo.Pipeline{matchFilter, sortFilter, groupFilter, sortFilter, limitFilter})
			cursor.All(props.Ctx, &rooms)
		}
	}

	return map[string][]bson.M{"rooms": rooms, "users": users}, ids
}

func ChatMessages(props *structs.Props, data *structs.Messages, r *http.Request, id int) map[string][]bson.M {
	var res = make(map[string][]bson.M)
	access := true
	filter := bson.M{"user": bson.M{"$in": []int{id, data.Target}}, "target": bson.M{"$in": []int{id, data.Target}}}

	if data.Access {
		if _, err := r.Cookie("premium"); err != nil {
			access = false
		}

		accesses := checkRoomAccess(props, filter)
		res["accesses"] = accesses

		if !access {
			for _, v := range accesses {
				if v["target"].(int32) == int32(id) {
					access = true
				}
			}
		}
	}

	if access {
		var messages []bson.M
		opts := options.Find().SetProjection(bson.M{
			"_id":        0,
			"user":       1,
			"text":       1,
			"created_at": 1,
			"viewed":     1,
		}).
			SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(25).SetSkip(data.Skip)

		cursor, _ := props.DB["messages"].Find(props.Ctx, filter, opts)
		cursor.All(props.Ctx, &messages)

		res["messages"] = messages
	}

	return res
}

// Проверить есть ли у юзера доступ к переписке
func checkRoomAccess(props *structs.Props, filter primitive.M) []bson.M {
	var accesses []bson.M
	opts := options.Find().SetProjection(bson.M{"_id": 0, "user": 1, "target": 1})
	cursor, _ := props.DB["access"].Find(props.Ctx, filter, opts)
	cursor.All(props.Ctx, &accesses)
	return accesses
}

// Получить кто онлайн (по юзеру)
func OnlineInChat(props *structs.Props, data *structs.Rooms) []bson.M {
	var users []bson.M
	opts := options.Find().SetProjection(bson.M{"online": 1})
	cur, _ := props.DB["users"].Find(props.Ctx, bson.M{"_id": bson.M{"$in": data.Nin}, "status": true}, opts)
	cur.All(props.Ctx, &users)
	return users
}
