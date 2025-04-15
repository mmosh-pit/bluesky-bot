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

	var agents []Agent

	res, err = client.Database(databaseName).Collection("mmosh-app-project").Find(*ctx, bson.D{{}})

	for res.Next(*ctx) {
		var agent Agent

		if err := res.Decode(&agent); err != nil {
			log.Printf("Error decoding tool: %v\n", err)
			continue
		}

		agents = append(agents, agent)
	}

	for i, value := range result {
		for _, agent := range agents {
			if agent.Key == value.Project {
				result[i].FounderUsername = agent.CreatorUsername

				break
			}
		}
	}

	return result, nil
}
