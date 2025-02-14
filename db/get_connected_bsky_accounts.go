package db

import (
	"log"

	"go.mongodb.org/mongo-driver/bson"
)

func GetConnectedBskyAccounts() ([]ToolData, error) {
	client, ctx := GetMongoClient()
	databaseName := "moral_panic_bot"

	collection := client.Database(databaseName).Collection("mmosh-app-project-tools")

	var result []ToolData

	res, err := collection.Find(*ctx, bson.D{{Key: "type", Value: "bsky"}})

	if err != nil {
		return nil, err
	}

	for res.Next(*ctx) {
		var tool ToolData

		if err := res.Decode(&tool); err != nil {
			log.Printf("Error decoding tool: %v\n", err)
			continue
		}

		result = append(result, tool)
	}

	return result, nil
}
