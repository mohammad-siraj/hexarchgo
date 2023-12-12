package libs

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ILogger interface {
}

type ILogRotator interface {
}

type LogRotator struct {
	fileName   string
	maxSize    int
	maxBackups int
	localTime  bool
	compress   bool
}

type Logger struct {
	logRotator LogRotator
	log        *zap.Logger
}

func NewLogger(isProduction bool) ILogger {

	encoderCfg := zap.NewDevelopmentEncoderConfig()
	atomicLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	if isProduction {
		encoderCfg = zap.NewProductionEncoderConfig()
		atomicLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             atomicLevel,
		Development:       !(isProduction),
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	return &Logger{
		log: zap.Must(config.Build()),
	}
}

func NewLogRotator(fileName string) ILogRotator {
	return LogRotator{
		fileName: fileName,
	}
}

func (l *Logger) SetLogRotator(rotator ILogRotator) {
	l.log.WithOptions()
}

func (l *Logger) Debug(ctx context.Context, message string, fields map[string]interface{}) {

}
