package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"marketService/internal/domain/entity"
	"marketService/pkg/logger"
	"time"
	"unicode/utf8"
)

var (
	ErrInvalidStoreName = errors.New("invalid store name")
)

type storeService struct {
	storeStorage StoreStorage

	logger logger.Logger
}

func NewStoreService(storeStorage StoreStorage, logger logger.Logger) StoreService {
	return &storeService{
		storeStorage: storeStorage,
		logger:       logger,
	}
}

func (s *storeService) Create(ctx context.Context, name string, ownerID uuid.UUID) error {
	const op = "storeService.Create"

	s.logger.Info(op, "creating store\n")
	createdStore := &entity.Store{
		ID:          uuid.New(),
		Name:        name,
		OwnerUserID: ownerID,
		CreatedAt:   time.Now(),
	}

	s.logger.Info(op, "validating store\n")
	if err := validateStore(createdStore); err != nil {
		s.logger.Error(op, err.Error())
		return err
	}

	s.logger.Info(op, "creating store\n")
	err := s.storeStorage.Create(ctx, createdStore)
	if err != nil {
		s.logger.Error(op, err.Error())
		return err
	}

	return nil
}

func (s *storeService) GetByID(ctx context.Context, storeID uuid.UUID) (*entity.Store, error) {
	const op = "storeService.GetByID"

	store, err := s.storeStorage.GetByID(ctx, storeID)
	if err != nil {
		s.logger.Error(op, err.Error())
		return nil, err
	}

	return store, nil
}

func (s *storeService) Update(ctx context.Context, storeName string, storeID uuid.UUID) error {
	const op = "storeService.Update"

	if err := validateStore(&entity.Store{Name: storeName}); err != nil {
		s.logger.Error(op, err.Error())
		return err
	}

	if err := s.storeStorage.Update(ctx, storeName, storeID); err != nil {
		s.logger.Error(op, err.Error())
		return err
	}

	return nil
}

func (s *storeService) Delete(ctx context.Context, storeID uuid.UUID) error {
	const op = "storeService.Delete"

	if err := s.storeStorage.Delete(ctx, storeID); err != nil {
		s.logger.Error(op, err.Error())
		return err
	}

	return nil
}

func validateStore(store *entity.Store) error {
	if utf8.RuneCountInString(store.Name) > 40 || utf8.RuneCountInString(store.Name) < 3 {
		return ErrInvalidStoreName
	}

	return nil
}
