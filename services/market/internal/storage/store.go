package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx"
	"marketService/internal/domain/entity"
	"marketService/pkg/logger"
)

var (
	ErrMaximalStoresCountReached = errors.New("maximal stores count reached")
)

type storeStorage struct {
	conn *pgx.Conn

	logging logger.Logger
}

func NewStoreStorage(conn *pgx.Conn, logging logger.Logger) *storeStorage {
	return &storeStorage{
		conn:    conn,
		logging: logging,
	}
}

func (s *storeStorage) Create(ctx context.Context, store *entity.Store) error {
	const op = "storeStorage.Create"

	var storesCount int
	err := s.conn.QueryRowEx(ctx, "SELECT stores_active FROM users WHERE user_id = $1", nil, store.OwnerUserID).
		Scan(&storesCount)
	if err != nil {
		s.logging.Error(op, err.Error())
		return err
	}

	if storesCount >= 3 {
		return ErrMaximalStoresCountReached
	}

	_, err = s.conn.ExecEx(ctx,
		"INSERT INTO stores (store_id, store_name, owner_user_id, products_count, is_active, created_at) VALUES($1, $2, $3, $4, $5, $6)",
		nil, store.ID, store.Name, store.OwnerUserID, store.ProductsCount, store.IsActive, store.CreatedAt)
	if err != nil {
		s.logging.Error(op, err.Error())
		return err
	}

	_, err = s.conn.ExecEx(ctx, "UPDATE users SET stores_active = stores_active + 1 WHERE user_id = $1", nil, store.OwnerUserID)
	if err != nil {
		s.logging.Error(op, err.Error())
		return err
	}

	return nil
}
