package logger

import (
	"context"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type ILogger interface {
}

type Logger struct {
	log *zap.Logger
}

func NewLogger(config Iconfigs) ILogger {

	encoderCfg := zap.NewDevelopmentEncoderConfig()
	atomicLevel := zap.NewAtomicLevelAt(zap.DebugLevel)
	if config.IsProduction() {
		encoderCfg = zap.NewProductionEncoderConfig()
		atomicLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.RFC3339TimeEncoder
	writeSyncer := genarateLogRotater(config)

	logCore := zapcore.NewCore(zapcore.NewJSONEncoder(encoderCfg), writeSyncer, atomicLevel)

	// loggerConfig := zap.Config{
	// 	Level:             atomicLevel,
	// 	Development:       !(config.IsProduction()),
	// 	DisableCaller:     false,
	// 	DisableStacktrace: false,
	// 	Sampling:          nil,
	// 	Encoding:          "json",
	// 	EncoderConfig:     encoderCfg,
	// 	OutputPaths: []string{
	// 		"stderr",
	// 	},
	// 	ErrorOutputPaths: []string{
	// 		"stderr",
	// 	},
	// 	InitialFields: map[string]interface{}{
	// 		"pid": os.Getpid(),
	// 	},
	// }

	// loggerConfig.WriteSyncer = genarateLogRotater(config)
	return &Logger{
		log: zap.New(logCore),
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

func (l *Logger) Debug(ctx context.Context, message string, fields map[string]interface{}) {

}
