package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
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
}, ioWriter io.Writer) (IHttpClient, error) {
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
			ioWriter.Write([]byte(string(s) + "\n"))
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

func (h *httpClient) Any(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.engine.Any(relativePath, handlerFunction...)
}

// Subgroups
func (h *subGroup) NewSubGroup(path string, handleFunctions ...gin.HandlerFunc) IHttpClient {
	subGroup := &subGroup{
		subGroup: h.subGroup.Group(path, handleFunctions...),
	}
	return subGroup
}
func (h *subGroup) Use(middleware http.HandlerFunc) {
	panic("unimplemented")
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

func (h *subGroup) Any(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.subGroup.Any(relativePath, handlerFunction...)
}

func (h *subGroup) Run(ConnString string) error {
	err := errors.New("unimplemented")
	return err
}
