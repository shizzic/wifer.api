package main

import (
	"log"
	"regexp"

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
	// экранирую @ символ, так как mongodb нужен он
	str_raw := conf.MONGO_CONNECTION_STRING
	mode := regexp.MustCompile("^(.*?)@(.*)$")
	replace := "${1}\\@$2"
	connect_uri := mode.ReplaceAllString(str_raw, replace)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connect_uri))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Print("CONNECTED TO MONGO\n")
	return client
}
