package service

import (
	"context"
	"fmt"
	"net/http"

	repo "github.com/gonan98/ecom-pc-api/internal/repository"
	"github.com/gonan98/ecom-pc-api/internal/types"
	"github.com/gonan98/ecom-pc-api/internal/util"
)

type ProductService struct {
	productRepo  *repo.ProductRepository
	brandRepo    *repo.BrandRepository
	categoryRepo *repo.CategoryRepository
}

func NewProductService(productRepo *repo.ProductRepository, brandRepo *repo.BrandRepository, categoryRepo *repo.CategoryRepository) *ProductService {
	return &ProductService{
		productRepo:  productRepo,
		brandRepo:    brandRepo,
		categoryRepo: categoryRepo,
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

func (s *ProductService) Create(ctx context.Context, req *types.CreateProductRequest) error {
	brand, err := s.brandRepo.GetByID(ctx, req.BrandID)
	if err != nil {
		return err
	}

	if brand.ID == 0 {
		return util.ResourceNotFound("brand", req.BrandID)
	}

	category, err := s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return err
	}

	if category.ID == 0 {
		return util.ResourceNotFound("category", req.CategoryID)
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
	brand, err := s.brandRepo.GetByID(ctx, req.BrandID)
	if err != nil {
		return err
	}

	if brand.ID == 0 {
		return util.ResourceNotFound("brand", req.BrandID)
	}

	category, err := s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return err
	}

	if category.ID == 0 {
		return util.ResourceNotFound("category", req.CategoryID)
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
