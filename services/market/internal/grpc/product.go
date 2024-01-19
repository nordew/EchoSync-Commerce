package grpcStore

import (
	"context"
	"github.com/google/uuid"
	grpcStore "github.com/nordew/EchoSync-protos/gen/go/store"
)

func (s *grpcServer) CreateProduct(ctx context.Context, req *grpcStore.CreateProductRequest) (*grpcStore.Empty, error) {
	const op = "grpcServer.CreateProduct"

	parsedStoreUUID, err := uuid.Parse(req.GetStoreId())
	if err != nil {
		s.logger.Error(op, err.Error())
		return nil, ErrInvalidUUID
	}

	err = s.productService.Create(context.Background(), parsedStoreUUID, req.GetProductName(), int(req.GetPrice()), int(req.GetQuantity()))
	if err != nil {
		s.logger.Error(op, "failed to create product", err.Error())
		return nil, err
	}

	return &grpcStore.Empty{}, nil
}

func (s *grpcServer) GetProductsByStoreID(ctx context.Context, req *grpcStore.GetProductsByStoreIDRequest) (*grpcStore.GetProductsByStoreIDResponse, error) {
	const op = "grpcServer.GetProductsByStoreID"

	parsedUUID, err := uuid.Parse(req.GetStoreId())
	if err != nil {
		s.logger.Error(op, err.Error())
		return nil, ErrInvalidUUID
	}

	products, err := s.productService.GetByStoreID(context.Background(), parsedUUID)
	if err != nil {
		s.logger.Error(op, "failed to get products", err.Error())
		return nil, err
	}

	var grpcProducts []*grpcStore.Product
	for _, product := range products {
		grpcProducts = append(grpcProducts, &grpcStore.Product{
			ProductId:   product.ProductID.String(),
			StoreId:     product.StoreID.String(),
			ProductName: product.ProductName,
			Price:       int32(product.Price),
			Quantity:    int32(product.Quantity),
		})
	}

	return &grpcStore.GetProductsByStoreIDResponse{
		Products: grpcProducts,
	}, nil
}

func (s *grpcServer) GetProductByID(ctx context.Context, req *grpcStore.GetProductByIDRequest) (*grpcStore.GetProductByIDResponse, error) {
	const op = "grpcServer.GetProductByID"

	parsedUUID, err := uuid.Parse(req.GetProductId())
	if err != nil {
		s.logger.Error(op, err.Error())
		return nil, ErrInvalidUUID
	}

	product, err := s.productService.GetByID(context.Background(), parsedUUID)
	if err != nil {
		s.logger.Error(op, "failed to get product", err.Error())
		return nil, err
	}

	return &grpcStore.GetProductByIDResponse{
		Product: &grpcStore.Product{
			ProductId:   product.ProductID.String(),
			StoreId:     product.StoreID.String(),
			ProductName: product.ProductName,
			Price:       int32(product.Price),
			Quantity:    int32(product.Quantity),
		},
	}, nil
}

func (s *grpcServer) UpdateProduct(ctx context.Context, req *grpcStore.UpdateProductRequest) (*grpcStore.Empty, error) {
	const op = "grpcServer.UpdateProduct"

	parsedUUID, err := uuid.Parse(req.Product.GetProductId())
	if err != nil {
		s.logger.Error(op, err.Error())
		return nil, ErrInvalidUUID
	}

	product, err := s.productService.GetByID(context.Background(), parsedUUID)
	if err != nil {
		s.logger.Error(op, err.Error())
		return nil, err
	}

	product.ProductName = req.Product.GetProductName()
	product.Price = int(req.Product.GetPrice())
	product.Quantity = int(req.Product.GetQuantity())

	err = s.productService.Update(context.Background(), product)
	if err != nil {
		s.logger.Error(op, err.Error())
		return nil, err
	}

	return &grpcStore.Empty{}, nil
}

func (s *grpcServer) DeleteProduct(ctx context.Context, req *grpcStore.DeleteProductRequest) (*grpcStore.Empty, error) {
	const op = "grpcServer.DeleteProduct"

	parsedUUID, err := uuid.Parse(req.GetProductId())
	if err != nil {
		s.logger.Error(op, err.Error())
		return nil, ErrInvalidUUID
	}

	err = s.productService.Delete(context.Background(), parsedUUID)
	if err != nil {
		s.logger.Error(op, err.Error())
		return nil, err
	}

	return &grpcStore.Empty{}, nil
}
