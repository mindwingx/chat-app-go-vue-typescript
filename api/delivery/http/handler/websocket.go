package handler

import (
	"chat-app/domain/entity"
	"chat-app/domain/port"
	"chat-app/domain/service"
	"chat-app/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type WebSocketHandler struct {
	logger   port.ILogger
	service  *service.WebSocketService
	upgrader websocket.Upgrader
}

type Received struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func NewWebSocketHandler(logger port.ILogger, service *service.WebSocketService) *WebSocketHandler {
	return &WebSocketHandler{
		logger:  logger,
		service: service,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // all origin passed - dev purpose
			},
		},
	}
}

func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(context.Background())

	username := r.URL.Query().Get("username")
	if username == "" {
		username = fmt.Sprintf("anonymous-%s", utils.RandomStr(5))
	}

	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error("handler.ws.upgrade.failed", zap.Error(err))
		cancel()
		return
	}

	client := entity.NewWebSocketClient(conn, username)
	h.service.RegisterClient(client)

	go h.writePump(ctx, client)
	go h.readPump(cancel, client)

	h.initializeOnlineUsers(client)
	h.initializePreviousMessages(client)
}

//

func (h *WebSocketHandler) writePump(ctx context.Context, client *entity.WebSocketClient) {
	for {
		select {
		case <-ctx.Done():
			h.service.UnregisterClient(client)
			h.logger.Warn("handler.ws.write-pump.closed", zap.String("client", client.Username()))
			return
		case message := <-client.SendChannel():
			err := client.Conn().WriteMessage(websocket.TextMessage, message)
			if err != nil {
				h.logger.Error("handler.ws.write-pump.failed", zap.Error(err))
				break
			}
		}
	}
}

func (h *WebSocketHandler) readPump(cancel context.CancelFunc, client *entity.WebSocketClient) {
	for {
		select {
		default:
			_, message, err := client.Conn().ReadMessage()
			if err != nil {
				cancel()
				h.logger.Error("handler.ws.read-pump.ctx.cancel",
					zap.String("client", client.Username()), zap.Error(err))
				return
			}

			var clientMsg Received
			if err = json.Unmarshal(message, &clientMsg); err != nil {
				h.logger.Error("handler.ws.read-pump.unmarshal.failed", zap.Error(err))
				continue
			}

			switch clientMsg.Type {
			case "typing":
				h.service.SetTypingUser(client.Username())
			case "message":
				if clientMsg.Value != "" {
					h.service.BroadcastMessage(client, clientMsg.Value)
				}
			}
		}
	}
}

//

func (h *WebSocketHandler) initializeOnlineUsers(client *entity.WebSocketClient) {
	usersResp := service.Response{
		Content: service.Content{
			Type:  service.OnlineUsersEvent,
			Extra: h.service.GetOnlineUsers(),
		},
		Time: time.Now().Format("15:04:05"),
	}

	bytes, _ := json.Marshal(usersResp)

	err := client.Conn().WriteMessage(websocket.TextMessage, bytes)
	if err != nil {
		h.logger.Error("handler.ws.write-pump.send.users.failed", zap.Error(err))
	}
}

func (h *WebSocketHandler) initializePreviousMessages(client *entity.WebSocketClient) {
	if len(h.service.GetLastTenMessages()) > 0 {
		for _, msg := range h.service.GetLastTenMessages() {
			bytes, _ := json.Marshal(msg)

			if err := client.Conn().WriteMessage(websocket.TextMessage, bytes); err != nil {
				h.logger.Error("handler.ws.write-pump.send.last-messages.failed", zap.Error(err))
			}
		}
	}
}
