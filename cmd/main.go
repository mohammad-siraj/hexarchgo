package main

import (
	"context"
	"log"
	"net"

	ports "github.com/mohammad-siraj/hexarchgo/internal/domain/adapters/ports"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/http"
)

func main() {
	c, err := http.NewHttpClient()
	if err != nil {
		log.Fatal(err)
	}

	grpcServer, err := http.NewGrpcGw("0.0.0.0:8080", c)
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	//user.RegisterUserHandlerServer(context.Background(),)

	ports.RegisterServicesToGrpcServer(grpcServer)

	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	go grpcServer.ServerPort(lis)

	grpcServer.ConnectClient("0.0.0.0:8080")

	ports.RegisterServiceHandlersToGrpcServer(ctx, grpcServer)

	grpcServer.StartGrpcServer(ctx, ":8080", ":8090", "0.0.0.0:8080", c)

	// subGroup := c.NewSubGroup("/v1")
	// {
	// 	//subGroup.Any("*{grpc_gateway}")
	// }

}
