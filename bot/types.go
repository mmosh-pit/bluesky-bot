package bot

import "time"

type RepoCommitEvent struct {
	Repo   string      `cbor:"repo"`
	Rev    string      `cbor:"rev"`
	Seq    int64       `cbor:"seq"`
	Since  string      `cbor:"since"`
	Time   string      `cbor:"time"`
	TooBig bool        `cbor:"tooBig"`
	Prev   interface{} `cbor:"prev"`
	Rebase bool        `cbor:"rebase"`
	Blocks []byte      `cbor:"blocks"`

	Ops []RepoOperation `cbor:"ops"`
}

type RepoOperation struct {
	Action string      `cbor:"action"`
	Path   string      `cbor:"path"`
	Reply  *Reply      `cbor:"reply"`
	Text   []byte      `cbor:"text"`
	CID    interface{} `cbor:"cid"`
}

type Reply struct {
	Parent Parent `json:"parent"`
	Root   Root   `json:"root"`
}

type Parent struct {
	Cid string `json:"cid"`
	Uri string `json:"uri"`
}

type Root struct {
	Cid string `json:"cid"`
	Uri string `json:"uri"`
}

type Post struct {
	Type      string          `json:"$type"`
	CreatedAt time.Time       `json:"createdAt"`
	Facets    []RichTextFacet `json:"facets"`
	Langs     []string        `json:"langs"`
	Text      string          `json:"text"`
}

type RichTextFacet struct {
	Type     string                 `json:"$type"`
	Features []RichTextFacetMention `json:"features"`
	Index    map[string]int         `json:"index"`
}

type RichTextFacetMention struct {
	Type string `json:"$type"`
	Did  string `json:"did"`
}
