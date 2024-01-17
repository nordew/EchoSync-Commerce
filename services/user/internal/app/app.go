package app

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	conn, err := psqlClient.NewPsqlClient(cfg.PGHost, cfg.PGDBName, cfg.PGUser, cfg.PGPassword, cfg.PGPort)
	if err != nil {
		logger.Error("failed to connect to pg", err)
		return fmt.Errorf("failed to connect to pg: %w", err)
	}

	userStorage := storage.NewUserStorage(conn, logger)

	authenticator := auth.NewAuth(cfg.JWTSignKey, logger)
	passwordHasher := hasher.NewPasswordHasher(cfg.HasherSalt)

	userService := service.NewAuthService(userStorage, authenticator, logger, passwordHasher)

	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Error("panic occurred", p)

			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(logger), loggingOpts...),
	))

	grpcAuth.Register(gRPCServer, userService)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GRPCPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	logger.Info("gRPC user server is listening", "port", cfg.GRPCPort)

	if err := gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("failed to serve: %w", err)
	}

	return nil
}

func InterceptorLogger(l logger.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Info(msg, fields...)
	})
}
