package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"

	"github.com/fxamacker/cbor/v2"
	"github.com/gorilla/websocket"
	"github.com/ipfs/go-cid"
	carv2 "github.com/ipld/go-car/v2"
	"github.com/joho/godotenv"
	"github.com/mmosh-pit/kinship-bsky-bot/bot"
)

const BOT_NAME = "@kinshipbot.bsky.social"

func main() {
	slog.Info("Starting bot")
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	err = Websocket()
	if err != nil {
		log.Fatal(err)
	}

}

func Websocket() error {
	wsURL := "wss://bsky.network/xrpc/com.atproto.sync.subscribeRepos"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		slog.Error("Failed to connect to WebSocket", "error", err)
		return err
	}
	defer conn.Close()

	slog.Info("Connected to WebSocket", "url", wsURL)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
				slog.Warn("WebSocket connection closed", "error", err)
				break
			}

			slog.Error("Error reading message from WebSocket", "error", err)

			conn.Close()
			conn, _, err = websocket.DefaultDialer.Dial(wsURL, nil)
			if err != nil {
				slog.Error("Failed to reconnect to WebSocket", "error", err)
				return err
			}
			slog.Info("Reconnected to WebSocket", "url", wsURL)
			continue
		}

		decoder := cbor.NewDecoder(bytes.NewReader(message))

		for {
			var evt bot.RepoCommitEvent
			err := decoder.Decode(&evt)
			if err == io.EOF {
				break
			}
			if err != nil {
				slog.Error("Error decoding CBOR message", "error", err)
				break
			}

			for _, op := range evt.Ops {
				if op.Action == "create" {
					if len(evt.Blocks) > 0 {
						err := handleCARBlocks(evt.Blocks, op, evt.Repo)
						if err != nil {
							slog.Error("Error handling CAR blocks", "error", err)
							return err
						}
					}
				}
			}
		}
	}
	return nil
}

func handleCARBlocks(blocks []byte, op bot.RepoOperation, did string) error {
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
					continue
				}

				if block.Cid().Equals(c) {
					var post bot.Post

          var test map[string]interface{}

					err := cbor.Unmarshal(block.RawData(), &post)

          _ = cbor.Unmarshal(block.RawData(), &test)

					if err != nil {
						continue
					}


					text := post.Text

					if strings.Contains(text, BOT_NAME) {

						resultingText := strings.ReplaceAll(text, BOT_NAME, "")

            post.DID = did

						bot.HandleBotTagged(resultingText, post, block.Cid().String(), op.Path)
					}
				}
			}
		}
	}

	return nil
}

func decodeCID(cidBytes []byte) (cid.Cid, error) {
	var c cid.Cid
	c, err := cid.Decode(string(cidBytes))
	if err != nil {
		return c, fmt.Errorf("error decoding CID: %w", err)
	}

	return c, nil
}
