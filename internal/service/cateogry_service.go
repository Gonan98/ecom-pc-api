package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gonan98/ecom-pc-api/internal/repository"
	"github.com/gonan98/ecom-pc-api/internal/types"
)

type CategoryService struct {
	categoryRepo *repository.CategoryRepository
}

func NewCategoryService(categoryRepo *repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		categoryRepo: categoryRepo,
	}
}

func (s *CategoryService) GetAll(ctx context.Context) ([]types.Category, error) {
	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s *CategoryService) GetByID(ctx context.Context, ID int) (*types.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	if category.ID == 0 {
		return nil, types.NewAPIError(http.StatusNotFound, fmt.Errorf("Category with ID=%d not found", ID))
	}

	return category, nil
}
