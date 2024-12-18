package main

import (
	"encoding/json"
	"log"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

var (
	topology map[string]any
	messages []int
)

func getConnectedNodes(n *maelstrom.Node) []string {
	nodes := topology[n.ID()].([]any)
	result := []string{}
	for _, n := range nodes {
		result = append(result, n.(string))
	}
	return result
}

func main() {
	n := maelstrom.NewNode()

	n.Handle("topology", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		topology = body["topology"].(map[string]any)

		return n.Reply(msg, map[string]any{
			"type": "topology_ok",
		})
	})

	n.Handle("broadcast", func(msg maelstrom.Message) error {
		var body map[string]any
		if err := json.Unmarshal(msg.Body, &body); err != nil {
			return err
		}

		messages = append(messages, int(body["message"].(float64)))

		for _, dest := range getConnectedNodes(n) {
			if dest == msg.Src {
				continue
			}
			go n.Send(dest, msg.Body)
		}

		response := map[string]any{
			"type": "broadcast_ok",
		}

		return n.Reply(msg, response)
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		body := map[string]any{
			"type":     "read_ok",
			"messages": messages,
		}

		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
