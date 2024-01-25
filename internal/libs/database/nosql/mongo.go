package nosql

import (
	"context"
	"errors"

	dbInterface "github.com/mohammad-siraj/hexarchgo/internal/libs/database"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoClient struct {
	client *mongo.Client
}

func NewMongoClient(connectionString string) (dbInterface.IDatabase, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectionString))
	if err != nil {
		return nil, err
	}
	return &mongoClient{
		client: client,
	}, nil
}

func (m *mongoClient) ExecWithContext(ctx context.Context, queryString string, opt ...interface{}) (string, error) {
	if len(opt) < 1 {
		return "", errors.New("not enough params provided")
	}

	return "OK", nil
}
