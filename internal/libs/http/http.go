package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HandlerFunc func(*context.Context)

// ServeHTTP implements http.Handler.
func (HandlerFunc) ServeHTTP(http.ResponseWriter, *http.Request) {
	panic("unimplemented")
}

type HttpClient struct {
	engine *gin.Engine
}

type SubGroup struct {
	subGroup gin.RouterGroup
}

type IHttpClient interface {
	Get(relativePath string, handlerFunction ...gin.HandlerFunc)
	Put(relativePath string, handlerFunction ...gin.HandlerFunc)
	Patch(relativePath string, handlerFunction ...gin.HandlerFunc)
	Delete(relativePath string, handlerFunction ...gin.HandlerFunc)
	NewSubGroup(path string, handleFunctions ...gin.HandlerFunc) IHttpClient
	Any(relativePath string, handlerFunction HandlerFunc)
}

func NewHttpClient() IHttpClient {
	client := &HttpClient{
		engine: gin.New(),
	}
	client.engine.Use(gin.Logger())
	return client
}

// basic router
func (h *HttpClient) NewSubGroup(path string, handleFunctions ...gin.HandlerFunc) IHttpClient {
	return &SubGroup{
		subGroup: *h.engine.Group(path, handleFunctions...),
	}
}

func (h *HttpClient) Get(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.engine.GET(relativePath, handlerFunction...)
}

func (h *HttpClient) Put(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.engine.PUT(relativePath, handlerFunction...)
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

func (h *HttpClient) Any(relativePath string, handlerFunction HandlerFunc) {
	h.engine.Any(relativePath, gin.WrapH(handlerFunction))
}

// Subgroups
func (h *SubGroup) NewSubGroup(path string, handleFunctions ...gin.HandlerFunc) IHttpClient {
	return &SubGroup{
		subGroup: *h.subGroup.Group(path, handleFunctions...),
	}
}

func (h *SubGroup) Get(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.subGroup.GET(relativePath, handlerFunction...)
}

func (h *SubGroup) Put(relativePath string, handlerFunction ...gin.HandlerFunc) {
	h.subGroup.PUT(relativePath, handlerFunction...)
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

func (h *SubGroup) Any(relativePath string, handlerFunction HandlerFunc) {
	h.subGroup.Any(relativePath, gin.WrapH(handlerFunction))
}
