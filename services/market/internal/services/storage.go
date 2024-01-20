package services

import (
	"context"
	"github.com/google/uuid"
	"marketService/internal/domain/entity"
)

type StoreStorage interface {
	// Create creates a new store
	Create(ctx context.Context, store *entity.Store) error

	// GetByOwnerID returns all stores by owner id
	GetByID(ctx context.Context, storeID uuid.UUID) (*entity.Store, error)

	// Update updates store
	Update(ctx context.Context, storeName string, storeID uuid.UUID) error

	// Delete deletes store
	Delete(ctx context.Context, storeID uuid.UUID) error
}

type ProductStorage interface {
	// Create creates a new product
	Create(ctx context.Context, product *entity.Product) error

	// GetByStoreID returns all products by store id
	GetByStoreID(ctx context.Context, storeID uuid.UUID) ([]*entity.Product, error)

	// GetByID returns product by id
	GetByID(ctx context.Context, productID uuid.UUID) (*entity.Product, error)

	// Update updates product
	Update(ctx context.Context, product *entity.Product) error

	// Delete deletes product
	Delete(ctx context.Context, productID uuid.UUID) error
}
