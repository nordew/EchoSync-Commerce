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

func (s *productService) GetByStoreID(ctx context.Context, storeID uuid.UUID) ([]*entity.Product, error) {
	const op = "productService.GetByStoreID"

	products, err := s.productStorage.GetByStoreID(ctx, storeID)
	if err != nil {
		s.logger.Error(op, err.Error())
		return nil, err
	}

	return products, nil
}

func (s *productService) GetByID(ctx context.Context, productID uuid.UUID) (*entity.Product, error) {
	const op = "productService.GetByID"

	product, err := s.productStorage.GetByID(ctx, productID)
	if err != nil {
		s.logger.Error(op, err.Error())
		return nil, err
	}

	return product, nil
}

func (s *productService) Update(ctx context.Context, product *entity.Product) error {
	const op = "productService.Update"

	if err := validateProduct(product); err != nil {
		s.logger.Error(op, err.Error())
		return err
	}

	err := s.productStorage.Update(ctx, product)
	if err != nil {
		s.logger.Error(op, err.Error())
		return err
	}

	return nil
}

func (s *productService) Delete(ctx context.Context, productID uuid.UUID) error {
	const op = "productService.Delete"

	if err := s.productStorage.Delete(ctx, productID); err != nil {
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
