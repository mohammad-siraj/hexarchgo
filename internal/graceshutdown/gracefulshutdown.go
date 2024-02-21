package graceshutdown

import (
	"context"
	"os"
	"os/signal"

	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
)

type IGracefulShutdown interface {
	InterruptShutdown(ctx context.Context, log logger.ILogger, opts ...interface {
		Exit(ctx context.Context, errChannel chan error)
	})
	ErrorHandler(log logger.ILogger)
}

type gracefulshutdown struct {
	errChannel      chan error
	interuptChannel chan os.Signal
	log             logger.ILogger
}

func NewGracefulShutDownHandler(errChannel chan error, loggerInstance logger.ILogger) IGracefulShutdown {
	interuptChannel := make(chan os.Signal, 1)
	signal.Notify(interuptChannel, os.Interrupt)
	return &gracefulshutdown{
		errChannel:      errChannel,
		interuptChannel: interuptChannel,
		log:             loggerInstance,
	}
}

func (g *gracefulshutdown) InterruptShutdown(ctx context.Context, log logger.ILogger, opts ...interface {
	Exit(ctx context.Context, errChannel chan error)
}) {
	<-g.interuptChannel
	g.log.Info(ctx, "Application received interrupt signal")
	for _, resource := range opts {
		log.Info(ctx, "Resource exit is happening")
		resource.Exit(ctx, g.errChannel)
	}
	os.Exit(0)
}

func (g *gracefulshutdown) ErrorHandler(log logger.ILogger) {
	select {
	case value := <-g.errChannel:
		log.Error(context.Background(), value.Error())
	}
	os.Exit(0)
}
