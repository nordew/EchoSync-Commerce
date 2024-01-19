package entity

import (
	"github.com/google/uuid"
	"time"
)

type Product struct {
	ProductID   uuid.UUID
	StoreID     uuid.UUID
	ProductName string
	Price       int
	Quantity    int
	CreatedAt   time.Time
}
