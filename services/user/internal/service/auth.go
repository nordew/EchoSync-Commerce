package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"sync"
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
	s.logger.Info("hashing password")
	hashedPassword, err := s.hasher.Hash(input.Password)
	if err != nil {
		s.logger.Error("failed to hash password", err)
		return err
	}

	s.logger.Info("creating user")
	user := entity.User{
		UserID:       uuid.New(),
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: hashedPassword,
		CreatedAt:    time.Now(),
	}

	s.logger.Info("saving user")
	if err := s.userStorage.Create(ctx, &user); err != nil {
		s.logger.Error("failed to create user", err)
		return err
	}

	return nil
}

func (s *authService) SignIn(ctx context.Context, input *entity.SignInInput) (string, string, error) {
	var (
		hashedPassword string
		user           *entity.User
		hashErr        error
		getUserErr     error
	)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		hashedPassword, hashErr = s.hasher.Hash(input.Password)
	}()

	go func() {
		user, getUserErr = s.userStorage.Get(ctx, input.Email)
	}()

	wg.Wait()

	if hashErr != nil {
		s.logger.Error("failed to hash password", hashErr)
		return "", "", hashErr
	}

	if getUserErr != nil {
		s.logger.Error("failed to get user", getUserErr)
		return "", "", getUserErr
	}

	if hashedPassword != user.PasswordHash {
		s.logger.Error("passwords do not match", nil)
		return "", "", ErrPasswordsDoNotMatch
	}

	accessToken, refreshToken, err := s.auth.GenerateTokens(&auth.GenerateTokenClaimsOptions{
		UserId: user.UserID.String(),
	})
	if err != nil {
		s.logger.Error("failed to generate tokens", err)
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
