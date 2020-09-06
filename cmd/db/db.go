package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var DB = "vobe-auth"
var Client *mongo.Client
var CTX *context.Context

var cachedCollections = make(map[string]*mongo.Collection)

func MongoConnect(uri string, db string) {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	DB = db
	Client = client
	CTX = &ctx
}

func GetCollection(coll string) *mongo.Collection {
	if val, ok := cachedCollections[coll]; ok {
		return val
	}
	cachedCollections[coll] = Client.Database(DB).Collection(coll)
	return GetCollection(coll)
}
