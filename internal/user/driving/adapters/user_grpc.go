package adapters

import (
	"context"

	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/database"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/http"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/middleware"
	user "github.com/mohammad-siraj/hexarchgo/internal/user/driving/adapters/proto/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type UserHandler struct {
	server http.IHttpClient
	log    logger.ILogger
	cache  database.IDatabase
	user.UnimplementedUserServer
}

func NewUserHandler(h http.IHttpClient, l logger.ILogger, cacheClient database.IDatabase) *UserHandler {
	return &UserHandler{
		server: h,
		log:    l,
		cache:  cacheClient,
	}
}

// NOTE: make sure the service name is same as proto definition
func (u *UserHandler) RegisterUser(ctx context.Context, in *user.UserRegisterRequest) (*user.UserRegisterReponse, error) {
	u.log.Info(ctx, "Request received", logger.NewLogFieldInput("methodName", "GetUser"))

	token, err := middleware.CreateToken("user")
	if err != nil {
		return &user.UserRegisterReponse{
			Status: "FAILED",
			UserId: "",
		}, nil
	}
	if _, err := u.cache.ExecWithContext(ctx, "SET "+in.Email+" "+token+"-"+in.Password); err != nil {
		u.log.Error(ctx, err.Error())
		return &user.UserRegisterReponse{
			Status: "FAILED",
			UserId: "",
		}, nil
	}
	header := metadata.New(map[string]string{"authentication": token})
	if err := grpc.SendHeader(ctx, header); err != nil {
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

	token, err := middleware.CreateToken("user")
	if err != nil {
		return &user.UserLoginResponse{
			Status: "FAILED",
			UserId: "",
		}, nil
	}
	header := metadata.New(map[string]string{"authentication": token})
	if err := grpc.SendHeader(ctx, header); err != nil {
		return nil, err
	}

	tokenData, err := u.cache.ExecWithContext(ctx, "GET "+in.Email)
	if err != nil {
		u.log.Error(ctx, err.Error())
		return &user.UserLoginResponse{
			Status: "FAILED",
			UserId: "",
		}, nil
	}
	return &user.UserLoginResponse{
		Status: "OK",
		UserId: tokenData,
	}, nil
}
