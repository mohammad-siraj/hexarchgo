package ports

import (
	"context"

	userHandler "github.com/mohammad-siraj/hexarchgo/internal/domain/user/driving/adapters/proto"
	user "github.com/mohammad-siraj/hexarchgo/internal/domain/user/driving/adapters/proto/service"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/http"
)

func RegisterServicesToGrpcServer(httpServer http.IHttpClient) {
	server := httpServer.GetGrpcServerInstanceForRegister().GetServerInstanceForRegister()
	handlerInstance := userHandler.NewUserHandler()
	user.RegisterUserServer(server, handlerInstance)
}

func RegisterServiceHandlersToGrpcServer(ctx context.Context, grpcGw http.IgrpcGw) error {
	err := user.RegisterUserHandler(ctx, grpcGw.GetMuxInstanceForRegister(), grpcGw.GetClientInstanceForRegister())
	return err
}
