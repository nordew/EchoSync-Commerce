package grpcStore

import (
	"context"
	"github.com/google/uuid"
	nordew "github.com/nordew/EchoSync-protos/gen/go/store"
	"google.golang.org/grpc"
	"marketService/internal/domain/entity"
	"marketService/internal/services"
	"marketService/pkg/logger"
	"time"
)

type grpcServer struct {
	nordew.UnimplementedStoreServiceServer

	storeService services.StoreService

	logger logger.Logger
}

func Register(s *grpc.Server, logger logger.Logger) {
	nordew.RegisterStoreServiceServer(s, &grpcServer{
		logger: logger,
	})
}

func NewStoreService(storeService services.StoreService, logger logger.Logger) *grpcServer {
	return &grpcServer{
		storeService: storeService,
		logger:       logger,
	}
}

func (s *grpcServer) CreateStore(ctx context.Context, req *nordew.CreateStoreRequest) (*nordew.Empty, error) {
	const op = "grpcServer.CreateStore"

	parsedUUID, err := uuid.Parse(req.OwnerId)
	if err != nil {
		s.logger.Error(op, "failed to parse UUID", err.Error())
		return nil, err
	}

	store := &entity.Store{
		Name:        req.GetName(),
		OwnerUserID: parsedUUID,
		CreatedAt:   time.Now(),
	}

	err = s.storeService.Create(ctx, store)
	if err != nil {
		s.logger.Error(op, "failed to create store", err.Error())
		return nil, err
	}

	return &nordew.Empty{}, nil
}
