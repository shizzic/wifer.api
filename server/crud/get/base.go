package get

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Узнать кол-во всех пользователей
func CountAll(DB map[string]*mongo.Collection, ctx context.Context) int64 {
	count, err := DB["users"].CountDocuments(ctx, bson.M{"status": true})

	if err != nil {
		return 0
	} else {
		return count
	}
}
