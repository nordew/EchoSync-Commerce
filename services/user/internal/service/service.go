package service

import (
	"context"
	"userService/internal/domain/entity"
)

type UserService interface {
	SignUp(ctx context.Context, input *entity.SignUpInput) error

	// SignIn returns a new access token and refresh token.
	SignIn(ctx context.Context, input *entity.SignInInput) (string, string, error)
}
