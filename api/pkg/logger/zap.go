package logger

import (
	"chat-app/pkg/utils"
	"fmt"
	"log"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type level int

const (
	_ level = iota
	DEBUG
	INFO
	WARN
	ERR
	FATAL
)

type Logger struct {
	*zap.Logger
}

func New() *Logger {
	var opts []zap.Option

	core := zapcore.NewTee(consoleCore(), fileCore())
	opts = append(opts, zap.AddCaller(), zap.AddCallerSkip(1))

	logger := zap.New(core, opts...).With(zap.String("service", "chat-app"))

	return &Logger{Logger: logger}
}

func (l *Logger) Stop() {
	_ = l.Logger.Sync() // ignore Sync error: because the stdout isn't flushable
	log.Printf("[logger] stopped")
}

func (l *Logger) C() *zap.Logger                          { return l.Logger }
func (l *Logger) Info(scope string, fields ...zap.Field)  { l.Logger.Info(scope, fields...) }
func (l *Logger) Error(scope string, fields ...zap.Field) { l.Logger.Error(scope, fields...) }
func (l *Logger) Warn(scope string, fields ...zap.Field)  { l.Logger.Warn(scope, fields...) }

// HELPERS

func encoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:      "ts",
		LevelKey:     "level",
		MessageKey:   "msg",
		CallerKey:    "caller",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeLevel:  zapcore.LowercaseLevelEncoder,
	}
}

func levelEnabler(lvl level) zap.LevelEnablerFunc {
	levels := map[level]zap.LevelEnablerFunc{
		DEBUG: zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return lvl >= zapcore.DebugLevel }),
		INFO:  zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return lvl == zapcore.InfoLevel }),
		WARN:  zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return lvl == zapcore.WarnLevel }),
		ERR:   zap.LevelEnablerFunc(func(lvl zapcore.Level) bool { return lvl == zapcore.ErrorLevel }),
	}

	return levels[lvl]
}

func consoleCore() zapcore.Core {
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig()),
		zapcore.AddSync(os.Stdout),
		levelEnabler(DEBUG),
	)
}

func fileCore() zapcore.Core {
	fileName := fmt.Sprintf("%s/logs/%s.log", utils.Root(), time.Now().Format("2006-01-02-15-01"))
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}

	// NOTE: there could be multiple Cores per level

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig()),
		zapcore.AddSync(file),
		levelEnabler(DEBUG),
	)
}
