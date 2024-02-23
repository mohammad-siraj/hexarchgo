package ports

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/mohammad-siraj/hexarchgo/internal/libs/database/cache"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/database/sql"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/http"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
	userDomain "github.com/mohammad-siraj/hexarchgo/internal/user/domain"
	userDriven "github.com/mohammad-siraj/hexarchgo/internal/user/driven"
	userDriving "github.com/mohammad-siraj/hexarchgo/internal/user/driving/adapters"
	user "github.com/mohammad-siraj/hexarchgo/internal/user/driving/adapters/proto/service"
)

type IPorter interface {
	RegisterRequestHandlers(ctx context.Context, isGrpcEnabled bool)
}

type porter struct {
	server         http.IHttpClient
	log            logger.ILogger
	cacheClient    cache.ICacheClient
	sqlClient      sql.ISqlDatabase
	grpcServerPort string
}

func NewPorter(h http.IHttpClient, l logger.ILogger, cacheClient cache.ICacheClient, sqlClient sql.ISqlDatabase, GRPCServerPort string) IPorter {
	return &porter{
		server:         h,
		log:            l,
		cacheClient:    cacheClient,
		sqlClient:      sqlClient,
		grpcServerPort: GRPCServerPort,
	}
}

func (p *porter) RegisterRequestHandlers(ctx context.Context, isGrpcEnabled bool) {

	//handler

	userDrivenHandler := userDriven.NewUserDriven(p.log, p.cacheClient, p.sqlClient)
	userDomainHandler := userDomain.NewUserDomain(p.log, userDrivenHandler)
	userHandler := userDriving.NewUserHandler(ctx, p.log, userDomainHandler)

	if isGrpcEnabled {
		//register service to grpc server
		user.RegisterUserServer(p.server.GetGrpcServerInstanceForRegister().GetServerInstanceForRegister(), userHandler)
		grpcServer := p.server.GetGrpcServerInstanceForRegister()
		lis, err := net.Listen("tcp", p.grpcServerPort)
		if err != nil {
			log.Fatal(err)
		}
		GRPCServerPortIp := fmt.Sprintf("0.0.0.0%s", p.grpcServerPort)
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

	userSubRoute := p.server.NewSubGroup("/auth")
	{
		userSubRoute.Post("/register", userHandler.RegisterUserHandler)
	}
}

func (p *porter) RegisterServiceHandlersToGrpcServer(ctx context.Context) error {
	grpcGw := p.server.GetGrpcServerInstanceForRegister()
	err := user.RegisterUserHandler(ctx, grpcGw.GetMuxInstanceForRegister(), grpcGw.GetClientInstanceForRegister())
	return err
}
