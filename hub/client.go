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

func (client *Client) readPump(h *Hub) {
	defer func() {
		h.unregister <- client
		err := client.conn.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			return
		}
		topic := ExtractTopic(message)
		h.broadcast <- &Message{
			topic:   topic,
			content: message,
		}
	}
}

func (client *Client) writePump() {
	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				err := client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					fmt.Println(err)
					return
				}
				return
			}
			err := client.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
