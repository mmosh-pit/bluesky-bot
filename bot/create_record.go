package bot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"time"
)

type CreateRecordProps struct {
	DIDResponse *DIDResponse
	Resource    string
	URI         string
	CID         string
	Text        string
	PostId      string
  Index int
  Total int
}

type CreateRecordResponse struct {
  URI        string `json:"uri"`
  CID         string `json:"cid"`
  Commit      Commit `json:"commit"`
  ValidationStatus string `json:"validationStatus"`
}

type Commit struct {
  CID  string `json:"cid"`
  Rev  string `json:"rev"`
}

func createRecord(r *CreateRecordProps) error {

  needMoreThreads := false

  if (len(r.Text) > 300) {
    needMoreThreads = true
  }

  resultingText := r.Text

  if needMoreThreads {
    resultingText = fmt.Sprintf("🧵%v of %v. %v", r.Index, r.Total, resultingText[:285]) 
  }


	body := map[string]interface{}{
		"collection": r.Resource,
		"repo":       r.DIDResponse.DID,

		"record": map[string]interface{}{
			"createdAt": time.Now(),
			"$type":     r.Resource,
			"text":      resultingText,
			"reply": map[string]interface{}{
				// "root": refs["root"],
				"root": map[string]string{
          "uri": fmt.Sprintf("at://%s/%s", r.PostId, r.URI),
					"cid": r.CID,
				},
				// "parent": refs["parent"],
				"parent": map[string]string{
          "uri": fmt.Sprintf("at://%s/%s", r.PostId, r.URI),
					"cid": r.CID,
				},
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
    log.Printf("Got error; %v\n", string(responseBody))
		return nil
	}


  var data CreateRecordResponse

  err = json.Unmarshal(responseBody, &data)

  if err != nil {
    log.Printf("Failed to decode create record data: %s\n", err)
  }

  if needMoreThreads {
    textHelper := r.Text

    if len(textHelper) > 285 {
      textHelper = textHelper[285:] 
    }

    resource := CreateRecordProps{
      DIDResponse: r.DIDResponse,
      Resource:    r.Resource,
      URI:         r.URI,
      CID:         data.CID,
      Text:        textHelper,
      PostId:      data.URI,
      Total: r.Total,
      Index: r.Index + 1,
    }

    createRecordReplies(resource)
  }

	return nil
}

func createRecordReplies(r CreateRecordProps) error {

  time.Sleep(time.Millisecond * 300)

  needMoreThreads := false

  if (len(r.Text) > 285) {
    needMoreThreads = true
  }

  resultingText := r.Text

  if needMoreThreads {
    resultingText = fmt.Sprintf("🧵%v of %v. %v", r.Index, r.Total, resultingText[:285]) 
  } else {
    resultingText = fmt.Sprintf("🧵%v of %v. %v", r.Index, r.Total, resultingText) 
  }

	body := map[string]interface{}{
		"collection": r.Resource,
		"repo":       r.DIDResponse.DID,

		"record": map[string]interface{}{
			"createdAt": time.Now(),
			"$type":     r.Resource,
			"text":      resultingText,
			"reply": map[string]interface{}{
				// "root": refs["root"],
				"root": map[string]string{
          "uri": r.PostId,
					"cid": r.CID,
				},
				// "parent": refs["parent"],
				"parent": map[string]string{
          "uri": r.PostId,
					"cid": r.CID,
				},
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
    log.Printf("Tried creating record reply, but got error: %v\n", string(responseBody))
		return nil
	}

  var data CreateRecordResponse

  err = json.Unmarshal(responseBody, &data)

  if err != nil {
    log.Printf("Failed to decode create record data: %s\n", err)
  }

  if needMoreThreads {
    textHelper := r.Text

    if len(textHelper) > 285 {
      textHelper = textHelper[285:] 
    }


    resource := CreateRecordProps{
      DIDResponse: r.DIDResponse,
      Resource:    r.Resource,
      URI:         r.URI,
      CID:         r.CID,
      Text:        textHelper,
      PostId:      r.PostId,
      Index: r.Index + 1,
      Total: r.Total,
    }

    createRecordReplies(resource)
  }

	return nil
}
