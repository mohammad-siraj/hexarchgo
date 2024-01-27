package http

import (
	"context"
	"net"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type RegisterServerFunc func(s grpc.ServiceRegistrar, server interface{})
type RegisterFuncHandlerFunc func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

type IgrpcGw interface {
	RegisterServer(registerServerFunc RegisterServerFunc, serverInstance interface{})
	RegisterFuncHandler(registerFuncHandlerFunc RegisterFuncHandlerFunc) error
}

type grpcGw struct {
	server *grpc.Server
	client *grpc.ClientConn
	mux    *runtime.ServeMux
}

func NewGrpcGw(ClientConnString string, MuxConnString string, opts ...interface{}) (IgrpcGw, error) {
	server := grpc.NewServer()

	clientConn, err := grpc.DialContext(
		context.Background(),
		//"0.0.0.0:8080"
		ClientConnString,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}
	gwmux := runtime.NewServeMux()
	return &grpcGw{
		server: server,
		client: clientConn,
		mux:    gwmux,
	}, nil
}

func (g *grpcGw) RegisterServer(registerServerFunc RegisterServerFunc, serverInstance interface{}) {
	registerServerFunc(g.server, serverInstance)

}

func (g *grpcGw) RegisterFuncHandler(registerFuncHandlerFunc RegisterFuncHandlerFunc) error {
	err := registerFuncHandlerFunc(context.Background(), g.mux, g.client)
	if err != nil {
		return err
	}
	return nil
}

func (g *grpcGw) GrpcServerServe(ServerConnString string) error {
	lis, err := net.Listen("tcp", ServerConnString)
	if err != nil {
		return err
	}
	if err := g.server.Serve(lis); err != nil {
		return err
	}
	return nil
}
