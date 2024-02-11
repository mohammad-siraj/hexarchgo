package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/middleware"
	"google.golang.org/grpc"
)

type HandlerFunc func(*context.Context)

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
	Get(relativePath string, handlerFunction ...http.HandlerFunc)
	Put(relativePath string, handlerFunction ...http.HandlerFunc)
	Patch(relativePath string, handlerFunction ...http.HandlerFunc)
	Delete(relativePath string, handlerFunction ...http.HandlerFunc)
	Post(relativePath string, handlerFunction ...http.HandlerFunc)
	NewSubGroup(path string, handleFunctions ...http.HandlerFunc) IHttpClient
	Any(relativePath string, handlerFunction ...gin.HandlerFunc)
	Use(handlerFunction ...http.HandlerFunc)
	GetGrpcServerInstanceForRegister() IgrpcGw
}

func NewHttpServer(isGrpcEnabled bool, logger interface {
	GetGrpcUnaryInterceptor() grpc.UnaryServerInterceptor
}, ioWriter io.Writer) (IHttpClient, error) {
	client := &httpClient{
		engine: gin.New(),
	}
	if isGrpcEnabled {
		grpcConnection, err := NewGrpcGw(GrpcOptions{
			logger: logger,
		}, middleware.GrpcAuthMiddleware(context.Background()))
		if err != nil {
			return nil, err
		}
		client.grpc = grpcConnection
	}

	if ioWriter == nil {
		client.engine.Use(gin.Logger())
	} else {
		client.engine.Use(jsonLoggerMiddleware(ioWriter))
	}
	return client, nil
}

func jsonLoggerMiddleware(ioWriter io.Writer) gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})

			log["status_code"] = params.StatusCode
			log["path"] = params.Path
			log["method"] = params.Method
			log["start_time"] = params.TimeStamp.Format("2006/01/02 - 15:04:05")
			log["remote_addr"] = params.ClientIP
			log["response_time"] = params.Latency.String()

			s, _ := json.Marshal(log)
			if _, err := ioWriter.Write([]byte(string(s) + "\n")); err != nil {
				return `{"error":"log format"}`
			}
			return string(s) + "\n"
		},
	)
}

// grpc changes
func (c *httpClient) GetGrpcServerInstanceForRegister() IgrpcGw {
	return c.grpc
}
func (c *subGroup) GetGrpcServerInstanceForRegister() IgrpcGw {
	return nil
}

// basic router
func (h *httpClient) NewSubGroup(path string, handleFunctions ...http.HandlerFunc) IHttpClient {
	ginHandleFunction := make([]gin.HandlerFunc, len(handleFunctions))
	for i, v := range handleFunctions {
		ginHandleFunction[i] = gin.WrapF(v)
	}
	return &subGroup{
		subGroup: h.engine.Group(path, ginHandleFunction...),
	}
}

func (h *httpClient) Run(ConnString string) error {
	err := h.engine.Run(ConnString)
	if err != nil {
		return err
	}
	return nil
}

func (h *httpClient) Get(relativePath string, handlerFunction ...http.HandlerFunc) {
	ginHandleFunction := make([]gin.HandlerFunc, len(handlerFunction))
	for i, v := range handlerFunction {
		ginHandleFunction[i] = gin.WrapF(v)
	}
	h.engine.GET(relativePath, ginHandleFunction...)
}

func (h *httpClient) Put(relativePath string, handlerFunction ...http.HandlerFunc) {
	ginHandleFunction := make([]gin.HandlerFunc, len(handlerFunction))
	for i, v := range handlerFunction {
		ginHandleFunction[i] = gin.WrapF(v)
	}
	h.engine.PUT(relativePath, ginHandleFunction...)
}

func (h *httpClient) Post(relativePath string, handlerFunction ...http.HandlerFunc) {
	ginHandleFunction := make([]gin.HandlerFunc, len(handlerFunction))
	for i, v := range handlerFunction {
		ginHandleFunction[i] = gin.WrapF(v)
	}
	h.engine.POST(relativePath, ginHandleFunction...)
}

