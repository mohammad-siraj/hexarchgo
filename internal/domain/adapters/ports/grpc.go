package ports

import (
	"context"

	handler "github.com/mohammad-siraj/hexarchgo/internal/domain/adapters/handlers"
	"github.com/mohammad-siraj/hexarchgo/internal/domain/adapters/proto/apis/user"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/http"
)

func RegisterServicesToGrpcServer(grpcGw http.IgrpcGw) {
	server := grpcGw.GetServerInstanceForRegister()
	handlerInstance := handler.NewUserHandler()
	user.RegisterUserServer(server, handlerInstance)
}

func RegisterServiceHandlersToGrpcServer(ctx context.Context, grpcGw http.IgrpcGw) error {
	err := user.RegisterUserHandler(ctx, grpcGw.GetMuxInstanceForRegister(), grpcGw.GetClientInstanceForRegister())
	return err
}
