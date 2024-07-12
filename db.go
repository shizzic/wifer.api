package main

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongo_clientent = connect_to_db()
	DB              = map[string]*mongo.Collection{
		"users":     mongo_clientent.Database("db").Collection("users"),
		"ensure":    mongo_clientent.Database("db").Collection("ensure"),
		"countries": mongo_clientent.Database("db").Collection("countries"),
		"cities":    mongo_clientent.Database("db").Collection("cities"),
		"templates": mongo_clientent.Database("db").Collection("templates"),
		"views":     mongo_clientent.Database("db").Collection("views"),
		"likes":     mongo_clientent.Database("db").Collection("likes"),
		"private":   mongo_clientent.Database("db").Collection("private"),
		"access":    mongo_clientent.Database("db").Collection("access"),
		"messages":  mongo_clientent.Database("db").Collection("messages"),
		"visits":    mongo_clientent.Database("db").Collection("visits"),
		"payments":  mongo_clientent.Database("db").Collection("payments"),
	}
)

func connect_to_db() *mongo.Client {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.MONGO_CONNECTION_STRING))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	return client
}
