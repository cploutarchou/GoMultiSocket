package hub

import (
	"encoding/json"
	"log"
)

// Message represents a generic message structure that can be serialized into JSON.
type Message struct {
	Topic string      `json:"topic"` // The topic name
	Data  interface{} `json:"data"`  // The data payload (can be any type)
}

// NewMessage creates a new Message instance.
func NewMessage(topic string, data interface{}) *Message {
	return &Message{
		Topic: topic,
		Data:  data,
	}
}

// ToJSON converts the message to JSON format.
func (m *Message) ToJSON() []byte {
	jsonMessage, err := json.Marshal(m)
	if err != nil {
		log.Printf("Error marshaling message to JSON: %v", err)
		return nil
	}
	return jsonMessage
}

// FromJSON creates a Message from a JSON byte slice.
func FromJSON(jsonMessage []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(jsonMessage, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}
