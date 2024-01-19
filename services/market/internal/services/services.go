package services

import (
	"context"
	"github.com/google/uuid"
)

type StoreService interface {
	// Create creates a new store
	Create(ctx context.Context, name string, ownerID uuid.UUID) error
}
