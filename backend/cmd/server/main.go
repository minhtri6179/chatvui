package main

import (
	"log"
	"net/http"

	"backend/internal/websocket"
)

func main() {
	wsManager := websocket.NewManager()
	go wsManager.Run()

	http.HandleFunc("/ws", wsManager.HandleWebSocket)

	log.Println("Server starting on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
