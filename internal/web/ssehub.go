package web

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

// SSEHub maintains the set of active clients and broadcasts messages to the clients.
type SSEHub struct {
	// Inbound messages from the clients.
	Broadcast chan *sse

	// Register requests from the clients.
	Register chan *sseClient

	// Unregister requests from clients.
	Unregister chan *sseClient

	// Registered clients.
	clients map[*sseClient]struct{}
}

type sse struct {
	event string
	data  string
}

func (evt *sse) ClientBytes() []byte {
	return []byte(fmt.Sprintf("event: %s\ndata: %s\n\n", evt.event, strings.ReplaceAll(evt.data, "\n", "")))
}

type sseClient struct {
	send chan *sse
}

// NewHub creates new SSEHub
func NewHub() *SSEHub {
	return &SSEHub{
		Broadcast:  make(chan *sse),
		Register:   make(chan *sseClient),
		Unregister: make(chan *sseClient),
		clients:    make(map[*sseClient]struct{}),
	}
}

// Run handles communication operations with SSEHub
func (h *SSEHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = struct{}{}
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

func (h *SSEHub) handler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		fmt.Fprintf(os.Stderr, "server sent events not supported\n")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	client := &sseClient{
		send: make(chan *sse),
	}
	h.Register <- client

	// for some reason we need to send some initial data – otherwise following updates are not acknowledged by the browser
	_, _ = w.Write((&sse{}).ClientBytes())
	flusher.Flush()

	for evt := range client.send {
		if _, err := w.Write(evt.ClientBytes()); err != nil {
			break
		}
		flusher.Flush()
	}
	h.Unregister <- client
}
