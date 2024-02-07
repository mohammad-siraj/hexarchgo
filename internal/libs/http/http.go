package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type HandlerFunc func(*context.Context)

type ginHandlerFunc func(*gin.Context)

// ServeHTTP implements http.Handler.
func (HandlerFunc) ServeHTTP(http.ResponseWriter, *http.Request) {
	panic("unimplemented")
}

type httpClient struct {
	engine *gin.Engine
	grpc   IgrpcGw
}

type subGroup struct {
	subGroup *gin.RouterGroup
}

type IHttpClient interface {
	Run(ConnString string) error
	Get(relativePath string, handlerFunction ...gin.HandlerFunc)
	Put(relativePath string, handlerFunction ...gin.HandlerFunc)
	Patch(relativePath string, handlerFunction ...gin.HandlerFunc)
	Delete(relativePath string, handlerFunction ...gin.HandlerFunc)
	Post(relativePath string, handlerFunction ...gin.HandlerFunc)
	NewSubGroup(path string, handleFunctions ...gin.HandlerFunc) IHttpClient
	Any(relativePath string, handlerFunction ...gin.HandlerFunc)
	GetGrpcServerInstanceForRegister() IgrpcGw
}

func NewHttpServer(isGrpcEnabled bool, logger interface {
	GetGrpcUnaryInterceptor() grpc.UnaryServerInterceptor
}) (IHttpClient, error) {
	client := &httpClient{
		engine: gin.New(),
	}
	if isGrpcEnabled {
		grpcConnection, err := NewGrpcGw(GrpcOptions{
			logger: logger,
		})
		if err != nil {
			return nil, err
		}
		client.grpc = grpcConnection
	}

	client.engine.Use(gin.Logger())
	return client, nil
}

// grpc changes
func (c *httpClient) GetGrpcServerInstanceForRegister() IgrpcGw {
	return c.grpc
}
func (c *subGroup) GetGrpcServerInstanceForRegister() IgrpcGw {
	return nil
}

// basic router
func (h *httpClient) NewSubGroup(path string, handleFunctions ...gin.HandlerFunc) IHttpClient {
	return &subGroup{
		subGroup: h.engine.Group(path, handleFunctions...),
	}
}

func (h *httpClient) Run(ConnString string) error {
	err := h.engine.Run(ConnString)
	if err != nil {
		return err
	}
	return nil
}

func (h *httpClient) Get(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.engine.GET(relativePath, handlerFunction...)
}

func (h *httpClient) Put(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.engine.PUT(relativePath, handlerFunction...)
}

func (h *httpClient) Post(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.engine.POST(relativePath, handlerFunction...)
}

func (h *httpClient) Patch(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.engine.PATCH(relativePath, handlerFunction...)
}

func (h *httpClient) Delete(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.engine.DELETE(relativePath, handlerFunction...)
}

// middlewear Compatibility
func (h *httpClient) Use(middlewareFunctions ...gin.HandlerFunc) {
	h.engine.Use(middlewareFunctions...)
}

func (h *httpClient) Any(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.engine.Any(relativePath, handlerFunction...)
}

// Subgroups
func (h *subGroup) NewSubGroup(path string, handleFunctions ...gin.HandlerFunc) IHttpClient {
	return &subGroup{
		subGroup: h.subGroup.Group(path, handleFunctions...),
	}
}

func (h *subGroup) Get(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.subGroup.GET(relativePath, handlerFunction...)
}

func (h *subGroup) Put(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.subGroup.PUT(relativePath, handlerFunction...)
}

func (h *subGroup) Post(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.subGroup.POST(relativePath, handlerFunction...)
}

func (h *subGroup) Patch(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.subGroup.PATCH(relativePath, handlerFunction...)
}

func (h *subGroup) Delete(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.subGroup.DELETE(relativePath, handlerFunction...)
}

// middlewear Compatibility
func (h *subGroup) Use(middlewareFunctions ...gin.HandlerFunc) {
	h.subGroup.Use(middlewareFunctions...)
}

func (h *subGroup) Any(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.subGroup.Any(relativePath, handlerFunction...)
}
func (h *subGroup) Run(ConnString string) error {
	err := errors.New("unimplemented")
	return err
}
