package db

type ToolData struct {
	Type    string   `json:"type" bson:"type"`
	Project string   `json:"project" bson:"project"`
	Data    BskyConn `json:"data" bson:"data"`
}

type BskyConn struct {
	Handle   string `json:"handle" bson:"handle"`
	Password string `json:"password" bson:"password"`
}

type SimpleProject struct {
	Name         string `json:"name" bson:"name"`
	SystemPrompt string `json:"system_prompt" bson:"system_prompt"`
}
