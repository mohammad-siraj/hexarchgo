package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(*context.Context)

type ginHandlerFunc func(*gin.Context)

// ServeHTTP implements http.Handler.
func (HandlerFunc) ServeHTTP(http.ResponseWriter, *http.Request) {
	panic("unimplemented")
}

type HttpClient struct {
	engine *gin.Engine
}

type SubGroup struct {
	subGroup *gin.RouterGroup
}

type IHttpClient interface {
	Run(ConnString string) error
	Get(relativePath string, handlerFunction ...gin.HandlerFunc)
	Put(relativePath string, handlerFunction ...gin.HandlerFunc)
	Patch(relativePath string, handlerFunction ...gin.HandlerFunc)
	Delete(relativePath string, handlerFunction ...gin.HandlerFunc)
	NewSubGroup(path string, handleFunctions ...gin.HandlerFunc) IHttpClient
	Any(relativePath string, handlerFunction ...gin.HandlerFunc)
}

func NewHttpClient() (IHttpClient, error) {
	client := &HttpClient{
		engine: gin.New(),
	}
	client.engine.Use(gin.Logger())
	return client, nil
}

// basic router
func (h *HttpClient) NewSubGroup(path string, handleFunctions ...gin.HandlerFunc) IHttpClient {
	return &SubGroup{
		subGroup: h.engine.Group(path, handleFunctions...),
	}
}

func (h *HttpClient) Run(ConnString string) error {
	err := h.engine.Run(ConnString)
	if err != nil {
		return err
	}
	return nil
}

func (h *HttpClient) Get(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.engine.GET(relativePath, handlerFunction...)
}

func (h *HttpClient) Put(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.engine.PUT(relativePath, handlerFunction...)
}

func (h *HttpClient) Post(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.engine.POST(relativePath, handlerFunction...)
}

func (h *HttpClient) Patch(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.engine.PATCH(relativePath, handlerFunction...)
}

func (h *HttpClient) Delete(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.engine.DELETE(relativePath, handlerFunction...)
}

// middlewear Compatibility
func (h *HttpClient) Use(middlewareFunctions ...gin.HandlerFunc) {
	h.engine.Use(middlewareFunctions...)
}

func (h *HttpClient) Any(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.engine.Any(relativePath, handlerFunction...)
}

// Subgroups
func (h *SubGroup) NewSubGroup(path string, handleFunctions ...gin.HandlerFunc) IHttpClient {
	return &SubGroup{
		subGroup: h.subGroup.Group(path, handleFunctions...),
	}
}

func (h *SubGroup) Get(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.subGroup.GET(relativePath, handlerFunction...)
}

func (h *SubGroup) Put(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.subGroup.PUT(relativePath, handlerFunction...)
}

func (h *SubGroup) Post(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.subGroup.POST(relativePath, handlerFunction...)
}

func (h *SubGroup) Patch(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.subGroup.PATCH(relativePath, handlerFunction...)
}

func (h *SubGroup) Delete(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.subGroup.DELETE(relativePath, handlerFunction...)
}

// middlewear Compatibility
func (h *SubGroup) Use(middlewareFunctions ...gin.HandlerFunc) {
	h.subGroup.Use(middlewareFunctions...)
}

func (h *SubGroup) Any(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.subGroup.Any(relativePath, handlerFunction...)
}
func (h *SubGroup) Run(ConnString string) error {
	err := errors.New("unimplemented")
	return err
}
