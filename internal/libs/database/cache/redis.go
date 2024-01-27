package cache

import (
	"context"
	"errors"
	"strings"

	dbInterface "github.com/mohammad-siraj/hexarchgo/internal/libs/database"
	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	client *redis.Client
}

func NewCacheClient(address, password string, db int) dbInterface.IDatabase {
	return &redisClient{
		client: redis.NewClient(&redis.Options{
			Addr:     address,
			Password: password,
			DB:       db,
		}),
	}
}

func (r *redisClient) ExecWithContext(ctx context.Context, queryString string, opt ...interface{}) (string, error) {
	output, err := r.queryStringParser(ctx, queryString)
	if err != nil {
		return "", err
	}
	return output, nil
}

func (r *redisClient) queryStringParser(ctx context.Context, queryString string) (string, error) {
	queryComponents := strings.Split(queryString, " ")
	switch strings.ToLower(queryComponents[0]) {
	case "SET":
		{
			if len(queryComponents) != 3 {
				return "", errors.New("invalid query")
			}
			key := queryComponents[1]
			value := queryComponents[2]
			output := r.client.Set(ctx, key, value, 0)
			if output.Err() != nil {
				return "", output.Err()
			}
			return "OK", nil
		}
	case "GET":
		{
			if len(queryComponents) != 2 {
				return "", errors.New("invalid query")
			}
			key := queryComponents[1]
			output := r.client.Get(ctx, key)
			if output.Err() != nil {
				return "", output.Err()
			}
			return output.Val(), nil
		}
	}
	return "", nil
}