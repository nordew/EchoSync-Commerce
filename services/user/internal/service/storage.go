package service

import (
	"context"
	"userService/internal/domain/entity"
)

// UserStorage is an interface for interacting with user data storage.
type UserStorage interface {
	// Create inserts a new user into the storage.
	Create(ctx context.Context, user *entity.User) error

	// Get retrieves a user from the storage based on the email.
	Get(ctx context.Context, email string) (*entity.User, error)
}
