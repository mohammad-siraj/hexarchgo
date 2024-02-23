package adapters

import (
	"context"

	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
	userDomain "github.com/mohammad-siraj/hexarchgo/internal/user/domain"
	user "github.com/mohammad-siraj/hexarchgo/internal/user/driving/adapters/proto/service"
)

type UserHandler struct {
	log        logger.ILogger
	userDomain userDomain.IUserDomain
	user.UnimplementedUserServer
}

func NewUserHandler(ctx context.Context, l logger.ILogger, userDomain userDomain.IUserDomain) *UserHandler {
	return &UserHandler{
		log:        l,
		userDomain: userDomain,
	}
}

// NOTE: make sure the service name is same as proto definition
func (u *UserHandler) RegisterUser(ctx context.Context, in *user.UserRegisterRequest) (*user.UserRegisterReponse, error) {
	u.log.Info(ctx, "Request received", logger.NewLogFieldInput("methodName", "GetUser"))

	err := u.userDomain.SendBackToken(ctx, in.Email)
	if err != nil {
		u.log.Error(ctx, "error in creating authentication token "+err.Error())
		return nil, err
	}
	return &user.UserRegisterReponse{
		Status: "OK",
		UserId: "1234",
	}, nil
}

// NOTE: make sure the service name is same as proto definition
func (u *UserHandler) LoginUser(ctx context.Context, in *user.UserLoginRequest) (*user.UserLoginResponse, error) {
	u.log.Info(ctx, "Request received", logger.NewLogFieldInput("methodName", "GetUser"))

	err := u.userDomain.SendBackToken(ctx, in.Email)
	if err != nil {
		u.log.Error(ctx, "error in creating authentication token "+err.Error())
		return nil, err
	}
	return &user.UserLoginResponse{
		Status: "OK",
		UserId: "1234567890",
	}, nil
}
