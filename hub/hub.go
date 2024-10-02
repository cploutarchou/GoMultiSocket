package hub

import (
	"github.com/gorilla/websocket"
	"sync"
)

type Hub struct {
	clients    map[*Client]bool
	topics     map[string]*Topic
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
	sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		topics:     make(map[string]*Topic),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
		RWMutex:    sync.RWMutex{},
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			if topic, ok := h.topics[message.topic]; ok {
				topic.Broadcast(message.content)
			}
		}
	}
}
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
