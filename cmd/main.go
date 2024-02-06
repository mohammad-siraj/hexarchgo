package main

import (
	"context"
	"fmt"
	"log"
	"net"

	ports "github.com/mohammad-siraj/hexarchgo/internal/domain/adapters/ports"
	httpServer "github.com/mohammad-siraj/hexarchgo/internal/libs/http"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
)

func main() {
	HTTPServerPort := ":8090"
	GRPCServerPort := ":8081"
	isGrpcEnabled := true
	ctx := context.Background()
	serverHttp, err := httpServer.NewHttpServer(isGrpcEnabled)
	if err != nil {
		log.Fatal(err)
	}
	loggerConfig := logger.NewlogConfigOptions(false)
	loggerInstance := logger.NewLogger(loggerConfig)
	loggerInstance.Info(ctx, "Starting server INSTANCE...\n")
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
