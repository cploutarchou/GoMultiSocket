package hub

import (
	"fmt"
	"github.com/gorilla/websocket"
)

type Client struct {
	conn   *websocket.Conn
	send   chan []byte
	topics map[string]bool
}

// readPump reads messages from the WebSocket connection and processes them.
func (client *Client) readPump(h *Hub) {
	defer func() {
		// Unregister client from the hub and close the WebSocket connection.
		h.unregister <- client
		err := client.conn.Close()
		if err != nil {
			fmt.Println("Error closing connection:", err)
		}
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			// If there's an error reading the message, exit the loop.
			fmt.Println("Error reading message:", err)
			return
		}

		// Extract the topic from the message.
		msg, parseErr := FromJSON(message)
		if parseErr != nil {
			fmt.Println("Error parsing JSON:", parseErr)
			continue
		}

		// Broadcast the message to the corresponding topic.
		h.broadcast <- msg
	}
}

// writePump writes messages from the `send` channel to the WebSocket connection.
func (client *Client) writePump() {
	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				// If the channel is closed, write a close message to the WebSocket and exit.
				err := client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					fmt.Println("Error writing close message:", err)
				}
				return
			}

			// Write the message to the WebSocket connection.
			err := client.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println("Error writing message:", err)
				return
			}
		}
	}
}
