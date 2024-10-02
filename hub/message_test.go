package hub

import (
	"encoding/json"
	"testing"
)

// Test that ToJSON correctly serializes a Message into JSON
func TestMessageToJSON(t *testing.T) {
	// Create a new message
	message := NewMessage("test-topic", map[string]interface{}{
		"user": "Alice",
		"msg":  "Hello, World!",
	})

	// Convert the message to JSON
	jsonMessage := message.ToJSON()

	if jsonMessage == nil {
		t.Fatalf("expected JSON message, got nil")
	}

	// Unmarshal the JSON back into a map to verify the structure
	var jsonData map[string]interface{}
	err := json.Unmarshal(jsonMessage, &jsonData)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON message: %v", err)
	}

	// Check the topic
	if jsonData["topic"] != "test-topic" {
		t.Errorf("expected topic 'test-topic', got %s", jsonData["topic"])
	}

	// Check the data fields
	data := jsonData["data"].(map[string]interface{})
	if data["user"] != "Alice" {
		t.Errorf("expected user 'Alice', got %s", data["user"])
	}

	if data["msg"] != "Hello, World!" {
		t.Errorf("expected message 'Hello, World!', got %s", data["msg"])
	}
}

// Test that FromJSON correctly deserializes a JSON string into a Message
func TestMessageFromJSON(t *testing.T) {
	// Create a JSON byte slice
	jsonString := `{"topic": "test-topic", "data": {"user": "Alice", "msg": "Hello, World!"}}`
	jsonMessage := []byte(jsonString)

	// Convert the JSON byte slice into a Message
	message, err := FromJSON(jsonMessage)
	if err != nil {
		t.Fatalf("failed to unmarshal JSON message: %v", err)
	}

	// Check the topic
	if message.Topic != "test-topic" {
		t.Errorf("expected topic 'test-topic', got %s", message.Topic)
	}

	// Check the data fields
	data := message.Data.(map[string]interface{})
	if data["user"] != "Alice" {
		t.Errorf("expected user 'Alice', got %s", data["user"])
	}

	if data["msg"] != "Hello, World!" {
		t.Errorf("expected message 'Hello, World!', got %s", data["msg"])
	}
}
