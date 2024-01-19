package storage

import (
	"context"
	"github.com/google/uuid"
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

func (s *productStorage) GetByID(ctx context.Context, productID uuid.UUID) (*entity.Product, error) {
	const op = "productStorage.GetByID"

	var product entity.Product
	err := s.conn.QueryRowEx(ctx, "SELECT * FROM products WHERE product_id = $1", nil, productID).
		Scan(&product.ProductID, &product.StoreID, &product.ProductName, &product.Price, &product.Quantity, &product.CreatedAt)
	if err != nil {
		s.logger.Error(op, "failed to get product", err.Error())
		return nil, err
	}

	return &product, nil
}

func (s *productStorage) GetByStoreID(ctx context.Context, storeID uuid.UUID) ([]*entity.Product, error) {
	const op = "productStorage.GetByStoreID"

	var products []*entity.Product
	rows, err := s.conn.QueryEx(ctx, "SELECT * FROM products WHERE store_id = $1", nil, storeID)
	if err != nil {
		s.logger.Error(op, "failed to get products", err.Error())
		return nil, err
	}

	for rows.Next() {
		var product entity.Product
		err = rows.Scan(&product.ProductID, &product.StoreID, &product.ProductName, &product.Price, &product.Quantity, &product.CreatedAt)
		if err != nil {
			s.logger.Error(op, "failed to scan product", err.Error())
			return nil, err
		}

		products = append(products, &product)
	}

	return products, nil
}

func (s *productStorage) Update(ctx context.Context, product *entity.Product) error {
	const op = "productStorage.Update"

	_, err := s.conn.ExecEx(ctx,
		"UPDATE products SET product_name = $1, price = $2, quantity = $3 WHERE product_id = $4",
		nil, product.ProductName, product.Price, product.Quantity, product.ProductID)
	if err != nil {
		s.logger.Error(op, "failed to update product", err.Error())
		return err
	}

	return nil
}

func (s *productStorage) Delete(ctx context.Context, productID uuid.UUID) error {
	const op = "productStorage.Delete"

	_, err := s.conn.ExecEx(ctx, "DELETE FROM products WHERE product_id = $1", nil, productID)
	if err != nil {
		s.logger.Error(op, "failed to delete product", err.Error())
		return err
	}

	return nil
}
