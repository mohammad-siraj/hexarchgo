package http

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
)

func (g *grpcGw) StartGrpcServer(ctx context.Context, Serverport string, HTTPServerPort string, ClientConnString string, server IHttpClient) error {
	// lis, err := net.Listen("tcp", Serverport)
	// if err != nil {
	// 	return err
	// }
	//go g.server.Serve(lis)
	fmt.Println(ClientConnString)
	// clientConn, err := grpc.DialContext(
	// 	ctx,
	// 	//"0.0.0.0:8080"
	// 	ClientConnString,
	// 	grpc.WithBlock(),
	// 	grpc.WithTransportCredentials(insecure.NewCredentials()),
	// )
	// if err != nil {
	// 	return err
	// }

	// g.client = clientConn
	server.NewSubGroup("v1/*{grpc_gateway}").Any("", gin.WrapH(g.mux))
	err := server.Run(HTTPServerPort)
	if err != nil {
		return err
	}
	return nil
}
