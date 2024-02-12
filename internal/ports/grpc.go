package ports

import (
	"context"

	userHandler "github.com/mohammad-siraj/hexarchgo/internal/user/driving/adapters"
	user "github.com/mohammad-siraj/hexarchgo/internal/user/driving/adapters/proto/service"
)

func (p *porter) RegisterServicesToGrpcServer() {
	server := p.server.GetGrpcServerInstanceForRegister().GetServerInstanceForRegister()
	handlerInstance := userHandler.NewUserHandler(p.server, p.log, p.cacheClient)
	user.RegisterUserServer(server, handlerInstance)
}

func (p *porter) RegisterServiceHandlersToGrpcServer(ctx context.Context) error {
	grpcGw := p.server.GetGrpcServerInstanceForRegister()
	err := user.RegisterUserHandler(ctx, grpcGw.GetMuxInstanceForRegister(), grpcGw.GetClientInstanceForRegister())
	return err
}
