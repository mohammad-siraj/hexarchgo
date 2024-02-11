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
	RegisterServer(registerServerFunc func(grpc.ServiceRegistrar, interface{}), serverInstance interface{})
	RegisterFuncHandler(registerFuncHandlerFunc RegisterFuncHandlerFunc) error
	GetServerInstanceForRegister() *grpc.Server
	GetMuxInstanceForRegister() *runtime.ServeMux
	GetClientInstanceForRegister() *grpc.ClientConn
	StartGrpcServer(ctx context.Context, server IHttpClient) error
	ServerPort(serve net.Listener)
	ConnectClient(clientConnectionString string) error
}

type grpcGw struct {
	server *grpc.Server
	client *grpc.ClientConn
	mux    *runtime.ServeMux
}

type GrpcOptions struct {
	logger interface {
		GetGrpcUnaryInterceptor() grpc.UnaryServerInterceptor
	}
}

func NewGrpcGw(opts GrpcOptions, authMiddleware grpc.UnaryServerInterceptor) (IgrpcGw, error) {
	server := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			opts.logger.GetGrpcUnaryInterceptor(),
			authMiddleware,
		),
	)
	gwmux := runtime.NewServeMux()
	return &grpcGw{
		server: server,
		mux:    gwmux,
	}, nil
}

func (g *grpcGw) ServerPort(serve net.Listener) {
	defer g.server.Stop()
	if err := g.server.Serve(serve); err != nil {
		return
	}
}

func (g *grpcGw) RegisterServer(registerServerFunc func(grpc.ServiceRegistrar, interface{}), serverInstance interface{}) {
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

func (g *grpcGw) GetServerInstanceForRegister() *grpc.Server {
	return g.server
}

func (g *grpcGw) GetMuxInstanceForRegister() *runtime.ServeMux {
	return g.mux
}

func (g *grpcGw) GetClientInstanceForRegister() *grpc.ClientConn {
	return g.client
}

func (g *grpcGw) ConnectClient(clientConnectionString string) error {
	conn, err := grpc.DialContext(
		context.Background(),
		clientConnectionString,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return err
	}
	g.client = conn
	return nil
}
