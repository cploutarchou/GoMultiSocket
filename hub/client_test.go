package hub

import (
	"encoding/json"
	"testing"
)

// Test broadcasting a message to the topic in JSON format
func TestTopicBroadcastJSON(t *testing.T) {
	// Create a topic manually
	topic := &Topic{
		clients: make(map[*Client]bool),
	}

	// Create a client and subscribe it to the topic
	client := &Client{
		send:   make(chan []byte, 256),
		topics: make(map[string]bool),
	}
	topic.clients[client] = true

	// Create a structured JSON message
	data := map[string]interface{}{
		"user": "Alice",
		"msg":  "Hello, World!",
	}
	message := NewMessage("chat", data)

	// Broadcast the message to the topic
	go topic.Broadcast(message)

	// Test that the client receives the correct JSON message
	receivedMsg := <-client.send
	var receivedJSON Message
	err := json.Unmarshal(receivedMsg, &receivedJSON)
	if err != nil {
		t.Errorf("failed to unmarshal received message: %v", err)
	}

	if receivedJSON.Topic != "chat" {
		t.Errorf("expected topic 'chat', got %s", receivedJSON.Topic)
	}

	if receivedJSON.Data.(map[string]interface{})["user"] != "Alice" {
		t.Errorf("expected user 'Alice', got %v", receivedJSON.Data.(map[string]interface{})["user"])
	}

	if receivedJSON.Data.(map[string]interface{})["msg"] != "Hello, World!" {
		t.Errorf("expected message 'Hello, World!', got %v", receivedJSON.Data.(map[string]interface{})["msg"])
	}
}
