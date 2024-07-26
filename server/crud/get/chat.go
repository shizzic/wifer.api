package get

import (
	"wifer/server/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ChatRooms(props *structs.Props, data *structs.Rooms, id int) (map[string][]bson.M, []int) {
	var res = make(map[string][]bson.M)
	var rooms []bson.M
	var ids = []int{}
	var users []bson.M

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

		for _, v := range rooms {
			user := int(v["user"].(int32))

			if user != id {
				ids = append(ids, user)
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
		nin := data.Nin
		nin = append(nin, id)
		opts := options.Find().SetProjection(bson.M{"username": 1, "avatar": 1, "online": 1})
		cur, _ := props.DB["users"].Find(props.Ctx, bson.M{"username": bson.M{"$regex": data.Username, "$options": "gi"}, "_id": bson.M{"$nin": nin}, "status": true}, opts)
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

	res["rooms"] = rooms
	res["users"] = users

	return res, ids
}
