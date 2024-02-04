package web

import (
	"fmt"
	"net/http"
	"os"
)

// WsHub maintains the set of active clients and broadcasts messages to the clients.
type WsHub struct {
	// Inbound messages from the clients.
	Broadcast chan []byte

	// Register requests from the clients.
	Register chan *WsClient

	// Unregister requests from clients.
	Unregister chan *WsClient

	// Registered clients.
	clients map[*WsClient]bool
}

// NewHub creates new WsHub
func NewHub() *WsHub {
	return &WsHub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *WsClient),
		Unregister: make(chan *WsClient),
		clients:    make(map[*WsClient]bool),
	}
}

// Run handles communication operations with WsHub
func (h *WsHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.Broadcast:
			for client := range h.clients {
				client.send <- message
			}
		}
	}
}

func (h *WsHub) handler(w http.ResponseWriter, r *http.Request) {
	wsClient, err := NewWsClient(h, w, r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create WebSocket client: %v", err)
		return
	}
	h.Register <- wsClient
	go wsClient.HandleReads()
	go wsClient.HandleWrites()
}
