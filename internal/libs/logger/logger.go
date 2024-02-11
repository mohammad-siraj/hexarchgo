package logger

import (
	"context"
	"io"
	"net/http"
	"os"
	"time"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	timeFormat = time.RFC3339
)

type ILogger interface {
	Sync()
	Debug(ctx context.Context, message string, fieldData ...ILogFieldInput)
	Info(ctx context.Context, message string, fieldData ...ILogFieldInput)
	Warn(ctx context.Context, message string, fieldData ...ILogFieldInput)
	Error(ctx context.Context, message string, fieldData ...ILogFieldInput)
	GetIoWriter() io.Writer
	RequestLog(request *http.Request)
	GetGrpcUnaryInterceptor() grpc.UnaryServerInterceptor
}

type Logger struct {
	log  *zap.Logger
	file zapcore.WriteSyncer
}

func NewLogger(config Iconfigs) ILogger {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   config.FileName(),
		MaxSize:    config.MaxSize(), // megabytes
		MaxBackups: config.MaxBackups(),
		MaxAge:     config.MaxAge(), // days
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
	log := zap.New(core)
	grpc_zap.ReplaceGrpcLogger(log)

	return &Logger{
		log:  log,
		file: file,
	}
}

func (l *Logger) GetIoWriter() io.Writer {
	return l.file
}

func (l *Logger) Sync() {
	if err := l.log.Sync(); err != nil {
		return
	}
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

// grpc logger
func (l *Logger) GetGrpcUnaryInterceptor() grpc.UnaryServerInterceptor {
	return grpc_zap.UnaryServerInterceptor(l.log)
}

func (l *Logger) RequestLog(request *http.Request) {
	start := time.Now()
	query := request.URL.RawQuery
	end := time.Now()
	latency := end.Sub(start)
	end = end.UTC()
	path := request.URL.Path

	fields := []zapcore.Field{
		zap.Int("status", request.Response.StatusCode),
		zap.String("method", request.Method),
		zap.String("path", path),
		zap.String("query", query),
		zap.String("ip", request.RemoteAddr),
		zap.String("user-agent", request.UserAgent()),
		zap.Duration("latency", latency),
	}
	fields = append(fields, zap.String("time", end.Format(timeFormat)))
	l.Info(context.Background(), request.Method+" "+path, NewLogFieldInput("headers", fields))
}

type ILogFieldInput interface {
	GetFieldName() string
	GetField() interface{}
}

type LogFieldInput struct {
	fieldName string
	fields    interface{}
}

func NewLogFieldInput(fieldName string, fields interface{}) ILogFieldInput {
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
