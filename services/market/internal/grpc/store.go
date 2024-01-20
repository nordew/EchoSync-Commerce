package grpcServer

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

func (s *grpcServer) GetMarket(ctx context.Context, req *grpcStore.GetStoreRequest) (*grpcStore.GetStoreResponse, error) {
	const op = "grpcServer.GetMarket"

	parsedUUID, err := uuid.Parse(req.GetStoreId())
	if err != nil {
		s.logger.Error(op, err.Error())
		return nil, ErrInvalidUUID
	}

	market, err := s.storeService.GetByID(context.Background(), parsedUUID)
	if err != nil {
		return nil, err
	}

	response := &grpcStore.GetStoreResponse{
		Store: &grpcStore.Store{
			StoreId: market.ID.String(),
			Name:    market.Name,
			OwnerId: market.OwnerUserID.String(),
		},
	}

	return response, nil
}

func (s *grpcServer) UpdateStore(ctx context.Context, req *grpcStore.UpdateStoreRequest) (*grpcStore.Empty, error) {
	const op = "grpcServer.UpdateStore"

	parsedUUID, err := uuid.Parse(req.GetStoreId())
	if err != nil {
		s.logger.Error(op, err.Error())
		return nil, ErrInvalidUUID
	}

	err = s.storeService.Update(context.Background(), req.GetName(), parsedUUID)
	if err != nil {
		s.logger.Error(op, "failed to update store", err.Error())
		return nil, err
	}

	return &grpcStore.Empty{}, nil
}

func (s *grpcServer) DeleteStore(ctx context.Context, req *grpcStore.DeleteStoreRequest) (*grpcStore.Empty, error) {
	const op = "grpcServer.DeleteStore"

	parsedUUID, err := uuid.Parse(req.GetStoreId())
	if err != nil {
		s.logger.Error(op, err.Error())
		return nil, ErrInvalidUUID
	}

	if err = s.storeService.Delete(context.Background(), parsedUUID); err != nil {
		s.logger.Error(op, "failed to delete store", err.Error())
		return nil, err
	}

	return &grpcStore.Empty{}, nil
}
