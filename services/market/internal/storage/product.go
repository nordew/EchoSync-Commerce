package storage

import (
	"context"
	"github.com/jackc/pgx"
	"marketService/internal/domain/entity"
	"marketService/pkg/logger"
)

type productStorage struct {
	conn *pgx.Conn

	logger logger.Logger
}

func NewProductStorage(conn *pgx.Conn, logger logger.Logger) *productStorage {
	return &productStorage{
		conn:   conn,
		logger: logger,
	}
}

func (s *productStorage) Create(ctx context.Context, product *entity.Product) error {
	const op = "productStorage.Create"

	_, err := s.conn.ExecEx(ctx,
		"INSERT INTO products (product_id, store_id, product_name, price, quantity, created_at) VALUES($1, $2, $3, $4, $5, $6)",
		nil, product.ProductID, product.StoreID, product.ProductName, product.Price, product.Quantity, product.CreatedAt)
	if err != nil {
		s.logger.Error(op, "failed to create product", err.Error())
		return err
	}

	return nil
}
