package bot

import (
	"bytes"
	"errors"
	"io"
	"log"
	"log/slog"
	"strings"

	"github.com/fxamacker/cbor/v2"
	carv2 "github.com/ipld/go-car/v2"
)

const BOT_NAME = "@kinshipbot.bsky.social"

func handleCARBlocks(blocks []byte, op RepoOperation) error {
	if len(blocks) == 0 {
		return errors.New("no blocks to process")
	}

	reader, err := carv2.NewBlockReader(bytes.NewReader(blocks))
	if err != nil {
		slog.Error("Error creating CAR block reader", "error", err)
		return err
	}

	for {
		block, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			slog.Error("Error reading CAR block", "error", err)
			break
		}

		if opTag, ok := op.CID.(cbor.Tag); ok {
			if cidBytes, ok := opTag.Content.([]byte); ok {
				c, err := decodeCID(cidBytes)
				if err != nil {
					slog.Error("Error decoding CID from bytes", "error", err)
					continue
				}

				if block.Cid().Equals(c) {
					var post Post
					err := cbor.Unmarshal(block.RawData(), &post)

					if err != nil {
						slog.Error("Error decoding CBOR block", "error", err)
						continue
					}

					text := post.Text
					log.Printf("Got a post with text: %v\n", text)

					if strings.Contains(text, BOT_NAME) {

						resultingText := strings.ReplaceAll(text, BOT_NAME, "")

						log.Printf("Found a post that contains a tag to kinship bot!!! %v\n", resultingText)

						HandleBotTagged(resultingText, post, block.Cid().String(), op.Path)
					}
				}
			}
		}
	}

	return nil
}
