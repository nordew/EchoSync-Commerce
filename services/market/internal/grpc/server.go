package grpcStore

import (
	"context"
	"github.com/google/uuid"
	grpcStore "github.com/nordew/EchoSync-protos/gen/go/store"
	"google.golang.org/grpc"
	"marketService/internal/services"
	"marketService/pkg/logger"
)

type grpcServer struct {
	grpcStore.UnimplementedStoreServiceServer

	storeService services.StoreService

	logger logger.Logger
}

func Register(s *grpc.Server, storeService services.StoreService, logger logger.Logger) {
	grpcStore.RegisterStoreServiceServer(s, &grpcServer{
		storeService: storeService,
		logger:       logger,
	})
}

func (s *grpcServer) CreateStore(ctx context.Context, req *grpcStore.CreateStoreRequest) (*grpcStore.Empty, error) {
	const op = "grpcServer.CreateStore"

	parsedUUID, err := uuid.Parse(req.GetOwnerId())
	if err != nil {
		s.logger.Error(op, "failed to parse UUID", err.Error())
		return nil, err
	}

	err = s.storeService.Create(context.Background(), req.GetName(), parsedUUID)
	if err != nil {
		s.logger.Error(op, "failed to create store", err.Error())
		return nil, err
	}

	return &grpcStore.Empty{}, nil
}
