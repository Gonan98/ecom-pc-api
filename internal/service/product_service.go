package service

import (
	"context"
	"errors"

	repo "github.com/gonan98/ecom-pc-api/internal/repository"
	"github.com/gonan98/ecom-pc-api/internal/types"
	"github.com/gonan98/ecom-pc-api/internal/util"
	"github.com/jackc/pgx/v5"
)

type ProductService struct {
	productRepo     *repo.ProductRepository
	brandService    *BrandService
	categoryService *CategoryService
}

func NewProductService(productRepo *repo.ProductRepository, brandService *BrandService, categoryService *CategoryService) *ProductService {
	return &ProductService{
		productRepo:     productRepo,
		brandService:    brandService,
		categoryService: categoryService,
	}
}

func (s *ProductService) GetAll(ctx context.Context) ([]types.Product, error) {
	return s.productRepo.GetAll(ctx)
}

func (s *ProductService) GetByID(ctx context.Context, ID int) (*types.Product, error) {
	product, err := s.productRepo.GetByID(ctx, ID)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, util.ResourceNotFound("product", ID)
	}

	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *ProductService) Create(ctx context.Context, req *types.CreateProductRequest) error {
	_, err := s.brandService.GetByID(ctx, req.BrandID)
	if err != nil {
		return err
	}

	_, err = s.categoryService.GetByID(ctx, req.CategoryID)
	if err != nil {
		return err
	}

	product := &types.Product{
		CategoryID:  req.CategoryID,
		BrandID:     req.BrandID,
		Name:        req.Name,
		Description: req.Description,
		ImageUrl:    req.ImageUrl,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	return s.productRepo.Create(ctx, product)
}

func (s *ProductService) Update(ctx context.Context, req *types.UpdateProductRequest, ID int) error {
	_, err := s.brandService.GetByID(ctx, req.BrandID)
	if err != nil {
		return err
	}

	_, err = s.categoryService.GetByID(ctx, req.CategoryID)
	if err != nil {
		return err
	}

	product, err := s.GetByID(ctx, ID)
	if err != nil {
		return err
	}

	product.CategoryID = req.CategoryID
	product.BrandID = req.BrandID
	product.Name = req.Name
	product.Description = req.Description
	product.ImageUrl = req.ImageUrl
	product.Price = req.Price
	product.Stock = req.Stock

	return s.productRepo.Update(ctx, product)
}

func (s *ProductService) Delete(ctx context.Context, ID int) error {
	if _, err := s.GetByID(ctx, ID); err != nil {
		return err
	}

	return s.productRepo.Delete(ctx, ID)
}
