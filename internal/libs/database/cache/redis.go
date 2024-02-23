package cache

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	dbInterface "github.com/mohammad-siraj/hexarchgo/internal/libs/database"
	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	client *redis.Client
}

type ICacheClient interface {
	dbInterface.IDatabase
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	Exit(ctx context.Context, errChannel chan error)
}

func NewCacheClient(address, password string, db int) ICacheClient {
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
		fmt.Println(output)
		return "", err
	}
	return output, nil
}

func (r *redisClient) queryStringParser(ctx context.Context, queryString string) (string, error) {
	queryComponents := strings.Split(queryString, " ")
	switch strings.ToUpper(queryComponents[0]) {
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

func (r *redisClient) Set(ctx context.Context, key string, value interface{}, timeOut time.Duration) error {
	err := r.client.Set(context.Background(), key, value, timeOut).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *redisClient) Get(ctx context.Context, key string) (string, error) {
	result := r.client.Get(ctx, key)
	if result.Err() != nil {
		return "", result.Err()
	}
	return result.Val(), nil
}

func (r *redisClient) Del(ctx context.Context, keys ...string) error {
	result := r.client.Del(context.Background(), keys...)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}

func (r *redisClient) Exit(ctx context.Context, errChannel chan error) {
	if err := r.client.Conn().Close(); err != nil {
		errChannel <- err
	}
}
