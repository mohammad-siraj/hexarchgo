package ports

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/mohammad-siraj/hexarchgo/internal/libs/database/cache"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/http"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
	user "github.com/mohammad-siraj/hexarchgo/internal/user/driving/adapters"
)

type IPorter interface {
	RegisterRequestHandlers()
	RegisterServicesToGrpcServer()
	RegisterServiceHandlersToGrpcServer(ctx context.Context) error
	SetUpGrpcGateWay(ctx context.Context, GRPCServerPort string)
}

type porter struct {
	server      http.IHttpClient
	log         logger.ILogger
	cacheClient cache.ICacheClient
}

func NewPorter(h http.IHttpClient, l logger.ILogger, cacheClient cache.ICacheClient) IPorter {
	return &porter{
		server:      h,
		log:         l,
		cacheClient: cacheClient,
	}
}

func (p *porter) RegisterRequestHandlers() {
	userHandler := user.NewUserHandler(p.server, p.log, p.cacheClient)
	userSubRoute := p.server.NewSubGroup("/auth")
	{
		userSubRoute.Post("/register", userHandler.RegisterUserHandler)
	}
}

func (p *porter) SetUpGrpcGateWay(ctx context.Context, GRPCServerPort string) {
	p.RegisterServicesToGrpcServer()
	grpcServer := p.server.GetGrpcServerInstanceForRegister()
	lis, err := net.Listen("tcp", GRPCServerPort)
	if err != nil {
		log.Fatal(err)
	}
	GRPCServerPortIp := fmt.Sprintf("0.0.0.0%s", GRPCServerPort)
	go grpcServer.ServerPort(lis)
	if err := grpcServer.ConnectClient(GRPCServerPortIp); err != nil {
		p.log.Error(context.Background(), "Failed to connect client with gRPC")
		return
	}
	if err := p.RegisterServiceHandlersToGrpcServer(ctx); err != nil {
		p.log.Error(context.Background(), "Failed to connect client with gRPC")
		return
	}
	if err := p.server.GetGrpcServerInstanceForRegister().StartGrpcServer(ctx, p.server); err != nil {
		p.log.Error(context.Background(), "Failed to connect client with gRPC")
		return
	}
}
