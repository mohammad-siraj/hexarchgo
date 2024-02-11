package main

import (
	"context"
	"fmt"
	"log"
	"net"

	httpServer "github.com/mohammad-siraj/hexarchgo/internal/libs/http"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
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

	serverHttp, err := httpServer.NewHttpServer(isGrpcEnabled, loggerInstance, loggerInstance.GetIoWriter())
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	if isGrpcEnabled {
		SetUpGrpcGateWay(ctx, serverHttp, loggerInstance, GRPCServerPort)
	}

	ports.RegisterRequestHandlers(serverHttp, loggerInstance)

	err = serverHttp.Run(HTTPServerPort)
	if err != nil {
		log.Fatal(err)
	}
}

func SetUpGrpcGateWay(ctx context.Context, serverHttp httpServer.IHttpClient, loggerInstance logger.ILogger, GRPCServerPort string) {
	ports.RegisterServicesToGrpcServer(serverHttp, loggerInstance)
	grpcServer := serverHttp.GetGrpcServerInstanceForRegister()
	lis, err := net.Listen("tcp", GRPCServerPort)
	if err != nil {
		log.Fatal(err)
	}
	GRPCServerPortIp := fmt.Sprintf("0.0.0.0%s", GRPCServerPort)
	go grpcServer.ServerPort(lis)
	if err := grpcServer.ConnectClient(GRPCServerPortIp); err != nil {
		loggerInstance.Error(context.Background(), "Failed to connect client with gRPC")
		return
	}
	if err := ports.RegisterServiceHandlersToGrpcServer(ctx, grpcServer); err != nil {
		loggerInstance.Error(context.Background(), "Failed to connect client with gRPC")
		return
	}
	if err := grpcServer.StartGrpcServer(ctx, serverHttp); err != nil {
		loggerInstance.Error(context.Background(), "Failed to connect client with gRPC")
		return
	}
}
