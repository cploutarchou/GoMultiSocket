package hub

import (
	"testing"
	"time"
)

// Test hub initialization
func TestNewHub(t *testing.T) {
	socketHub := NewHub()
	if socketHub == nil {
		t.Error("hub should not be nil")
	}

	if len(socketHub.clients) != 0 {
		t.Error("initial clients map should be empty")
	}

	if len(socketHub.topics) != 0 {
		t.Error("initial topics map should be empty")
	}
}

// Test hub registration
func TestHubClientRegistration(t *testing.T) {
	socketHub := NewHub()

	client := &Client{
		send:   make(chan []byte, 256),
		topics: make(map[string]bool),
	}

	// Start the hub's Run loop in a separate goroutine
	go socketHub.Run()

	// Register the client
	socketHub.register <- client

	// Add a small delay to ensure the client registration is processed
	time.Sleep(50 * time.Millisecond)

	// Check if the client was registered successfully
	if !socketHub.clients[client] {
		t.Error("client should be registered in the hub")
	}

	// Unregister the client
	socketHub.unregister <- client

	// Add a small delay to ensure the client unregistration is processed
	time.Sleep(50 * time.Millisecond)

	// Check if the client was unregistered successfully
	if socketHub.clients[client] {
		t.Error("client should be unregistered from the hub")
	}
}
