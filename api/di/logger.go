package di

import (
	"chat-app/domain/port"
	"chat-app/pkg/logger"
)

func ProvideLogger() port.ILogger {
	return logger.New()
}
