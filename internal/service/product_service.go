package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gonan98/ecom-pc-api/internal/model"
	"github.com/gonan98/ecom-pc-api/internal/repository"
)

type ProductService struct {
	productRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
	}
}

func (s *ProductService) GetAll(ctx context.Context) ([]model.Product, error) {
	return s.productRepo.GetAll(ctx)
}

func (s *ProductService) GetByID(ctx context.Context, ID int) (*model.Product, error) {
	product, err := s.productRepo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	if product.ID == 0 {
		return nil, model.NewAPIError(http.StatusNotFound, fmt.Errorf("Product with ID = %d not found", ID))
	}

	return product, nil
}
