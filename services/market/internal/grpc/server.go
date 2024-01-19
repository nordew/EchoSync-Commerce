package grpcStore

import (
	"errors"
	grpcStore "github.com/nordew/EchoSync-protos/gen/go/store"
	"google.golang.org/grpc"
	"marketService/internal/services"
	"marketService/pkg/logger"
)

var (
	ErrInvalidUUID = errors.New("invalid uuid")
)

type grpcServer struct {
	grpcStore.UnimplementedStoreServiceServer
	grpcStore.UnimplementedProductServiceServer

	storeService   services.StoreService
	productService services.ProductService

	logger logger.Logger
}

func Register(s *grpc.Server, storeService services.StoreService, productService services.ProductService, logger logger.Logger) {
	grpcStore.RegisterStoreServiceServer(s, &grpcServer{
		storeService:   storeService,
		productService: productService,
		logger:         logger,
	})
}
