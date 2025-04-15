package db

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func SavePostAsHistoryInDb(project, text string) {
	client, ctx := GetMongoClient()
	databaseName := "moral_panic_bot"

	collection := client.Database(databaseName).Collection("mmosh-app-project-tools")

	update := bson.D{{Key: "$push", Value: bson.D{{Key: "messages", Value: map[string]string{
		"text":      text,
		"createdAt": time.Now().String(),
	}}}}}

	filter := bson.D{{Key: "project", Value: project}}

	collection.UpdateOne(*ctx, filter, update)
}
