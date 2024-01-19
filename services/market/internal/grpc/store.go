package grpcStore

import (
	"context"
	"github.com/google/uuid"
	grpcStore "github.com/nordew/EchoSync-protos/gen/go/store"
)

func (s *grpcServer) CreateStore(ctx context.Context, req *grpcStore.CreateStoreRequest) (*grpcStore.Empty, error) {
	const op = "grpcServer.CreateStore"

	parsedUUID, err := uuid.Parse(req.GetOwnerId())
	if err != nil {
		s.logger.Error(op, err.Error())
		return nil, ErrInvalidUUID
	}

	err = s.storeService.Create(context.Background(), req.GetName(), parsedUUID)
	if err != nil {
		s.logger.Error(op, "failed to create store", err.Error())
		return nil, err
	}

	return &grpcStore.Empty{}, nil
}
