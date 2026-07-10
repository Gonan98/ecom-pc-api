package service

import (
	"context"
	"fmt"
	"net/http"

	repo "github.com/gonan98/ecom-pc-api/internal/repository"
	"github.com/gonan98/ecom-pc-api/internal/types"
)

type ProductService struct {
	productRepo *repo.ProductRepository
}

func NewProductService(productRepo *repo.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) GetAll(ctx context.Context) ([]types.Product, error) {
	return s.productRepo.GetAll(ctx)
}

func (s *ProductService) GetByID(ctx context.Context, ID int) (*types.Product, error) {
	product, err := s.productRepo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	if product.ID == 0 {
		return nil, types.NewAPIError(http.StatusNotFound, fmt.Errorf("Product with ID = %d not found", ID))
	}

	return product, nil
}
