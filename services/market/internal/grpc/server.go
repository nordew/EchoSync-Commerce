package grpcStore

import (
	"context"
	"github.com/google/uuid"
<<<<<<< HEAD
	nordew "github.com/nordew/EchoSync-protos/gen/go/store"
=======
	grpcStore "github.com/nordew/EchoSync-protos/gen/go/store"
>>>>>>> marketService
	"google.golang.org/grpc"
	"marketService/internal/domain/entity"
	"marketService/internal/services"
	"marketService/pkg/logger"
	"time"
)

type grpcServer struct {
<<<<<<< HEAD
	nordew.UnimplementedStoreServiceServer
=======
	grpcStore.UnimplementedStoreServiceServer
>>>>>>> marketService

	storeService services.StoreService

	logger logger.Logger
}

func Register(s *grpc.Server, logger logger.Logger) {
<<<<<<< HEAD
	nordew.RegisterStoreServiceServer(s, &grpcServer{
=======
	grpcStore.RegisterStoreServiceServer(s, &grpcServer{
>>>>>>> marketService
		logger: logger,
	})
}

func NewStoreService(storeService services.StoreService, logger logger.Logger) *grpcServer {
	return &grpcServer{
		storeService: storeService,
		logger:       logger,
	}
}

<<<<<<< HEAD
func (s *grpcServer) CreateStore(ctx context.Context, req *nordew.CreateStoreRequest) (*nordew.Empty, error) {
=======
func (s *grpcServer) CreateStore(ctx context.Context, req *grpcStore.CreateStoreRequest) (*grpcStore.Empty, error) {
>>>>>>> marketService
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

	return &grpcStore.Empty{}, nil
}
