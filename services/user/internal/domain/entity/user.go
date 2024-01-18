package entity

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	UserID       uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	RefreshToken string
	StoresActive int
	CreatedAt    time.Time
}

type SignUpInput struct {
	Username string
	Email    string
	Password string
}

type SignInInput struct {
	Email    string
	Password string
}
