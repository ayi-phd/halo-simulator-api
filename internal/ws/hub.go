package ws

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"halo-simulator/internal/events"
)

var upgrader = websocket.Upgrader{
	// Allow all origins for development.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// message is the JSON envelope sent to dashboard clients.
type message struct {
	Type    events.EventType `json:"type"`
	Payload any              `json:"payload"`
}

// Hub owns the connected-client registry. Only its Run goroutine
// touches the clients map — no mutex needed.
type Hub struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	bus        *events.Bus
}

func NewHub(bus *events.Bus) *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		bus:        bus,
	}
}

// Run starts the hub's dispatch loop. Call in its own goroutine.
func (h *Hub) Run(ctx context.Context) {
	sub := h.bus.Subscribe()
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
			log.Printf("ws: client connected (%d total)", len(h.clients))

		case c := <-h.unregister:
			if _, ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
				log.Printf("ws: client disconnected (%d total)", len(h.clients))
			}

		case e := <-sub:
			data, err := json.Marshal(message{Type: e.Type, Payload: e.Payload})
			if err != nil {
				log.Printf("ws: marshal error: %v", err)
				continue
			}
			for c := range h.clients {
				select {
				case c.send <- data:
				default:
					// Client send buffer full — drop the slow client.
					close(c.send)
					delete(h.clients, c)
				}
			}

		case <-ctx.Done():
			return
		}
	}
}

// ServeWS upgrades an HTTP connection to WebSocket and registers the client
// with the hub. Meant to be used directly as an http.HandlerFunc.
func (h *Hub) ServeWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("ws: upgrade error: %v", err)
		return
	}
	c := &Client{hub: h, conn: conn, send: make(chan []byte, 64)}
	h.register <- c
	go c.writePump()
	go c.readPump()
}
