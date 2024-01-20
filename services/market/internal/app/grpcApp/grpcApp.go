package grpcApp

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"marketService/internal/services"
	"marketService/pkg/logger"
	"runtime"
	"runtime/debug"
	"strings"

	grpcStore "marketService/internal/grpc"
)

type grpcApp struct {
	storeService   services.StoreService
	productService services.ProductService
	GRPCServer     *grpc.Server
	logger         logger.Logger
}

func New(storeService services.StoreService, productService services.ProductService, logger logger.Logger) *grpcApp {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

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

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(logger), loggingOpts...),
	))

	grpcStore.Register(gRPCServer, storeService, productService, logger)

	return &grpcApp{
		storeService:   storeService,
		productService: productService,
		GRPCServer:     gRPCServer,
		logger:         logger,
	}
}

func InterceptorLogger(l logger.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...interface{}) {
		l.Info(msg, fields...)
	})
}
