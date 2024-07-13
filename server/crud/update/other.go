package update

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Увеличить кол-во посещений приложения, за каждый клик по ссылке
func Visit(DB map[string]*mongo.Collection, ctx context.Context) {
	var data bson.M
	opts := options.FindOne().SetProjection(bson.M{"count": 1})
	DB["visits"].FindOne(ctx, bson.M{"_id": 1}, opts).Decode(&data)

	DB["visits"].UpdateOne(ctx, bson.M{"_id": 1}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "count", Value: data["count"].(int32) + 1}}},
	})
}

// Обнулить онлайн всем пользователям
func ResetOnlineForUsers(props Props) {
	props.DB["users"].UpdateMany(props.Ctx, bson.M{"online": true},
		bson.D{{Key: "$set", Value: bson.D{{Key: "online", Value: false}}}},
	)
}
