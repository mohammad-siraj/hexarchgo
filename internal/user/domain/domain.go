package domain

import (
	"context"

	"time"

	"github.com/mohammad-siraj/hexarchgo/internal/common"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/middleware"
	userDriven "github.com/mohammad-siraj/hexarchgo/internal/user/driven"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type IUserDomain interface {
	SendBackToken(ctx context.Context, email string) error
}

type User struct {
	email    string
	password string
}

type UserDomain struct {
	log    logger.ILogger
	driven userDriven.IUserDriven
}

func NewUserDomain(log logger.ILogger, userDriven userDriven.IUserDriven) IUserDomain {
	return &UserDomain{
		log:    log,
		driven: userDriven,
	}
}

func (u *UserDomain) SendBackToken(ctx context.Context, email string) error {
	u.log.Info(ctx, "Request received", logger.NewLogFieldInput("methodName", "SendBackToken"))
	token, err := middleware.CreateToken("user")
	if err != nil {
		return err
	}
	userLoginTimeout, err := time.ParseDuration(common.UserLoginTimeOut)
	if err != nil {
		u.log.Error(ctx, "login time not parsable "+err.Error())
		return err
	}

	if err := u.driven.SetTimeoutForUserRegister(ctx, email, token, userLoginTimeout); err != nil {
		u.log.Error(ctx, "Error in setting redis key "+err.Error())
		return err
	}
	header := metadata.New(map[string]string{"authentication": token})
	if err := grpc.SendHeader(ctx, header); err != nil {
		return err
	}
	return nil
}
