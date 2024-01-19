package services

import (
	"context"
	"marketService/internal/domain/entity"
)

type StoreStorage interface {
	// Create creates a new store
	Create(ctx context.Context, store *entity.Store) error
}

type ProductStorage interface {
	// Create creates a new product
	Create(ctx context.Context, product *entity.Product) error
}
