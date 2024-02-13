package domain

import (
	"context"
	"os"
	"time"

	"github.com/mohammad-siraj/hexarchgo/internal/libs/database/cache"
	sqlDb "github.com/mohammad-siraj/hexarchgo/internal/libs/database/sql"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	userLoginTime = os.Getenv("USER_LOGIN_TIME") //nolint: gosec
)

type IUserDomain interface {
	SendBackToken(ctx context.Context, token string) error
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

func (u *UserDomain) SendBackToken(ctx context.Context, token string) error {
	u.log.Info(ctx, "Request received", logger.NewLogFieldInput("methodName", "SendBackToken"))
	token, err := middleware.CreateToken("user")
	if err != nil {
		return err
	}
	if err := u.cache.Set(ctx, token, "value", 1000*time.Millisecond); err != nil {
		u.log.Error(ctx, err.Error())
		return err
	}
	header := metadata.New(map[string]string{"authentication": token})
	if err := grpc.SendHeader(ctx, header); err != nil {
		return err
	}
	return nil
}
