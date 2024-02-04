package web

import (
	"fmt"
	"net/http"
	"os"
	"time"

	ws "github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 3 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = ws.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type WsClient struct {
	hub  *WsHub
	conn *ws.Conn
	send chan []byte
}

func NewWsClient(hub *WsHub, w http.ResponseWriter, r *http.Request) (*WsClient, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	return &WsClient{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}, nil
}

func (c *WsClient) HandleReads() {
	defer func() {
		c.hub.Unregister <- c
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return err
	})
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		// we don't process any client messages except for pongs, the loop simply blocks until the connection is closed
	}
}

// HandleWrites writes messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *WsClient) HandleWrites() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(ws.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(ws.TextMessage)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot open new writer for websocket connection: %v", err)
				return
			}
			_, err = w.Write(message)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot open new writer for websocket connection: %v", err)
				return
			}
			if err := w.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to close writer for websocket connection: %v", err)
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(ws.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
