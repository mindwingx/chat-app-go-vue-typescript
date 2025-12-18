package http

import (
	"chat-app/domain/port"
	"net/http"

	"go.uber.org/zap"
)

func RecoveryMiddleware(logger port.ILogger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Warn("recovered", zap.Any("panic", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
