package grpcStore

import (
	"context"
	"errors"
	"github.com/google/uuid"
	grpcStore "github.com/nordew/EchoSync-protos/gen/go/store"
	"google.golang.org/grpc"
	"marketService/internal/services"
	"marketService/pkg/logger"
)

var (
	ErrInvalidUUId = errors.New("invalid uuid")
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
		s.logger.Error(op, err.Error())
		return nil, ErrInvalidUUId
	}

	err = s.storeService.Create(context.Background(), req.GetName(), parsedUUID)
	if err != nil {
		s.logger.Error(op, "failed to create store", err.Error())
		return nil, err
	}

	return &grpcStore.Empty{}, nil
}
