package ws

import (
	"log"

	"github.com/gorilla/websocket"
)

// Client represents a single connected WebSocket dashboard client.
type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

// readPump reads from the WebSocket solely to detect disconnects.
// All server → client traffic flows through writePump.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		if _, _, err := c.conn.ReadMessage(); err != nil {
			break
		}
	}
}

// writePump drains the send channel and writes each message to the WebSocket.
// It exits when the send channel is closed by the hub (client dropped or shutdown).
func (c *Client) writePump() {
	defer c.conn.Close()
	for msg := range c.send {
		if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			log.Printf("ws: write error: %v", err)
			return
		}
	}
}
