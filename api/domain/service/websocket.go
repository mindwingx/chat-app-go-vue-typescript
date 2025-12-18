package service

import (
	"chat-app/domain/entity"
	"chat-app/domain/port"
	"chat-app/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type EventType string

const (
	NotificationEvent EventType = "notification"
	OnlineUsersEvent  EventType = "online-users"
	TypingEvent       EventType = "typing"
	MessageEvent      EventType = "message"
)

const (
	JoinAction  string = "joined"
	LeaveAction string = "left"
)

type WebSocketService struct {
	logger          port.ILogger
	mutex           sync.Mutex
	clients         sync.Map // used for concurrent access
	typingUsers     map[string]struct{}
	lastTenMessages []Response
	broadcast       chan []byte
	retry           map[string][]byte
}

type (
	Response struct {
		Username string  `json:"username"`
		Content  Content `json:"content"`
		Time     string  `json:"time"`
	}

	Content struct {
		Type  EventType `json:"type"`
		Id    string    `json:"id,omitempty"`
		Value string    `json:"value,omitempty"`
		Extra any       `json:"extra,omitempty"`
	}
)

func NewWebSocketService(logger port.ILogger) *WebSocketService {
	return &WebSocketService{
		logger:          logger,
		clients:         sync.Map{},
		typingUsers:     make(map[string]struct{}),
		lastTenMessages: make([]Response, 0),
		broadcast:       make(chan []byte, 100),
		retry:           map[string][]byte{},
	}
}

//

func (ws *WebSocketService) GetOnlineUsers() []string {
	users := make([]string, 0)

	ws.clients.Range(func(c, username interface{}) (res bool) {
		res = true
		users = append(users, username.(string))
		return
	})

	return users
}

func (ws *WebSocketService) GetLastTenMessages() []Response {
	return ws.lastTenMessages
}

func (ws *WebSocketService) SetTypingUser(username string) {
	ws.mutex.Lock()
	ws.typingUsers[username] = struct{}{}
	ws.mutex.Unlock()
}

//

func (ws *WebSocketService) RegisterClient(client *entity.WebSocketClient) {
	ws.clients.Store(client, client.Username())
	ws.logger.Info("service.ws.client.registered", zap.Any("username", client.Username()))

	notifValue := fmt.Sprintf("%s %s", client.Username(), JoinAction)
	ws.BroadcastEvent(NotificationEvent, client.Username(), notifValue, ws.GetOnlineUsers())
}

func (ws *WebSocketService) UnregisterClient(client *entity.WebSocketClient) {
	if _, loaded := ws.clients.LoadAndDelete(client); loaded {
		client.Close()
		ws.logger.Info("service.ws.client.unregistered", zap.Any("username", client.Username()))

		notifValue := fmt.Sprintf("%s %s", client.Username(), LeaveAction)
		ws.BroadcastEvent(NotificationEvent, client.Username(), notifValue, ws.GetOnlineUsers())
	}
}

//

func (ws *WebSocketService) BroadcastEvent(eventType EventType, username, value string, extra any) {
	defer func() {
		if err := recover(); err != nil {
			ws.logger.Error("service.ws.broadcast.eventType",
				zap.Any("eventType", eventType),
				zap.Any("value", value),
				zap.Any("extra", extra),
				zap.Any("error", err),
			)
		}
	}()

	messageBytes, _ := json.Marshal(Response{
		Username: username,
		Content: Content{
			Type:  eventType,
			Value: value,
			Extra: extra,
		},
		Time: time.Now().Format("15:04:05"),
	})

	ws.broadcast <- messageBytes
}

func (ws *WebSocketService) BroadcastMessage(client *entity.WebSocketClient, content string) {
	username, exists := ws.clients.Load(client)
	if !exists {
		return
	}

	messageBytes, _ := json.Marshal(Response{
		Username: username.(string),
		Content: Content{
			Type:  MessageEvent,
			Id:    utils.RandomStr(10),
			Value: content,
		},
		Time: time.Now().Format("15:04:05"),
	})

	ws.broadcast <- messageBytes
}

//

func (ws *WebSocketService) Broadcast(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			ws.logger.Warn("service.ws.broadcast.channel.closed")
			return
		case broadcast := <-ws.broadcast:
			var response Response
			_ = json.Unmarshal(broadcast, &response)

			if response.Content.Type == MessageEvent {
				ws.mutex.Lock()
				if len(ws.lastTenMessages) == 10 {
					ws.lastTenMessages = ws.lastTenMessages[1:]
				}
				ws.lastTenMessages = append(ws.lastTenMessages, response)
				ws.mutex.Unlock()
			}

			ws.clients.Range(func(client, username interface{}) (res bool) {
				res = true

				defer func() {
					if err := recover(); err != nil {
						ws.logger.Warn("service.broadcast.recovered", zap.Any("panic", err))
					}
				}()

				wsClient := client.(*entity.WebSocketClient)

				if response.Username == wsClient.Username() {
					return
				}

				select {
				case wsClient.SendChannel() <- broadcast:
					// Response sent successfully
				case <-time.After(1 * time.Second):
					if ws.retry[response.Content.Id] == nil {
						ws.logger.Warn("service.ws.broadcast.channel-buffer.full", zap.Any("username", username))
						ws.retry[response.Content.Id] = broadcast
					}
				}

				return
			})
		}
	}
}

