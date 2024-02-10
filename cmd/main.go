package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/gin-gonic/gin"
	httpServer "github.com/mohammad-siraj/hexarchgo/internal/libs/http"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
	"github.com/mohammad-siraj/hexarchgo/internal/ports"
)

func main() {
	HTTPServerPort := ":8090"
	GRPCServerPort := ":8081"
	isGrpcEnabled := true
	ctx := context.Background()
	StartServer(ctx, isGrpcEnabled, GRPCServerPort, HTTPServerPort)

}

func StartServer(ctx context.Context, isGrpcEnabled bool, GRPCServerPort string, HTTPServerPort string) {
	loggerConfig := logger.NewlogConfigOptions(false)

	//logger configs
	loggerConfig.WithFilename("log/app.log")
	loggerConfig.WithIsCompressed(true)
	loggerConfig.WithIsLocalTime(true)
	loggerConfig.WithMaxAge(1)
	loggerConfig.WithMaxBackUp(2)

	//logger
	loggerInstance := logger.NewLogger(loggerConfig)
	loggerInstance.Info(ctx, "Starting server ...\n")

	serverHttp, err := httpServer.NewHttpServer(isGrpcEnabled, loggerInstance, loggerInstance.GetIoWriter())
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}
	if isGrpcEnabled {
		ports.RegisterServicesToGrpcServer(serverHttp)
		grpcServer := serverHttp.GetGrpcServerInstanceForRegister()
		lis, err := net.Listen("tcp", GRPCServerPort)
		if err != nil {
			log.Fatal(err)
		}
		GRPCServerPortIp := fmt.Sprintf("0.0.0.0%s", GRPCServerPort)
		go grpcServer.ServerPort(lis)
		grpcServer.ConnectClient(GRPCServerPortIp)
		ports.RegisterServiceHandlersToGrpcServer(ctx, grpcServer)
		grpcServer.StartGrpcServer(ctx, serverHttp)
	}
	// serverHttp.Use(func(w http.ResponseWriter, r *http.Request) {
	// 	loggerInstance.RequestLog(r)
	// })

	serverHttp.Get("/token", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	}))

	subGroup := serverHttp.NewSubGroup("/auth", func(ctx *gin.Context) {

	})
	{
		subGroup.Get("/test", gin.WrapF(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World!"))
		}))
	}

	err = serverHttp.Run(HTTPServerPort)
	if err != nil {
		log.Fatal(err)
	}
}
