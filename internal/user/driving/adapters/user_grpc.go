package adapters

import (
	"context"

	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/http"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
	user "github.com/mohammad-siraj/hexarchgo/internal/user/driving/adapters/proto/service"
)

type UserHandler struct {
	server http.IHttpClient
	log    logger.ILogger
	user.UnimplementedUserServer
}

// mustEmbedUnimplementedUserServer implements user.UserServer.
func (u *UserHandler) mustEmbedUnimplementedUserServer() {
	panic("unimplemented")
}

func NewUserHandler(h http.IHttpClient, l logger.ILogger) *UserHandler {
	return &UserHandler{
		server: h,
		log:    l,
	}
}

// NOTE: make sure the service name is same as proto definition
func (u *UserHandler) RegisterUser(ctx context.Context, in *user.UserRegisterRequest) (*user.UserRegisterReponse, error) {
	u.log.Info(ctx, "Request received", logger.NewLogFieldInput("methodName", "GetUser"))
	return &user.UserRegisterReponse{
		Status: "OK",
		UserId: "1234",
	}, nil
}
