package libs

import (
	"context"

	"go.uber.org/zap"
)

type ILogger interface {
}

type Logger struct {
	log *zap.Logger
}

func NewLogger() ILogger {
	logConfig := zap.Must(zap.NewProduction())
	return &Logger{
		log: logConfig,
	}
}

func (l *Logger) Debug(ctx context.Context, message string, fields map[string]interface{}) {
}
