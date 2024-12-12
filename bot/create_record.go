package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

type CreateRecordProps struct {
	DIDResponse *DIDResponse
	Resource    string
	URI         string
	CID         string
	Text        string
	PostId      string
}

func createRecord(r *CreateRecordProps) error {
	postId := strings.Split(r.URI, "/")
	refs, err := getReplyRefs(r.DIDResponse.DID, r.Resource, postId[len(postId)-1])

	if err != nil {
		log.Printf("error getting reply references; %v\n", err)
		return err
	}

	log.Printf("Got refs: %v\n", refs)

	body := map[string]interface{}{
		"collection": r.Resource,
		"repo":       r.DIDResponse.DID,

		"record": map[string]interface{}{
			"createdAt": time.Now(),
			"$type":     r.Resource,
			"text":      r.Text[:300],
			"reply": map[string]interface{}{
				"root": refs["root"],
				// "root": map[string]string{
				// 	"uri": r.URI,
				// 	"cid": r.CID,
				// },
				"parent": refs["parent"],
				// "parent": map[string]string{
				// 	"uri": r.URI,
				// 	"cid": r.CID,
				// },
			},
		},
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		slog.Error("Error marshalling request", "error", err, "resource", r.Resource)
		return err
	}

	url := fmt.Sprintf("%s/com.atproto.repo.createRecord", API_URL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		slog.Error("Error creating request", "error", err, "r.Resource", r.Resource)
		return nil
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", r.DIDResponse.AccessJwt))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		slog.Error("Error sending request", "error", err, "r.Resource", r.Resource)
		return nil
	}

	responseBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusOK {
		log.Printf("error: %v\n", string(responseBody))
		return nil
	}

	slog.Info("Published successfully", "resource", r.Resource)

	return nil
}
