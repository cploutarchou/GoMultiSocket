# GoMultiSocket
![Go Report Card](https://goreportcard.com/badge/github.com/cploutarchou/GoMultiSocket)

A SOLID Go WebSocket implementation with support for multiple topics.

## Features

- WebSocket connection management
- Support for multiple topics (channels)
- Broadcast messages to specific topics
- Lightweight and scalable
- Easy to extend

## Installation

To install the GoMultiSocket package, run the following command:

```bash
go get github.com/cploutarchou/GoMultiSocket/hub
```

Then, import the package into your project:

```go
import "github.com/cploutarchou/GoMultiSocket/hub"
```
## Usage
Here is a simple example of how to use GoMultiSocket to create a WebSocket server that supports multiple topics:

### Example:
```go

package main

import (
	"github.com/cploutarchou/GoMultiSocket/hub"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// Upgrader to convert HTTP connections to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Serve WebSocket connections and register them to the hub
func serveWs(hub *hub.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	hub.HandleClientConnection(conn)
}

func main() {
	// Create a new hub instance and start its event loop
	hub := hub.NewHub()
	go hub.Run()

	// HTTP handler to handle WebSocket requests
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	// Start the WebSocket server
	log.Println("WebSocket server started on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

```

## Sending and Receiving Messages
When a client connects to the WebSocket, they can send and receive JSON-encoded messages with the following structure:

### Message Format:

```json
{
	"topic": "example-topic",
	"data": {
		"user": "Alice",
		"message": "Hello, World!"
	}
}

```

1. Subscribe to a topic: Once a client sends a message with a specific topic (e.g., "example-topic"), it subscribes to that topic.
2. Broadcasting messages: The server will broadcast messages to all clients subscribed to the same topic.

### Example Client (JavaScript):
```js

const socket = new WebSocket("ws://localhost:8081/ws");

socket.onopen = function() {
    console.log("Connected to WebSocket");

    // Send a message to the server
    socket.send(JSON.stringify({
        topic: "chat",
        data: { user: "Alice", message: "Hello, everyone!" }
    }));
};

socket.onmessage = function(event) {
    const message = JSON.parse(event.data);
    console.log("Received message:", message);
};

socket.onclose = function() {
    console.log("Disconnected from WebSocket");
};

```

## License
This project is licensed under the MIT License - see the [LICENSE](LICENCE) file for details.




