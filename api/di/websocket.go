package di

import (
	"chat-app/domain/port"
	"chat-app/domain/service"
)

func ProvideWebSocketService(logger port.ILogger) *service.WebSocketService {
	return service.NewWebSocketService(logger)
}
