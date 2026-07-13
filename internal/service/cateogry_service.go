package service

import (
	"context"
	"fmt"
	"net/http"

	repo "github.com/gonan98/ecom-pc-api/internal/repository"
	"github.com/gonan98/ecom-pc-api/internal/types"
)

type CategoryService struct {
	categoryRepo *repo.CategoryRepository
}

func NewCategoryService(categoryRepo *repo.CategoryRepository) *CategoryService {
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

func (s *CategoryService) Create(ctx context.Context, req *types.CreateCategoryRequest) error {
	category := &types.Category{
		Name:        req.Name,
		Description: &req.Description,
	}

	return s.categoryRepo.Create(ctx, category)
}

func (s *CategoryService) Update(ctx context.Context, req *types.UpdateCategoryRequest, ID int) error {

	category, err := s.GetByID(ctx, ID)
	if err != nil {
		return err
	}

	category.Name = req.Name
	category.Description = &req.Description

	return s.categoryRepo.Update(ctx, category)
}

func (s *CategoryService) Delete(ctx context.Context, ID int) error {
	_, err := s.GetByID(ctx, ID)
	if err != nil {
		return err
	}

	return s.categoryRepo.Delete(ctx, ID)
}
