package app

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"net"
	"userService/internal/config"
	grpcAuth "userService/internal/grpc"
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

	logger.Info("config", "config", cfg)

	conn, err := psqlClient.NewPsqlClient(cfg.PGHost, cfg.PGDBName, cfg.PGUser, cfg.PGPassword, cfg.PGPort)
	if err != nil {
		logger.Error("failed to connect to pg", err)
		return fmt.Errorf("failed to connect to pg: %w", err)
	}

	userStoage := storage.NewUserStorage(conn, logger)

	authenticator := auth.NewAuth(cfg.JWTSignKey, logger)
	passwordHasher := hasher.NewPasswordHasher(cfg.HasherSalt)

	userService := service.NewAuthService(userStoage, authenticator, logger, passwordHasher)

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Error("panic occurred", p)
			return nil
		}),
	}

	grpcServer := grpc.NewServer(grpc.ChainUnaryInterceptor(recovery.UnaryServerInterceptor(recoveryOpts...)))
	grpcAuth.Register(grpcServer, userService)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	logger.Info("gRPC user server is listening", "port", cfg.GRPCPort)

	if err := grpcServer.Serve(l); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}
