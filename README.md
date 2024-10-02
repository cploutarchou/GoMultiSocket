# GoMultiSocket

A SOLID Go WebSocket implementation with support for multiple topics.

## Features

- WebSocket connection management
- Support for multiple topics (channels)
- Broadcast messages to specific topics
- Lightweight and scalable

## Installation

```bash
go get github.com/cploutarchou/GoMultiSocket/hub
```

## Usage

```go
package main

import (
	"github.com/cploutarchou/GoMultiSocket/hub"
	"github.com/gorilla/websocket"

	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func serveWs(hub *socket.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	hub.HandleClientConnection(conn)
}

func main() {
	hub := socket.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}


```
## License
MIT [License](LICENCE) Copyright (c) 2024 Christos Ploutarchou 