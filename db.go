package main

import (
	"log"
	"wifer/mongorestore"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// Ставлю бд с облака, если ее не существует (или прямо с папки init_dump в корне проекта, если это локалка и эта папка существует)
	if dbs, err := client.ListDatabaseNames(ctx, bson.M{"name": primitive.Regex{Pattern: "^db$"}}); err != nil || len(dbs) == 0 {
		mongorestore.Start(&props)
	}

	props.DB = map[string]*mongo.Collection{
		"users":     client.Database("db").Collection("users"),
		"ensure":    client.Database("db").Collection("ensure"),
		"countries": client.Database("db").Collection("countries"),
		"cities":    client.Database("db").Collection("cities"),
		"templates": client.Database("db").Collection("templates"),
		"views":     client.Database("db").Collection("views"),
		"likes":     client.Database("db").Collection("likes"),
		"private":   client.Database("db").Collection("private"),
		"access":    client.Database("db").Collection("access"),
		"messages":  client.Database("db").Collection("messages"),
		"visits":    client.Database("db").Collection("visits"),
		"payments":  client.Database("db").Collection("payments"),
	}

	return client
}
