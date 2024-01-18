package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"marketService/internal/domain/entity"
	"marketService/pkg/logger"
	"unicode/utf8"
)

var (
	ErrInvalidStoreName = errors.New("invalid store name")
)

type storeService struct {
	storeStorage StoreStorage

	logger logger.Logger
}

func NewStoreService(storeStorage StoreStorage, logger logger.Logger) *storeService {
	return &storeService{
		storeStorage: storeStorage,
		logger:       logger,
	}
}

func (s *storeService) Create(ctx context.Context, store *entity.Store) error {
	const op = "storeService.Create"

	createdStore := &entity.Store{
		ID:          uuid.New(),
		Name:        store.Name,
		OwnerUserID: store.OwnerUserID,
		CreatedAt:   store.CreatedAt,
	}

	if err := validateStore(createdStore); err != nil {
		s.logger.Error(op, err.Error())
		return err
	}

	err := s.storeStorage.Create(ctx, createdStore)
	if err != nil {
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
