package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type ToolData struct {
	Type            string   `json:"type" bson:"type"`
	Project         string   `json:"project" bson:"project"`
	Data            BskyConn `json:"data" bson:"data"`
	FounderUsername string   `json:"founder_username"`
}

type BskyConn struct {
	Handle       string `json:"handle" bson:"handle"`
	Password     string `json:"password" bson:"password"`
	Instructions string `json:"instructions" bson:"instructions"`
}

type SimpleProject struct {
	Name         string `json:"name" bson:"name"`
	SystemPrompt string `json:"system_prompt" bson:"system_prompt"`
}

type Agent struct {
	Id              *primitive.ObjectID `bson:"_id" json:"id"`
	Name            string              `bson:"name" json:"name"`
	Desc            string              `bson:"desc" json:"desc"`
	Image           string              `bson:"image" json:"image"`
	Symbol          string              `bson:"symbol" json:"symbol"`
	Key             string              `bson:"key" json:"key"`
	SystemPrompt    string              `bson:"system_prompt" json:"system_prompt"`
	CreatorUsername string              `bson:"creatorUsername" json:"creatorUsername"`
	Type            string              `bson:"type" json:"type"`
}
