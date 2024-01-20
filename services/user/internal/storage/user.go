package storage

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx"
	"userService/internal/domain/entity"
	"userService/pkg/logger"
)

type userStorage struct {
	conn *pgx.Conn

	logger logger.Logger
}

func NewUserStorage(conn *pgx.Conn, logger logger.Logger) *userStorage {
	return &userStorage{
		conn:   conn,
		logger: logger,
	}
}

func (s *userStorage) Create(ctx context.Context, user *entity.User) error {
	_, err := s.conn.ExecEx(ctx,
		"INSERT INTO users (user_id, username, email, password_hash, refresh_token, stores_active, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
		nil, user.UserID, user.Username, user.Email, user.PasswordHash, user.RefreshToken, user.StoresActive, user.CreatedAt)
	if err != nil {
		s.logger.Error("failed to create user", err)
		return err
	}

	return nil
}

func (s *userStorage) Get(ctx context.Context, email, refreshToken string) (*entity.User, error) {
	var user entity.User

	_, err := s.conn.ExecEx(ctx, "UPDATE users SET refresh_token=$1 WHERE email=$2", nil, refreshToken, email)
	if err != nil {
		s.logger.Error("failed to update refresh token", err)
		return nil, err
	}

	err = s.conn.QueryRowEx(ctx, "SELECT * FROM users WHERE email=$1", nil, email).
		Scan(&user.UserID, &user.Username, &user.Email, &user.PasswordHash, &user.RefreshToken, &user.StoresActive, &user.CreatedAt)
	if err != nil {
		s.logger.Error("failed to get user", err)
		return nil, err
	}

	return &user, nil
}

func (s *userStorage) CreateRefreshToken(ctx context.Context, userID uuid.UUID, token string) error {
	_, err := s.conn.ExecEx(ctx, "UPDATE users SET refresh_token=$1 WHERE user_id=$2", nil, token, userID)
	if err != nil {
		s.logger.Error("failed to create refresh token", err)
		return err
	}

	return nil
}
