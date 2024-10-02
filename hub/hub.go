package hub

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Hub struct {
	clients      map[*Client]bool  // Connected clients
	topics       map[string]*Topic // Topics with subscribed clients
	register     chan *Client      // Channel to register clients
	unregister   chan *Client      // Channel to unregister clients
	broadcast    chan *Message     // Channel to broadcast messages
	sync.RWMutex                   // RWMutex for thread safety
}

// NewHub creates a new Hub instance.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		topics:     make(map[string]*Topic),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

// Run starts the main loop that listens for registration, unregistration, and broadcast requests.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.Lock()
			h.clients[client] = true
			h.Unlock()

		case client := <-h.unregister:
			h.Lock()
			if _, ok := h.clients[client]; ok {
				// Unregister client from the hub
				delete(h.clients, client)
				close(client.send)
			}
			h.Unlock()

		case message := <-h.broadcast:
			h.RLock()
			if topic, ok := h.topics[message.Topic]; ok {
				// Broadcast the message to all clients subscribed to the topic
				topic.Broadcast(message)
			}
			h.RUnlock()
		}
	}
}

// HandleClientConnection registers a new client and starts the read and write goroutines for it.
func (h *Hub) HandleClientConnection(conn *websocket.Conn) {
	client := &Client{
		conn:   conn,
		send:   make(chan []byte, 256),
		topics: make(map[string]bool),
	}

	h.register <- client

	// Start goroutines to handle sending and receiving
	go client.writePump()
	go client.readPump(h)
}

// CreateOrGetTopic retrieves an existing topic or creates a new one if it doesn't exist.
func (h *Hub) CreateOrGetTopic(topicName string) *Topic {
	h.Lock()
	defer h.Unlock()
	if topic, ok := h.topics[topicName]; ok {
		return topic
	}

	newTopic := &Topic{
		clients: make(map[*Client]bool),
	}
	h.topics[topicName] = newTopic
	return newTopic
}
