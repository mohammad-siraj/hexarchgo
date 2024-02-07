package handlers

import (
	"context"
	"fmt"

	user "github.com/mohammad-siraj/hexarchgo/internal/domain/user/driving/adapters/proto/service"
)

type UserHandler struct {
	user.UnimplementedUserServer
}

// mustEmbedUnimplementedUserServer implements user.UserServer.
func (u *UserHandler) mustEmbedUnimplementedUserServer() {
	panic("unimplemented")
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// NOTE: make sure the service name is same as proto definition
func (u *UserHandler) RegisterUser(ctx context.Context, in *user.UserRegisterRequest) (*user.UserRegisterReponse, error) {
	fmt.Println("its here")
	return &user.UserRegisterReponse{
		Status: "OK",
		UserId: "1234",
	}, nil
}
