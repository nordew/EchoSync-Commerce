package entity

import (
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"time"
)

var validate *validator.Validate

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
	Username string `validate:"required,gte=2"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,gte=6"`
}

type SignInInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required,gte=6"`
}

func (i SignUpInput) Validate() error {
	return validate.Struct(i)
}

func (i SignInInput) Validate() error {
	return validate.Struct(i)
}
