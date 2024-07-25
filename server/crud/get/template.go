package get

import (
	"wifer/server/structs"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Получаю шаблоны для поиска
func Templates(props *structs.Props, id int) (text bson.M) {
	opts := options.FindOne().SetProjection(bson.M{"data": 1})
	props.DB["templates"].FindOne(props.Ctx, bson.M{"_id": id}, opts).Decode(&text)

	return text
}
