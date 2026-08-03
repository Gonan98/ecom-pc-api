package service

import (
	"context"
	"errors"

	repo "github.com/gonan98/ecom-pc-api/internal/repository"
	"github.com/gonan98/ecom-pc-api/internal/types"
	"github.com/gonan98/ecom-pc-api/internal/util"
	"github.com/jackc/pgx/v5"
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
	return s.categoryRepo.GetAll(ctx)
}

func (s *CategoryService) GetByID(ctx context.Context, ID int) (*types.Category, error) {
	category, err := s.categoryRepo.GetByID(ctx, ID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, util.ResourceNotFound("category", ID)
	}

	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s *CategoryService) Create(ctx context.Context, req *types.CreateCategoryRequest) error {
	category := &types.Category{
		Name:        req.Name,
		Description: req.Description,
	}

	return s.categoryRepo.Create(ctx, category)
}

func (s *CategoryService) Update(ctx context.Context, req *types.UpdateCategoryRequest, ID int) error {

	category, err := s.GetByID(ctx, ID)
	if err != nil {
		return err
	}

	category.Name = req.Name
	category.Description = req.Description

	return s.categoryRepo.Update(ctx, category)
}

func (s *CategoryService) Delete(ctx context.Context, ID int) error {
	_, err := s.GetByID(ctx, ID)
	if err != nil {
		return err
	}

	return s.categoryRepo.Delete(ctx, ID)
}
