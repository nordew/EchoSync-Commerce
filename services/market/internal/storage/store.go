package storage

import (
	"context"
	"github.com/jackc/pgx"
	"marketService/internal/domain/entity"
	"marketService/pkg/logger"
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
	_, err := s.conn.ExecEx(ctx,
		"INSERT INTO stores (store_id, store_name, owner_user_id, products_count, is_active, created_at) VALUES($1, $2, $3, $4, $5, $6)",
		nil, store.ID, store.Name, store.OwnerUserID, store.ProductsCount, store.IsActive, store.CreatedAt)
	if err != nil {
		s.logging.Error("failed to create store", err.Error())
		return err
	}

	return nil
}
