package domain

import (
	"context"

	"time"

	"github.com/mohammad-siraj/hexarchgo/internal/common"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/database/cache"
	sqlDb "github.com/mohammad-siraj/hexarchgo/internal/libs/database/sql"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type IUserDomain interface {
	SendBackToken(ctx context.Context) error
}

type UserDomain struct {
	log   logger.ILogger
	cache cache.ICacheClient
	sqlDb sqlDb.ISqlDatabase
}

func NewUserDomain(log logger.ILogger, cache cache.ICacheClient, sqlDb sqlDb.ISqlDatabase) IUserDomain {
	return &UserDomain{
		log:   log,
		cache: cache,
		sqlDb: sqlDb,
	}
}

func (u *UserDomain) SendBackToken(ctx context.Context) error {
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

	if err := u.cache.Set(ctx, token, "value", userLoginTimeout); err != nil {
		u.log.Error(ctx, err.Error())
		return err
	}
	header := metadata.New(map[string]string{"authentication": token})
	if err := grpc.SendHeader(ctx, header); err != nil {
		return err
	}
	return nil
}
