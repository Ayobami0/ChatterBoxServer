package service

import (
	"log"
	"sync"

	ws "github.com/gorilla/websocket"
)

type WebsocketService struct {
	hub map[string]*ws.Conn
	mu  sync.RWMutex
}

func NewWebsocketService() *WebsocketService {
	return &WebsocketService{
		hub: map[string]*ws.Conn{},
		mu:  sync.RWMutex{},
	}
}

func (w *WebsocketService) AddConnection(id string, connection *ws.Conn) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.hub[id] = connection
}

func (w *WebsocketService) RemoveConnection(id string) {
	w.mu.Lock()
	defer w.mu.Unlock()
	conn, ok := w.hub[id]

	if ok {
		conn.WriteMessage(ws.CloseMessage, ws.FormatCloseMessage(ws.CloseNormalClosure, ""))
		conn.Close()
		delete(w.hub, id)
	}
}

func (w *WebsocketService) RemoveConnectionWithError(id, error string, closeCode int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	conn, ok := w.hub[id]

	if ok {
		conn.WriteMessage(ws.CloseMessage, ws.FormatCloseMessage(closeCode, error))
		conn.Close()
		delete(w.hub, id)
	}
}

func (w *WebsocketService) BroadcastMessage(message interface{}, ids ...string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, id := range ids {
		conn, ok := w.hub[id]

    if !ok {
      continue
    }

		if err := conn.WriteJSON(message); err != nil {
      log.Println(err, id)
      w.RemoveConnection(id)
			continue
		}
	}
}
