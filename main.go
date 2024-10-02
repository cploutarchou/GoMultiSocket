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

func serveWs(hub *hub.Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	hub.HandleClientConnection(conn)
}

func main() {
	socketHub := hub.NewHub()
	go socketHub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(socketHub, w, r)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}
