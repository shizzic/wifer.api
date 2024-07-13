package get

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Countries(props *Props) []bson.M {
	var data []bson.M
	cursor, _ := props.DB["countries"].Find(props.Ctx, bson.M{})
	cursor.All(props.Ctx, &data)

	return data
}

func Cities(props *Props, country_id int) []bson.M {
	var data []bson.M
	opts := options.Find().SetProjection(bson.M{"_id": 1, "title": 1})
	cursor, _ := props.DB["cities"].Find(props.Ctx, bson.M{"country_id": country_id}, opts)
	cursor.All(props.Ctx, &data)

	return data
}
