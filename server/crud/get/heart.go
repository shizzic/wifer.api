package get

import (
	"wifer/server/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Получить список пользователей имеющих хотя бы 1 взаимодействие с юзером и наоборот
func TargetList(props *structs.Props, data *structs.Target, id int) (int, map[string][]bson.M) {
	result := make(map[string][]bson.M)
	q := -1

	if data.Which == 0 {
		q, result = viewList(props, data, id)
	}

	if data.Which == 1 {
		q, result = likeList(props, data, id)
	}

	if data.Which == 2 {
		q, result = privateList(props, data, id)
	}

	if data.Which == 3 {
		q, result = accessList(props, data, id)
	}

	return q, result
}

func viewList(props *structs.Props, data *structs.Target, id int) (int, map[string][]bson.M) {
	res := make(map[string][]bson.M)
	q := -1
	var list []bson.M
	var ids []int32
	var key string
	var targets []bson.M

	projection := bson.M{"_id": 0, "created_at": 1, "viewed": 1}
	filter := bson.M{}

	if data.Mode {
		projection["target"] = 1
		filter["user"] = id
		key = "target"
	} else {
		projection["user"] = 1
		filter["target"] = id
		key = "user"
	}

	if data.Count {
		count, err := props.DB["views"].CountDocuments(props.Ctx, filter)
		if err != nil {
			q = 0
		} else {
			q = int(count)
		}
	}

	opts1 := options.Find().SetProjection(projection).SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(data.Limit).SetSkip(data.Skip)
	cursor, _ := props.DB["views"].Find(props.Ctx, filter, opts1)
	cursor.All(props.Ctx, &targets)
	ids = retrieveTargets(targets, key)

	if data.Mode {
		props.DB["views"].UpdateMany(props.Ctx, bson.M{"user": id, "target": bson.M{"$in": ids}}, bson.D{{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}}})
	} else {
		props.DB["views"].UpdateMany(props.Ctx, bson.M{"user": bson.M{"$in": ids}, "target": id}, bson.D{{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}}})
	}

	opts2 := options.Find().SetProjection(bson.M{"username": 1, "title": 1, "age": 1, "weight": 1, "height": 1, "body": 1, "ethnicity": 1, "public": 1, "private": 1, "avatar": 1, "premium": 1, "country_id": 1, "city_id": 1, "online": 1, "is_about": 1})
	cur, _ := props.DB["users"].Find(props.Ctx, bson.M{"_id": bson.M{"$in": ids}, "status": true}, opts2)
	cur.All(props.Ctx, &list)
	res["users"] = list
	res["targets"] = targets

	return q, res
}

func likeList(props *structs.Props, data *structs.Target, id int) (int, map[string][]bson.M) {
	res := make(map[string][]bson.M)
	q := -1
	var list []bson.M
	var ids []int32
	var key string
	var targets []bson.M

	projection := bson.M{"_id": 0, "created_at": 1, "viewed": 1}
	filter := bson.M{}

	if data.Mode {
		projection["target"] = 1
		projection["text"] = 1
		filter["user"] = id
		key = "target"
	} else {
		projection["user"] = 1
		filter["target"] = id
		key = "user"
	}

	if data.Count {
		count, err := props.DB["likes"].CountDocuments(props.Ctx, filter)
		if err != nil {
			q = 0
		} else {
			q = int(count)
		}
	}

	opts1 := options.Find().SetProjection(projection).SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(data.Limit).SetSkip(data.Skip)
	cursor, _ := props.DB["likes"].Find(props.Ctx, filter, opts1)
	cursor.All(props.Ctx, &targets)
	ids = retrieveTargets(targets, key)

	if data.Mode {
		props.DB["likes"].UpdateMany(props.Ctx, bson.M{"user": id, "target": bson.M{"$in": ids}}, bson.D{{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}}})
	} else {
		props.DB["likes"].UpdateMany(props.Ctx, bson.M{"user": bson.M{"$in": ids}, "target": id}, bson.D{{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}}})
	}

	opts2 := options.Find().SetProjection(bson.M{"username": 1, "title": 1, "age": 1, "weight": 1, "height": 1, "body": 1, "ethnicity": 1, "public": 1, "private": 1, "avatar": 1, "premium": 1, "country_id": 1, "city_id": 1, "online": 1, "is_about": 1})

	cur, _ := props.DB["users"].Find(props.Ctx, bson.M{"_id": bson.M{"$in": ids}, "status": true}, opts2)
	cur.All(props.Ctx, &list)
	res["users"] = list
	res["targets"] = targets

	return q, res
}

