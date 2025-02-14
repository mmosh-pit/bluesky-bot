package db

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoClient *mongo.Client
var mongoContext context.Context

func InitializeMongoConnection() {
	uri := os.Getenv("MONGO_URI")

	mongoContext := context.TODO()

	client, err := mongo.Connect(mongoContext, options.Client().
		ApplyURI(uri))

	if err != nil {
		log.Fatalf("Couldn't connect to mongo: %v\n", err)
		return
	}

	mongoClient = client
	log.Printf("MongoDB successfully connnected")
}

func GetMongoClient() (*mongo.Client, *context.Context) {
	return mongoClient, &mongoContext
}

func DisconnectMongoClient() {
	mongoClient.Disconnect(mongoContext)

	mongoClient = nil
	mongoContext = nil
}
