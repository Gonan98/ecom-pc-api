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

func (s *BrandService) Create(ctx context.Context, req *types.CreateBrandRequest) error {

	brand := &types.Brand{
		Name:    req.Name,
		Website: &req.Website,
	}

	return s.brandRepo.Create(ctx, brand)
}

func (s *BrandService) Update(ctx context.Context, req *types.UpdateBrandRequest, ID int) error {
	brand, err := s.GetByID(ctx, ID)
	if err != nil {
		return err
	}

	brand.Name = req.Name
	brand.Website = &req.Website

	return s.brandRepo.Update(ctx, brand)
}

func (s *BrandService) Delete(ctx context.Context, ID int) error {
	if _, err := s.GetByID(ctx, ID); err != nil {
		return err
	}

	return s.brandRepo.Delete(ctx, ID)
}
