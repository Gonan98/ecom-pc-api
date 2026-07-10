package service

import (
	"context"
	"fmt"
	"net/http"

	repo "github.com/gonan98/ecom-pc-api/internal/repository"
	"github.com/gonan98/ecom-pc-api/internal/types"
)

type BrandService struct {
	brandRepo *repo.BrandRepository
}

func NewBrandService(brandRepo *repo.BrandRepository) *BrandService {
	return &BrandService{
		brandRepo: brandRepo,
	}
}

func (s *BrandService) GetAll(ctx context.Context) ([]types.Brand, error) {
	brands, err := s.brandRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return brands, nil
}

func (s *BrandService) GetByID(ctx context.Context, ID int) (*types.Brand, error) {
	brand, err := s.brandRepo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	if brand.ID == 0 {
		return nil, types.NewAPIError(http.StatusNotFound, fmt.Errorf("Brand with ID=%d not found", ID))
	}

	return brand, nil
}