func (h *httpClient) Patch(relativePath string, handlerFunction ...http.HandlerFunc) {
	ginHandleFunction := make([]gin.HandlerFunc, len(handlerFunction))
	for i, v := range handlerFunction {
		ginHandleFunction[i] = gin.WrapF(v)
	}
	h.engine.PATCH(relativePath, ginHandleFunction...)
}

func (h *httpClient) Delete(relativePath string, handlerFunction ...http.HandlerFunc) {
	ginHandleFunction := make([]gin.HandlerFunc, len(handlerFunction))
	for i, v := range handlerFunction {
		ginHandleFunction[i] = gin.WrapF(v)
	}
	h.engine.DELETE(relativePath, ginHandleFunction...)
}

func (h *httpClient) Use(handlerFunction ...http.HandlerFunc) {
	ginHandleFunction := make([]gin.HandlerFunc, len(handlerFunction))
	for i, v := range handlerFunction {
		ginHandleFunction[i] = gin.WrapF(v)
	}
	h.engine.Use(ginHandleFunction...)
}

func (h *httpClient) Any(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.engine.Any(relativePath, handlerFunction...)
}

// Subgroups
func (h *subGroup) NewSubGroup(path string, handleFunctions ...http.HandlerFunc) IHttpClient {

	ginHandleFunction := make([]gin.HandlerFunc, len(handleFunctions))
	for i, v := range handleFunctions {
		ginHandleFunction[i] = gin.WrapF(v)
	}
	subGroup := &subGroup{
		subGroup: h.subGroup.Group(path, ginHandleFunction...),
	}
	return subGroup
}

// func (h *subGroup) Use(middleware http.HandlerFunc) {
// 	panic("unimplemented")
// }

func (h *subGroup) Get(relativePath string, handlerFunction ...http.HandlerFunc) {
	ginHandleFunction := make([]gin.HandlerFunc, len(handlerFunction))
	for i, v := range handlerFunction {
		ginHandleFunction[i] = gin.WrapF(v)
	}
	h.subGroup.GET(relativePath, ginHandleFunction...)
}

func (h *subGroup) Put(relativePath string, handlerFunction ...http.HandlerFunc) {
	ginHandleFunction := make([]gin.HandlerFunc, len(handlerFunction))
	for i, v := range handlerFunction {
		ginHandleFunction[i] = gin.WrapF(v)
	}
	h.subGroup.PUT(relativePath, ginHandleFunction...)
}

func (h *subGroup) Post(relativePath string, handlerFunction ...http.HandlerFunc) {
	ginHandleFunction := make([]gin.HandlerFunc, len(handlerFunction))
	for i, v := range handlerFunction {
		ginHandleFunction[i] = gin.WrapF(v)
	}
	h.subGroup.POST(relativePath, ginHandleFunction...)
}

func (h *subGroup) Patch(relativePath string, handlerFunction ...http.HandlerFunc) {
	ginHandleFunction := make([]gin.HandlerFunc, len(handlerFunction))
	for i, v := range handlerFunction {
		ginHandleFunction[i] = gin.WrapF(v)
	}
	h.subGroup.PATCH(relativePath, ginHandleFunction...)
}

func (h *subGroup) Delete(relativePath string, handlerFunction ...http.HandlerFunc) {
	ginHandleFunction := make([]gin.HandlerFunc, len(handlerFunction))
	for i, v := range handlerFunction {
		ginHandleFunction[i] = gin.WrapF(v)
	}
	h.subGroup.DELETE(relativePath, ginHandleFunction...)
}

func (h *subGroup) Use(handlerFunction ...http.HandlerFunc) {
	ginHandleFunction := make([]gin.HandlerFunc, len(handlerFunction))
	for i, v := range handlerFunction {
		ginHandleFunction[i] = gin.WrapF(v)
	}
	h.subGroup.Use(ginHandleFunction...)
}

func (h *subGroup) Any(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.subGroup.Any(relativePath, handlerFunction...)
}

func (h *subGroup) Run(ConnString string) error {
	err := errors.New("unimplemented")
	return err
}
