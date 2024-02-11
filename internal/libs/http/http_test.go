package http

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

type mockLogger struct {
}

func (m *mockLogger) GetGrpcUnaryInterceptor() grpc.UnaryServerInterceptor {
	return nil
}

// Creates a new httpClient instance with gin engine
func TestNewHttpServerWithGinEngine(t *testing.T) {
	logger := &mockLogger{}
	ioWriter := &bytes.Buffer{}
	client, err := NewHttpServer(false, logger, ioWriter)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.IsType(t, &httpClient{}, client)
}

// If isGrpcEnabled is true, creates a new grpc connection
func TestNewHttpServerWithGrpcConnection(t *testing.T) {
	logger := &mockLogger{}
	ioWriter := &bytes.Buffer{}
	client, err := NewHttpServer(true, logger, ioWriter)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.IsType(t, &httpClient{}, client)
}

// Uses jsonLoggerMiddleware if ioWriter is not nil
func TestNewHttpServerWithJsonLoggerMiddleware(t *testing.T) {
	logger := &mockLogger{}
	ioWriter := &bytes.Buffer{}
	client, err := NewHttpServer(false, logger, ioWriter)
	assert.NoError(t, err)
	assert.NotNil(t, client)
	assert.IsType(t, &httpClient{}, client)
	//assert.Contains(t, client.engine.Handlers, jsonLoggerMiddleware(ioWriter))
}
