package bot

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
)

const url = "https://mmoshapi-471939176450.us-central1.run.app/generate/"

func HandleBotTagged(text string, post Post, cid, path string) {
	client := http.Client{}

	body := map[string]any{
		"username":   "Visitor",
		"prompt":     text,
		"namespaces": []string{"PUBLIC"},
		"metafield":  "",
	}

	encoded, err := json.Marshal(body)

	if err != nil {
		log.Printf("Could not encode request body: %v\n", err)
		return
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(encoded))

	if err != nil {
		log.Printf("Could not create the POST request: %v\n", err)
		return
	}

	request.Header.Set("Content-Type", "application/json")

	res, err := client.Do(request)

	if err != nil {
		log.Printf("Error sending POST request: %v\n", err)
		return
	}

	defer res.Body.Close()

	responseBody, _ := io.ReadAll(res.Body)

	token, err := getToken()

	log.Printf("Path: %v\n CID: %v\n", path, cid)

	postId := post.Facets[len(post.Facets)-1].Features[len(post.Facets[len(post.Facets)-1].Features)-1].Did

	resource := &CreateRecordProps{
		DIDResponse: token,
		Resource:    "app.bsky.feed.post",
		URI:         path,
		CID:         cid,
		Text:        string(responseBody),
		PostId:      postId,
	}

	err = createRecord(resource)
	if err != nil {
		slog.Error("Error creating record", "error", err, "resource", resource.Resource)
	}
}
