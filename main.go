package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/fxamacker/cbor/v2"
	"github.com/go-co-op/gocron/v2"
	"github.com/gorilla/websocket"
	"github.com/ipfs/go-cid"
	carv2 "github.com/ipld/go-car/v2"
	"github.com/joho/godotenv"
	"github.com/mmosh-pit/kinship-bsky-bot/bot"
	"github.com/mmosh-pit/kinship-bsky-bot/db"
)

const BOT_NAME = "@kinshipbot.bsky.social"

var botsHandle = map[string]db.ToolData{}

func main() {
	slog.Info("Starting bot")
	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file")
	}

	db.InitializeMongoConnection()

	go updatingBotsCron()

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

						bot.HandleBotTagged(resultingText, post, block.Cid().String(), op.Path, "", "", "", "", "")

						continue
					}

					for _, value := range botsHandle {

						botName := fmt.Sprintf("@%s", value.Data.Handle)

						if strings.Contains(text, botName) {
							resultingText := strings.ReplaceAll(text, botName, "")

							log.Printf("Gonna send with bot: %v\n", value)

							post.DID = did

							bot.HandleBotTagged(resultingText, post, block.Cid().String(), op.Path, value.Project, value.Data.Handle, value.Data.Password, value.Data.Instructions, value.FounderUsername)

							break
						}
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

func updatingBotsCron() {
	s, err := gocron.NewScheduler()
	if err != nil {
		log.Fatalf("Could not create the cron Scheduler: %v\n", err)
	}

	// add a job to the scheduler
	_, err = s.NewJob(
		gocron.DurationJob(
			5*time.Minute,
		),
		gocron.NewTask(
			func() {
				log.Println("Executing...")

				accounts, err := db.GetConnectedBskyAccounts()

				if err != nil {
					log.Printf("Got error trying to retreive accounts: %v\n", err)
					return
				}

				updatedHandles := map[string]db.ToolData{}

				for _, value := range accounts {
					updatedHandles[value.Data.Handle] = value
				}

				botsHandle = updatedHandles
				log.Printf("Updated handles: %v\n", botsHandle)
			},
		),
	)

	if err != nil {

		log.Fatalf("Could not start the cron Scheduler: %v\n", err)
	}
	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	select {}

	// when you're done, shut it down
	// log.Println("Shutting down???")
	// err = s.Shutdown()
	// if err != nil {
	// 	log.Fatalf("Could not shutdown the cron Scheduler: %v\n", err)
	// }
}
