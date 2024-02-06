package http

import (
	"context"

	"github.com/gin-gonic/gin"
)

func (g *grpcGw) StartGrpcServer(ctx context.Context, server IHttpClient) error {
	server.NewSubGroup("v1/*{grpc_gateway}").Any("", gin.WrapH(g.mux))
	return nil
}
