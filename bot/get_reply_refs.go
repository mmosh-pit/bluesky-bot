package bot

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	neturl "net/url"
)

const pdsUrl = "https://bsky.social"

func getReplyRefs(repo, resource, rkey string) (map[string]map[string]string, error) {
	parentUrl := fmt.Sprintf("%s/xrpc/com.atproto.repo.getRecord?repo=%s&collection=%s&rkey=%v", pdsUrl, repo, resource, rkey)

	log.Printf("Parent URL: %v\n", parentUrl)
	req, err := http.NewRequest("GET", parentUrl, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var parentData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&parentData)
	if err != nil {
		return nil, err
	}

	rootData := parentData

	log.Printf("Got parent data: %v\n", parentData)

	if parentData["value"] != nil {
		parentReply := parentData["value"].(map[string]interface{})["reply"]
		if parentReply != nil {
			rootUri := parentReply.(map[string]interface{})["root"].(map[string]interface{})["uri"].(string)
			rootParts := parseUri(rootUri)
			rootUrl := fmt.Sprintf("%s/xrpc/com.atproto.repo.getRecord", pdsUrl)
			req, err := http.NewRequest("GET", rootUrl, nil)
			if err != nil {
				return nil, err
			}
			q := req.URL.Query()
			for k, v := range rootParts {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()

			resp, err := client.Do(req)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			err = json.NewDecoder(resp.Body).Decode(&rootData)
			if err != nil {
				return nil, err
			}
		}
	}

	return map[string]map[string]string{
		"root": {
			"uri": rootData["uri"].(string),
			"cid": rootData["cid"].(string),
		},
		"parent": {
			"uri": parentData["uri"].(string),
			"cid": parentData["cid"].(string),
		},
	}, nil
}

func parseUri(uri string) map[string]string {
	// Implement parsing logic here, similar to Python's parse_uri function
	// You might use a URL parsing library or manually extract the necessary parts.
	// For example:
	u, err := neturl.Parse(uri)
	if err != nil {
		// Handle error
	}
	return map[string]string{
		"repo":       u.Host,
		"collection": u.Path,
		"rkey":       u.Fragment,
	}
}
