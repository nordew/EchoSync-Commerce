package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"time"
	"userService/internal/domain/entity"
	"userService/pkg/auth"
	"userService/pkg/hasher"
	"userService/pkg/logger"
)

var (
	ErrPasswordsDoNotMatch = errors.New("passwords do not match")
)

type authService struct {
	userStorage UserStorage

	auth auth.Authenticator

	logger logger.Logger

	hasher hasher.Hasher
}

func NewAuthService(userStorage UserStorage, auth auth.Authenticator, logger logger.Logger, hasher hasher.Hasher) UserService {
	return &authService{
		userStorage: userStorage,
		auth:        auth,
		logger:      logger,
		hasher:      hasher,
	}
}

func (s *authService) SignUp(ctx context.Context, input *entity.SignUpInput) error {
	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		s.logger.Error("failed to hash password", err)
		return err
	}

	user := entity.User{
		UserID:       uuid.New(),
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	if err := s.userStorage.Create(ctx, &user); err != nil {
		s.logger.Error("failed to create user", err)
		return err
	}

	return nil
}

func (s *authService) SignIn(ctx context.Context, input *entity.SignInInput) (string, string, error) {
	hashedPassword, err := s.hasher.Hash(input.Password)

	user, err := s.userStorage.Get(ctx, input.Email)

	accessToken, refreshToken, err := s.auth.GenerateTokens(&auth.GenerateTokenClaimsOptions{
		UserId: user.UserID.String(),
	})
	if err != nil {
		s.logger.Error("failed to generate tokens", err)
		return "", "", err
	}

	if hashedPassword != user.PasswordHash {
		s.logger.Error("passwords do not match", nil)
		return "", "", ErrPasswordsDoNotMatch
	}

	if err := s.userStorage.CreateRefreshToken(ctx, user.UserID, refreshToken); err != nil {
		s.logger.Error("failed to create refresh token", err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) RefreshTokens(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := s.auth.ParseToken(refreshToken)
	if err != nil {
		s.logger.Error("failed to parse refresh token", err)
		return "", "", err
	}

	accessToken, newRefreshToken, err := s.auth.GenerateTokens(&auth.GenerateTokenClaimsOptions{
		UserId: claims.Sub,
	})
	if err != nil {
		s.logger.Error("failed to generate tokens", err)
		return "", "", err
	}

	parsedUUID, err := uuid.Parse(claims.Sub)
	if err != nil {
		s.logger.Error("failed to parse uuid", err)
		return "", "", err
	}

	if err := s.userStorage.CreateRefreshToken(ctx, parsedUUID, newRefreshToken); err != nil {
		s.logger.Error("failed to create refresh token", err)
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}