func (ws *WebSocketService) RetryFailedItems(ctx context.Context) {
	var (
		cond     = sync.NewCond(&ws.mutex)
		ready    = true
		progress = false
	)

	for {
		select {
		case <-ctx.Done():
			ws.logger.Warn("service.ws.broadcast.retry.closed")
			return
		default:
			if len(ws.retry) > 0 && progress == false {
				ready = false
				progress = true

				go func() {
					for _, msg := range ws.retry {
						func([]byte) {
							defer func() {
								if err := recover(); err != nil {
									ws.logger.Error("service.ws.broadcast.retry.failed", zap.Any("error", err))
								}
							}()

							var resp Response
							_ = json.Unmarshal(msg, &resp)

							ws.broadcast <- msg
							delete(ws.retry, resp.Content.Id)
						}(msg)
					}

					time.Sleep(3 * time.Second)

					ws.mutex.Lock()
					ready = true
					cond.Signal()
					ws.mutex.Unlock()
				}()
			}

			if !ready {
				ws.mutex.Lock()
				cond.Wait()
				progress = false
				ws.mutex.Unlock()
			}
		}
	}
}

func (ws *WebSocketService) ReleaseTypingUsers(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			ws.logger.Warn("service.ws.broadcast.typing-users.closed")
			return
		default:
			users := make([]string, 0, len(ws.typingUsers))

			if len(ws.typingUsers) > 0 {
				ws.mutex.Lock()

				for k := range ws.typingUsers {
					users = append(users, k)
				}
				ws.typingUsers = make(map[string]struct{})
				ws.mutex.Unlock()
			}

			ws.BroadcastEvent(TypingEvent, "", "", users)
			time.Sleep(time.Duration(500) * time.Millisecond)
		}
	}
}

//

func (ws *WebSocketService) Terminate() {
	ws.logger.Info("service.ws.clients.connections.closing")

	clientCount := 0
	ws.clients.Range(func(key, value interface{}) bool {
		clientCount++
		return true
	})

	if clientCount == 0 {
		ws.logger.Info("service.ws.clients.connections.close", zap.Int("active-connections", clientCount))
		return
	}

	closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "service stopped")

	ws.clients.Range(func(key, value interface{}) bool {
		client := key.(*entity.WebSocketClient)

		if client.Conn() != nil {
			err := client.Conn().WriteMessage(websocket.CloseMessage, closeMsg)
			if err != nil {
				ws.logger.Error("service.ws.connection.close.msg.failed", zap.Error(err))
			}

			err = client.Conn().Close()
			if err != nil {
				ws.logger.Error("service.ws.connection.close.failure", zap.Error(err))
			}
		}

		ws.clients.Delete(client)
		client.Close()

		return true
	})

	close(ws.broadcast)

	ws.logger.Info("service.ws.connections.terminated", zap.Any("total-clients", clientCount))
}
