package app

import (
	"fmt"
	"marketService/internal/app/grpcApp"
	"marketService/internal/config"
	"marketService/internal/services"
	"marketService/internal/storage"
	psqlClient "marketService/pkg/client/psql"
	"marketService/pkg/logger"
	"net"
)

func Run() error {
	logging := logger.NewLogger()

	cfg, err := config.NewConfig("main", "yaml", "configs")
	if err != nil {
		logging.Error("failed to load config", err.Error())
		return err
	}

	conn, err := psqlClient.NewPsqlClient(cfg.PGHost, cfg.PGDBName, cfg.PGUser, cfg.PGPassword, cfg.PGPort)
	if err != nil {
		logging.Error("failed to create psql client", err.Error())
		return err
	}

	storeStorage := storage.NewStoreStorage(conn, logging)
	storeService := services.NewStoreService(storeStorage, logging)

	grpcSrv := grpcApp.New(storeService, logging)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		logging.Error("failed to listen", err.Error())
		return err
	}

	logging.Info("gRPC market server is listening", "port", cfg.GRPCPort)

	if err := grpcSrv.GRPCServer.Serve(lis); err != nil {
		logging.Error("failed to serve", err.Error())
		return err
	}

	return nil
}
