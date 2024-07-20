package create

import (
	"wifer/server/structs"

	"go.mongodb.org/mongo-driver/bson"
)

func Template(props *Props, data *structs.Template, id int) {
	props.DB["templates"].DeleteOne(props.Ctx, bson.M{"_id": id})
	props.DB["templates"].InsertOne(props.Ctx, bson.D{
		{Key: "_id", Value: id},
		{Key: "data", Value: data.Text},
	})
}
