package hub

import "sync"

type Topic struct {
	clients map[*Client]bool
	sync.RWMutex
}

// Broadcast sends the JSON-encoded message to all clients subscribed to the topic.
func (t *Topic) Broadcast(message *Message) {
	jsonMessage := message.ToJSON()
	if jsonMessage == nil {
		return
	}

	t.RLock()
	defer t.RUnlock()
	for client := range t.clients {
		select {
		case client.send <- jsonMessage:
		default:
			close(client.send)
			delete(t.clients, client)
		}
	}
}
