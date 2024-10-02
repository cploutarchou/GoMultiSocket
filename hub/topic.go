package hub

import "sync"

type Topic struct {
	clients map[*Client]bool
	sync.RWMutex
}

func (t *Topic) Broadcast(content []byte) {
	t.RLock()
	defer t.RUnlock()
	for client := range t.clients {
		select {
		case client.send <- content:
		default:
			close(client.send)
			delete(t.clients, client)
		}
	}
}

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
