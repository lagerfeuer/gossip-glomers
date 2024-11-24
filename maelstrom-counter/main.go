package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

const (
	key            = "counter"
	defaultTimeout = 1 * time.Second
)

type LocalStorage struct {
	mutex *sync.Mutex
	value int
}

// update
// helper function to broadcast the new counter value to all nodes
func update(n *maelstrom.Node, l *LocalStorage) error {
	body := map[string]any{
		"type":  "sync",
		"node":  n.ID(),
		"value": l.value,
	}

	for _, node := range n.NodeIDs() {
		if node == n.ID() {
			continue
		}

		go func() {
			for {
				ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
				defer cancel()
				_, error := n.SyncRPC(ctx, node, body)
				if error == nil {
					break
				}
			}
		}()
	}
	return nil
}

func main() {
	n := maelstrom.NewNode()
	kv := maelstrom.NewSeqKV(n)
	store := LocalStorage{value: 0, mutex: &sync.Mutex{}}

	n.Handle("add", func(msg maelstrom.Message) error {
		ctx := context.Background()
		var body map[string]any
		if error := json.Unmarshal(msg.Body, &body); error != nil {
			return error
		}

		delta := int(body["delta"].(float64))

		store.mutex.Lock()
		// HACK: there's nicer ways to do this than to use an endless loop
		for {
			previous, error := kv.ReadInt(ctx, key)
			if error != nil {
				previous = 0
			}

			error = kv.CompareAndSwap(ctx, key, previous, previous+delta, true)
			if error == nil {
				break
			}
		}
		go update(n, &store)
		store.mutex.Unlock()

		return n.Reply(msg, map[string]any{
			"type": "add_ok",
		})
	})

	n.Handle("read", func(msg maelstrom.Message) error {
		store.mutex.Lock()
		value, error := kv.ReadInt(context.Background(), key)
		if error != nil {
			value = 0
		}

		if store.value > value {
			value = store.value
		}
		store.mutex.Unlock()

		body := map[string]any{
			"type":  "read_ok",
			"value": value,
		}

		return n.Reply(msg, body)
	})

	n.Handle("sync", func(msg maelstrom.Message) error {
		var body map[string]any
		error := json.Unmarshal(msg.Body, &body)
		if error != nil {
			return error
		}

		value := int(body["value"].(float64))

		store.mutex.Lock()
		store.value = value
		store.mutex.Unlock()

		return n.Reply(msg, map[string]any{
			"type": "sync_ok",
		})
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