func privateList(props *structs.Props, data *structs.Target, id int) (int, map[string][]bson.M) {
	res := make(map[string][]bson.M)
	q := -1
	var list []bson.M
	var ids []int32
	var key string
	var targets []bson.M

	projection := bson.M{"_id": 0, "created_at": 1, "viewed": 1}
	filter := bson.M{}

	if data.Mode {
		projection["target"] = 1
		filter["user"] = id
		key = "target"
	} else {
		projection["user"] = 1
		filter["target"] = id
		key = "user"
	}

	if data.Count {
		count, err := props.DB["private"].CountDocuments(props.Ctx, filter)
		if err != nil {
			q = 0
		} else {
			q = int(count)
		}
	}

	opts1 := options.Find().SetProjection(projection).SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(data.Limit).SetSkip(data.Skip)
	cursor, _ := props.DB["private"].Find(props.Ctx, filter, opts1)
	cursor.All(props.Ctx, &targets)
	ids = retrieveTargets(targets, key)

	if data.Mode {
		props.DB["private"].UpdateMany(props.Ctx, bson.M{"user": id, "target": bson.M{"$in": ids}}, bson.D{{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}}})
	} else {
		props.DB["private"].UpdateMany(props.Ctx, bson.M{"user": bson.M{"$in": ids}, "target": id}, bson.D{{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}}})
	}

	opts2 := options.Find().SetProjection(bson.M{"username": 1, "title": 1, "age": 1, "weight": 1, "height": 1, "body": 1, "ethnicity": 1, "public": 1, "private": 1, "avatar": 1, "premium": 1, "country_id": 1, "city_id": 1, "online": 1, "is_about": 1})

	cur, _ := props.DB["users"].Find(props.Ctx, bson.M{"_id": bson.M{"$in": ids}, "status": true}, opts2)
	cur.All(props.Ctx, &list)
	res["users"] = list
	res["targets"] = targets

	return q, res
}

func accessList(props *structs.Props, data *structs.Target, id int) (int, map[string][]bson.M) {
	res := make(map[string][]bson.M)
	q := -1
	var list []bson.M
	var ids []int32
	var key string
	var targets []bson.M

	projection := bson.M{"_id": 0, "created_at": 1, "viewed": 1}
	filter := bson.M{}

	if data.Mode {
		projection["target"] = 1
		filter["user"] = id
		key = "target"
	} else {
		projection["user"] = 1
		filter["target"] = id
		key = "user"
	}

	if data.Count {
		count, err := props.DB["access"].CountDocuments(props.Ctx, filter)
		if err != nil {
			q = 0
		} else {
			q = int(count)
		}
	}

	opts1 := options.Find().SetProjection(projection).SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(data.Limit).SetSkip(data.Skip)
	cursor, _ := props.DB["access"].Find(props.Ctx, filter, opts1)
	cursor.All(props.Ctx, &targets)
	ids = retrieveTargets(targets, key)

	if data.Mode {
		props.DB["access"].UpdateMany(props.Ctx, bson.M{"user": id, "target": bson.M{"$in": ids}}, bson.D{{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}}})
	} else {
		props.DB["access"].UpdateMany(props.Ctx, bson.M{"user": bson.M{"$in": ids}, "target": id}, bson.D{{Key: "$set", Value: bson.D{{Key: "viewed", Value: true}}}})
	}

	opts2 := options.Find().SetProjection(bson.M{"username": 1, "title": 1, "age": 1, "weight": 1, "height": 1, "body": 1, "ethnicity": 1, "public": 1, "private": 1, "avatar": 1, "premium": 1, "country_id": 1, "city_id": 1, "online": 1, "is_about": 1})
	cur, _ := props.DB["users"].Find(props.Ctx, bson.M{"_id": bson.M{"$in": ids}, "status": true}, opts2)
	cur.All(props.Ctx, &list)

	res["users"] = list
	res["targets"] = targets

	return q, res
}

// Получаю массив id таргетов
func retrieveTargets(data []bson.M, key string) []int32 {
	res := []int32{}
	for _, value := range data {
		res = append(res, value[key].(int32))
	}
	return res
}
