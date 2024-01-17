package grpcApp

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime"
	"runtime/debug"
	"strings"
	grpcAuth "userService/internal/grpc"
	"userService/internal/service"
	"userService/pkg/logger"
)

type grpcApp struct {
	userService service.UserService
	GRPCServer  *grpc.Server
}

func New(logger logger.Logger, userService service.UserService) *grpcApp {
	logger.Info("starting gRPC user server")
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	logger.Info("recovering from panics")
	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			logger.Error("panic occurred", p)

			if runtimeErr, ok := p.(runtime.Error); ok && strings.Contains(runtimeErr.Error(), "nil pointer dereference") {
				logger.Error("nil pointer dereference occurred", "stack_trace", string(debug.Stack()))

				return status.Errorf(codes.Internal, "internal error: nil pointer dereference")
			}

			logger.Error("unhandled panic", "err", p)
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	logger.Info("creating gRPC server")
	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(logger), loggingOpts...),
	))

	logger.Info("registering gRPC server")
	grpcAuth.Register(gRPCServer, userService)

	return &grpcApp{
		userService: userService,
		GRPCServer:  gRPCServer,
	}
}

func InterceptorLogger(l logger.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...interface{}) {
		l.Info(msg, fields...)
	})
}
