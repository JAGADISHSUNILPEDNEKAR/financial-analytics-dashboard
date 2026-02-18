package handlers

import (
	"go.uber.org/zap"
)

type WebSocketHub struct {
	clients    map[*Client]bool
	Broadcast  chan []byte
	Register   chan *Client
	Unregister chan *Client
	logger     *zap.Logger
}

func NewWebSocketHub(logger *zap.Logger) *WebSocketHub {
	return &WebSocketHub{
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		logger:     logger,
	}
}

func (h *WebSocketHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.clients[client] = true
			h.logger.Info("Client registered", zap.String("user_id", client.userID))
		case client := <-h.Unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				h.logger.Info("Client unregistered", zap.String("user_id", client.userID))
			}
		case message := <-h.Broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
