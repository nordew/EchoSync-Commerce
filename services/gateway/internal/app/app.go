package app

import (
	v1 "gateway/internal/controller/http/v1"
	"gateway/pkg/logging"
	nordew "github.com/nordew/EchoSync-protos/gen/go/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func Run() error {
	logger := logging.NewLogger()

	logger.Info("starting gRPC user client")
	userServerConn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	userClient := nordew.NewUserClient(userServerConn)

	handler := v1.NewHandler(logger, userClient)

	app := handler.Init()

	if err := app.Listen(":8080"); err != nil {
		return err
	}
	logger.Info("started gRPC user client")

	return nil
}
