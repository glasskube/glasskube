package sse

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

// sseHub maintains the set of active clients and broadcasts messages to the clients.
type sseHub struct {
	// Inbound messages from the clients.
	broadcast chan *sse

	// register requests from the clients.
	register chan *sseClient

	// unregister requests from clients.
	unregister chan *sseClient

	// Registered clients.
	clients sync.Map // map[*sseClient]struct{}
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

// newHub creates new sseHub
func newHub() *sseHub {
	return &sseHub{
		broadcast:  make(chan *sse),
		register:   make(chan *sseClient),
		unregister: make(chan *sseClient),
		clients:    sync.Map{},
	}
}

// Run handles communication operations with sseHub
func (h *sseHub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients.Store(client, struct{}{})
		case client := <-h.unregister:
			close(client.send)
			h.clients.Delete(client)
		case message := <-h.broadcast:
			h.clients.Range(func(key, value any) bool {
				if client, ok := key.(*sseClient); ok {
					client.send <- message
				}
				return true
			})
		}
	}
}

func (h *sseHub) handler(w http.ResponseWriter) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		fmt.Fprintf(os.Stderr, "server sent events not supported\n")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	client := &sseClient{
		send: make(chan *sse),
	}
	h.register <- client

	// for some reason we need to send some initial data – otherwise following updates are not acknowledged by the browser
	_, _ = w.Write((&sse{}).ClientBytes())
	flusher.Flush()

	for evt := range client.send {
		if _, err := w.Write(evt.ClientBytes()); err != nil {
			break
		}
		flusher.Flush()
	}
	h.unregister <- client
}
