package app

import (
	"fmt"
	"net"
	"userService/internal/app/grpcApp"
	"userService/internal/config"
	"userService/internal/service"
	"userService/internal/storage"
	"userService/pkg/auth"
	psqlClient "userService/pkg/client/psql"
	"userService/pkg/hasher"
	"userService/pkg/logger"
)

func Run() error {
	logger := logger.NewLogger()

	cfg, err := config.NewConfig("main", "yaml", "./configs")
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	conn, err := psqlClient.NewPsqlClient(cfg.PGHost, cfg.PGDBName, cfg.PGUser, cfg.PGPassword, cfg.PGPort)
	if err != nil {
		logger.Error("failed to connect to pg", err)
		return fmt.Errorf("failed to connect to pg: %w", err)
	}

	userStorage := storage.NewUserStorage(conn, logger)

	authenticator := auth.NewAuth(cfg.JWTSignKey, logger)
	passwordHasher := hasher.NewPasswordHasher(cfg.HasherSalt)

	userService := service.NewAuthService(userStorage, authenticator, logger, passwordHasher)

	srv := grpcApp.New(logger, userService)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	logger.Info("gRPC user server is listening", "port", cfg.GRPCPort)

	if err := srv.GRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
