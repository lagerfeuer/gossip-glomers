package main

import (
	"fmt"
	"log"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()

	n.Handle("generate", func(msg maelstrom.Message) error {
		body := map[string]any{
			"type": "generate_ok",
			"id":   fmt.Sprintf("%s%d", n.ID(), time.Now().UnixNano()),
		}

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
