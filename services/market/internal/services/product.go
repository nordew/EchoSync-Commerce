package services

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"marketService/internal/domain/entity"
	"marketService/pkg/logger"
	"time"
)

var (
	ErrInvalidProductPrice    = errors.New("invalid product price")
	ErrInvalidProductQuantity = errors.New("invalid product quantity")
)

type productService struct {
	productStorage ProductStorage

	logger logger.Logger
}

func NewProductService(productStorage ProductStorage, logger logger.Logger) ProductService {
	return &productService{
		productStorage: productStorage,
		logger:         logger,
	}
}

func (s *productService) Create(ctx context.Context, storeID uuid.UUID, productName string, price, quantity int) error {
	const op = "productService.Create"

	createdProduct := &entity.Product{
		ProductID:   uuid.New(),
		StoreID:     storeID,
		ProductName: productName,
		Price:       price,
		Quantity:    quantity,
		CreatedAt:   time.Now(),
	}

	if err := validateProduct(createdProduct); err != nil {
		s.logger.Error(op, err.Error())
		return err
	}

	err := s.productStorage.Create(ctx, createdProduct)
	if err != nil {
		s.logger.Error(op, err.Error())
		return err
	}

	return nil
}

func validateProduct(product *entity.Product) error {
	if product.Price < 0 {
		return ErrInvalidProductPrice
	} else if product.Quantity < 0 {
		return ErrInvalidProductQuantity
	}

	return nil
}
