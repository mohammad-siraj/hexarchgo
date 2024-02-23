package driven

import (
	"context"
	"time"

	"github.com/mohammad-siraj/hexarchgo/internal/libs/database/cache"
	sqlDb "github.com/mohammad-siraj/hexarchgo/internal/libs/database/sql"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
)

type IUserDriven interface {
	SetTimeoutForUserRegister(ctx context.Context, email, token string, userLoginTimeout time.Duration) error
}

type userDriven struct {
	log   logger.ILogger
	cache cache.ICacheClient
	sqlDb sqlDb.ISqlDatabase
}

func NewUserDriven(log logger.ILogger, cache cache.ICacheClient, sqlDb sqlDb.ISqlDatabase) IUserDriven {
	return &userDriven{
		log:   log,
		cache: cache,
		sqlDb: sqlDb,
	}
}

func (u *userDriven) SetTimeoutForUserRegister(ctx context.Context, email, token string, userLoginTimeout time.Duration) error {
	if err := u.cache.Set(ctx, email, token, userLoginTimeout); err != nil {
		u.log.Error(ctx, "Error in setting redis key "+err.Error())
		return err
	}
	return nil
}
