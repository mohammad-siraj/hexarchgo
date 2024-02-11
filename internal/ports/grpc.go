package ports

import (
	"context"

	"github.com/mohammad-siraj/hexarchgo/internal/libs/http"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
	userHandler "github.com/mohammad-siraj/hexarchgo/internal/user/driving/adapters"
	user "github.com/mohammad-siraj/hexarchgo/internal/user/driving/adapters/proto/service"
)

func RegisterServicesToGrpcServer(httpServer http.IHttpClient, log logger.ILogger) {
	server := httpServer.GetGrpcServerInstanceForRegister().GetServerInstanceForRegister()
	handlerInstance := userHandler.NewUserHandler(httpServer, log)
	user.RegisterUserServer(server, handlerInstance)
}

func RegisterServiceHandlersToGrpcServer(ctx context.Context, grpcGw http.IgrpcGw) error {
	err := user.RegisterUserHandler(ctx, grpcGw.GetMuxInstanceForRegister(), grpcGw.GetClientInstanceForRegister())
	return err
}
