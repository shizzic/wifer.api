package get

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Count how many users there are
func CountAll(DB map[string]*mongo.Collection, ctx context.Context) int64 {
	count, err := DB["users"].CountDocuments(ctx, bson.M{"status": true})

	if err != nil {
		return 0
	} else {
		return count
	}
}

// Update visit-count for whole app (each link open)
func Visit(DB map[string]*mongo.Collection, ctx context.Context) {
	var data bson.M
	opts := options.FindOne().SetProjection(bson.M{"count": 1})
	DB["visits"].FindOne(ctx, bson.M{"_id": 1}, opts).Decode(&data)

	DB["visits"].UpdateOne(ctx, bson.M{"_id": 1}, bson.D{
		{Key: "$set", Value: bson.D{{Key: "count", Value: data["count"].(int32) + 1}}},
	})
}
