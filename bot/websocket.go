package bot

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/fxamacker/cbor"
	"github.com/gorilla/websocket"
	"golang.org/x/exp/slog"
)

func Websocket() error {
	wsURL := "wss://bsky.network/xrpc/com.atproto.sync.subscribeRepos"
	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		slog.Error("Failed to connect to WebSocket", "error", err)
		return err
	}
	defer conn.Close()

	slog.Info("Connected to WebSocket", "url", wsURL)

	f, _ := os.OpenFile("logs.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)

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
			var evt RepoCommitEvent
			err := decoder.Decode(&evt)
			if err == io.EOF {
				break
			}
			if err != nil {
				slog.Error("Error decoding CBOR message", "error", err)
				break
			}

			for _, op := range evt.Ops {
				f.WriteString(fmt.Sprintf("%v\n %v\n %v\n", op.Action, string(op.Text), op.Path))
				if op.Action == "create" {
					if len(evt.Blocks) > 0 {
						err := handleCARBlocks(evt.Blocks, op)
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
