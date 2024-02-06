package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ILogger interface {
	Sync()
	Debug(ctx context.Context, message string, fieldData ...ILogFieldInput)
	Info(ctx context.Context, message string, fieldData ...ILogFieldInput)
	Warn(ctx context.Context, message string, fieldData ...ILogFieldInput)
	Error(ctx context.Context, message string, fieldData ...ILogFieldInput)
}

type Logger struct {
	log *zap.Logger
}

func NewLogger(config Iconfigs) ILogger {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})

	level := zap.NewAtomicLevelAt(zap.InfoLevel)

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)

	return &Logger{
		log: zap.New(core),
	}
}

func genarateLogRotater(c Iconfigs) zapcore.WriteSyncer {
	ll := &lumberjack.Logger{
		Filename:   c.FileName(),
		MaxSize:    c.MaxSize(), //MB
		MaxBackups: c.MaxBackups(),
		MaxAge:     c.MaxAge(), //days
		Compress:   c.IsCompressed(),
	}
	return zapcore.AddSync(ll)
}

func (l *Logger) Sync() {
	l.log.Sync()
}

func (l *Logger) Debug(ctx context.Context, message string, fieldData ...ILogFieldInput) {
	if len(fieldData) == 1 {
		l.log.Debug(message, zap.Any(fieldData[0].GetFieldName(), fieldData[0].GetField()))
		return
	}
	l.log.Debug(message)
}

func (l *Logger) Info(ctx context.Context, message string, fieldData ...ILogFieldInput) {
	if len(fieldData) == 1 {
		l.log.Info(message, zap.Any(fieldData[0].GetFieldName(), fieldData[0].GetField()))
		return
	}
	l.log.Info(message)
}

func (l *Logger) Warn(ctx context.Context, message string, fieldData ...ILogFieldInput) {
	if len(fieldData) == 1 {
		l.log.Warn(message, zap.Any(fieldData[0].GetFieldName(), fieldData[0].GetField()))
		return
	}
	l.log.Warn(message)
}

func (l *Logger) Error(ctx context.Context, message string, fieldData ...ILogFieldInput) {
	if len(fieldData) == 1 {
		l.log.Error(message, zap.Any(fieldData[0].GetFieldName(), fieldData[0].GetField()))
		return
	}
	l.log.Error(message)
}

type ILogFieldInput interface {
	GetFieldName() string
	GetField() interface{}
}

type LogFieldInput struct {
	fieldName string
	fields    interface{}
}

func NewLogFieldInput(fieldName string, fields interface{}) LogFieldInput {
	return LogFieldInput{
		fieldName: fieldName,
		fields:    fields,
	}
}

func (i LogFieldInput) GetFieldName() string {
	return i.fieldName
}

func (i LogFieldInput) GetField() interface{} {
	return i.fields
}
