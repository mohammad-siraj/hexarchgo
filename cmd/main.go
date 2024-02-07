package main

import (
	"context"
	"fmt"
	"log"
	"net"

	httpServer "github.com/mohammad-siraj/hexarchgo/internal/libs/http"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
	ports "github.com/mohammad-siraj/hexarchgo/internal/ports"
)

func main() {
	HTTPServerPort := ":8090"
	GRPCServerPort := ":8081"
	isGrpcEnabled := true
	ctx := context.Background()
	loggerConfig := logger.NewlogConfigOptions(false)
	loggerInstance := logger.NewLogger(loggerConfig)
	loggerInstance.Info(ctx, "Starting server ...\n")
	serverHttp, err := httpServer.NewHttpServer(isGrpcEnabled, loggerInstance)
	if err != nil {
		log.Fatal(err)
	}
	if isGrpcEnabled {
		ports.RegisterServicesToGrpcServer(serverHttp)
		grpcServer := serverHttp.GetGrpcServerInstanceForRegister()
		lis, err := net.Listen("tcp", GRPCServerPort)
		if err != nil {
			log.Fatal(err)
		}
		GRPCServerPortIp := fmt.Sprintf("0.0.0.0%s", GRPCServerPort)
		go grpcServer.ServerPort(lis)
		grpcServer.ConnectClient(GRPCServerPortIp)
		ports.RegisterServiceHandlersToGrpcServer(ctx, grpcServer)
		grpcServer.StartGrpcServer(ctx, serverHttp)
	}

	err = serverHttp.Run(HTTPServerPort)
	if err != nil {
		log.Fatal(err)
	}

}
