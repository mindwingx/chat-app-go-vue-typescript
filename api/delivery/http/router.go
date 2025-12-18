package http

import (
	"chat-app/delivery/http/handler"
	"chat-app/domain/port"
	"chat-app/domain/service"
	"net/http"
)

type Router struct {
	mux       *http.ServeMux
	wsHandler *handler.WebSocketHandler
}

func NewRouter(logger port.ILogger, websocketService *service.WebSocketService) *Router {
	return &Router{
		mux:       http.NewServeMux(),
		wsHandler: handler.NewWebSocketHandler(logger, websocketService),
	}
}

func (r *Router) SetRoutes() {
	r.mux.HandleFunc("/handshake", handler.Handshake)
	r.mux.HandleFunc("/ws", r.wsHandler.HandleWebSocket)
}

func (r *Router) Mux() *http.ServeMux {
	return r.mux
}
