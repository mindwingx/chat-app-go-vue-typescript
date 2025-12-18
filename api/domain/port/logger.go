package port

import "go.uber.org/zap"

type ILogger interface {
	C() *zap.Logger
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Stop()
}
