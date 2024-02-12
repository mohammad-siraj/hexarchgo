package main

import (
	"context"
	"log"

	httpServer "github.com/mohammad-siraj/hexarchgo/internal/libs/http"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/middleware"
	"github.com/mohammad-siraj/hexarchgo/internal/ports"
)

func main() {
	HTTPServerPort := ":8090"
	GRPCServerPort := ":8081"
	isGrpcEnabled := true
	ctx := context.Background()
	StartServer(ctx, isGrpcEnabled, GRPCServerPort, HTTPServerPort)

}

func StartServer(ctx context.Context, isGrpcEnabled bool, GRPCServerPort string, HTTPServerPort string) {
	loggerConfig := logger.NewlogConfigOptions(false)

	//logger configs
	loggerConfig.WithFilename("log/app.log")
	loggerConfig.WithIsCompressed(true)
	loggerConfig.WithIsLocalTime(true)
	loggerConfig.WithMaxAge(1)
	loggerConfig.WithMaxBackUp(1)

	//logger
	loggerInstance := logger.NewLogger(loggerConfig)
	loggerInstance.Info(ctx, "Starting server ...\n")

	serverHttp, err := httpServer.NewHttpServer(loggerInstance.GetIoWriter(), httpServer.NewGrpcOptions(loggerInstance, middleware.GrpcAuthMiddleware(ctx)))
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	porter := ports.NewPorter(serverHttp, loggerInstance)
	if isGrpcEnabled {
		porter.SetUpGrpcGateWay(ctx, GRPCServerPort)
	}

	porter.RegisterRequestHandlers()

	err = serverHttp.Run(HTTPServerPort)
	if err != nil {
		log.Fatal(err)
	}
}
