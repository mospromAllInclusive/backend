package web_sockets

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Message struct {
	Topic       string      `json:"topic"`
	EventAction string      `json:"eventAction"`
	Payload     interface{} `json:"payload"`
	EventTime   time.Time   `json:"eventTime"`
}

type client struct {
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	mu          sync.RWMutex
	subscribers map[string]map[*client]struct{}
	register    chan registration
	unregister  chan registration
	broadcast   chan Message
	ctx         context.Context
	cancel      context.CancelFunc
}

type registration struct {
	topic  string
	client *client
}

func NewHub(ctx context.Context) *Hub {
	_ctx, cancel := context.WithCancel(ctx)
	h := &Hub{
		subscribers: make(map[string]map[*client]struct{}),
		register:    make(chan registration),
		unregister:  make(chan registration),
		broadcast:   make(chan Message, 1024),
		ctx:         _ctx,
		cancel:      cancel,
	}
	go h.run()
	return h
}

func (h *Hub) run() {
	for {
		select {
		case <-h.ctx.Done():
			h.closeAll()
			return
		case r := <-h.register:
			h.mu.Lock()
			if _, ok := h.subscribers[r.topic]; !ok {
				h.subscribers[r.topic] = make(map[*client]struct{})
			}
			h.subscribers[r.topic][r.client] = struct{}{}
			h.mu.Unlock()
		case r := <-h.unregister:
			h.mu.Lock()
			if set, ok := h.subscribers[r.topic]; ok {
				if _, exists := set[r.client]; exists {
					delete(set, r.client)
					close(r.client.send)
					if len(set) == 0 {
						delete(h.subscribers, r.topic)
					}
				}
			}
			h.mu.Unlock()
		case msg := <-h.broadcast:
			b, err := json.Marshal(msg)
			if err != nil {
				log.Printf("ws marshal error: %v", err)
				continue
			}
			h.mu.RLock()
			set := h.subscribers[msg.Topic]
			for c := range set {
				select {
				case c.send <- b:
				default:
					go func(topic string, cl *client) { h.unregister <- registration{topic: topic, client: cl} }(msg.Topic, c)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) closeAll() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for topic, set := range h.subscribers {
		for c := range set {
			close(c.send)
			_ = c.conn.Close()
		}
		delete(h.subscribers, topic)
	}
}

func (h *Hub) Broadcast(topic, eventAction string, payload interface{}) {
	select {
	case h.broadcast <- Message{
		Topic:       topic,
		EventAction: eventAction,
		Payload:     payload,
		EventTime:   time.Now().UTC(),
	}:
	default:
		log.Printf("ws broadcast queue full (topic=%s, eventAction=%s)", topic, eventAction)
	}
}

func (h *Hub) BroadcastMany(topics []string, eventAction string, payload interface{}) {
	for _, topic := range topics {
		h.Broadcast(topic, eventAction, payload)
	}
}

func (h *Hub) Shutdown() { h.cancel() }
