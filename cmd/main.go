package main

import (
	"context"
	"log"

	"github.com/mohammad-siraj/hexarchgo/internal/libs/database/cache"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/database/sql"
	httpServer "github.com/mohammad-siraj/hexarchgo/internal/libs/http"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/logger"
	"github.com/mohammad-siraj/hexarchgo/internal/libs/middleware"
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
	loggerConfig.WithMaxBackUp(1)

	//logger
	loggerInstance := logger.NewLogger(loggerConfig)
	loggerInstance.Info(ctx, "Starting server ...\n")

	//Http client
	serverHttp, err := httpServer.NewHttpServer(loggerInstance.GetIoWriter(), httpServer.NewGrpcOptions(loggerInstance, middleware.GrpcAuthMiddleware(ctx)))
	if err != nil {
		loggerInstance.Error(ctx, err.Error())
	}

	//cache client
	cacheClient := cache.NewCacheClient("localhost:6379", "", 1)

	//sql database client
	sqlClient, err := sql.NewDatabase("postgresql://postgres:postgres@localhost:5432/mainserver?sslmode=disable")
	if err != nil {
		loggerInstance.Error(ctx, "error while initiating sql client "+err.Error())
	}

	_, err = sqlClient.ExecWithContext(ctx, "INSERT INTO users (id, username,email,password,created_at) VALUES ('1', 'John Doe','John Doe','John Doe','2024-02-14T18:27:58+00:00')")
	if err != nil {
		loggerInstance.Error(ctx, "error while initiating sql client "+err.Error())
	}

	porter := ports.NewPorter(serverHttp, loggerInstance, cacheClient)
	if isGrpcEnabled {
		porter.SetUpGrpcGateWay(ctx, GRPCServerPort)
	}

	porter.RegisterRequestHandlers()

	//event handler test
	// BrokerConfig := make([]string, 0)
	// BrokerConfig = append(BrokerConfig, "0.0.0.0:9092")
	// //var Offset int64 = -1
	// topic := "test_topic"
	// eventbroker, err := eventbroker.NewEventBroker(BrokerConfig)
	// if err != nil {
	// 	loggerInstance.Error(ctx, "error while initiating event broker client "+err.Error())
	// }
	// partition, offset, err := eventbroker.SyncSendMessageToTopic("hello there siraj\n", topic, false)
	// if err != nil {
	// 	loggerInstance.Error(ctx, "error while send message to topic client "+err.Error())
	// }
	// ch, err := eventbroker.ConsumeTopic(topic, 0)
	// if err != nil {
	// 	loggerInstance.Error(ctx, "error while consumer initialization client "+err.Error())
	// }
	// go func() {
	// 	for msg := range ch {
	// 		fmt.Println("The value is here ", msg)
	// 		loggerInstance.Info(ctx, "value is "+msg)
	// 	}
	// }()
	err = serverHttp.Run(HTTPServerPort)
	if err != nil {
		log.Fatal(err)
	}
}
