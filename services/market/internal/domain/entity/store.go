package entity

import (
	"github.com/google/uuid"
	"time"
)

type Store struct {
	ID            uuid.UUID
	Name          string
	OwnerUserID   uuid.UUID
	ProductsCount int
	IsActive      bool
	CreatedAt     time.Time
}
