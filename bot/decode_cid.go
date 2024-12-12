package bot

import (
	"fmt"

	"github.com/ipfs/go-cid"
)

func decodeCID(cidBytes []byte) (cid.Cid, error) {
	var c cid.Cid
	c, err := cid.Decode(string(cidBytes))
	if err != nil {
		return c, fmt.Errorf("error decoding CID: %w", err)
	}

	return c, nil
}
