package hub

import (
	"encoding/json"
	"testing"
)

// Test topic creation
func TestCreateOrGetTopic(t *testing.T) {
	// Initialize the Hub
	socketHub := NewHub()
	topicName := "test-topic"

	// Test topic creation
	topic := socketHub.CreateOrGetTopic(topicName)
	if topic == nil {
		t.Error("expected topic to be created")
	}

	// Check if the topic was added to the Hub
	if len(socketHub.topics) != 1 {
		t.Errorf("expected 1 topic, got %d", len(socketHub.topics))
	}

	// Test fetching the same topic returns the same instance
	sameTopic := socketHub.CreateOrGetTopic(topicName)
	if topic != sameTopic {
		t.Error("expected the same topic instance")
	}
}

// Test broadcasting a message to the topic
func TestTopicBroadcast(t *testing.T) {
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

	// Unmarshal the received JSON message
	var receivedJSON Message
	err := json.Unmarshal(receivedMsg, &receivedJSON)
	if err != nil {
		t.Fatalf("failed to unmarshal received message: %v", err)
	}

	// Test that the message topic and data are correct
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
