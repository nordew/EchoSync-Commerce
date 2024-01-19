package app

import (
	"gateway/internal/config"
	"gateway/pkg/auth"
	"gateway/pkg/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	v1 "gateway/internal/controller/http/v1"
	grpcStore "github.com/nordew/EchoSync-protos/gen/go/store"
	grpcUser "github.com/nordew/EchoSync-protos/gen/go/user"
)

func Run() error {
	logger := logging.NewLogger()

	cfg, err := config.NewConfig("main", "yaml", "configs")
	if err != nil {
		logger.Error("failed to load config", err)
		return err
	}

	logger.Info("starting gRPC user client")
	userServerConn, err := grpc.Dial(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	storeClientConn, err := grpc.Dial(":50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	authenticator := auth.NewAuth(cfg.JWTSignKey, logger)

	userClient := grpcUser.NewUserClient(userServerConn)
	storeClient := grpcStore.NewStoreServiceClient(storeClientConn)
	productClient := grpcStore.NewProductServiceClient(storeClientConn)

	handler := v1.NewHandler(logger, userClient, storeClient, productClient, authenticator)

	app := handler.Init()

	if err := app.Listen(":8080"); err != nil {
		return err
	}
	logger.Info("started gRPC user client")

	return nil
}
