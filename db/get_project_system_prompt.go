package db

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func GetProjectSystemPrompt(project string) (string, error) {
	client, ctx := GetMongoClient()
	databaseName := "moral_panic_bot"

	collection := client.Database(databaseName).Collection("mmosh-app-project")

	var result SimpleProject

	err := collection.FindOne(*ctx, bson.D{{Key: "key", Value: project}}).Decode(&result)

	if err != nil {
		log.Printf("Error trying to get system prompt: %v\n", err)
		return "", err
	}

	return result.SystemPrompt, nil
}
