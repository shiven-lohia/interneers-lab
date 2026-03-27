package service

import (
	"context"
	"errors"

	"github.com/shiven-lohia/interneers-lab/pkg/products/entity"
	"github.com/shiven-lohia/interneers-lab/pkg/products/repository"
)

type ProductCategoryService struct {
	repo repository.ProductCategoryRepository
}

func NewProductCategoryService(repo repository.ProductCategoryRepository) *ProductCategoryService {
	return &ProductCategoryService{
		repo: repo,
	}
}

func (s *ProductCategoryService) CreateCategory(ctx context.Context, c entity.ProductCategory) (entity.ProductCategory, error) {
	if c.Title == "" {
		return entity.ProductCategory{}, errors.New("Title is required")
	}
	
	return s.repo.Create(ctx, c)
}

func (s *ProductCategoryService) GetCategoryByID(ctx context.Context, id string) (entity.ProductCategory, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductCategoryService) GetAllCategories(ctx context.Context) ([]entity.ProductCategory, error) {
	return s.repo.GetAll(ctx)
}

func (s *ProductCategoryService) UpdateCategory(ctx context.Context, id string, c entity.ProductCategory) (entity.ProductCategory, error) {
	if c.Title == "" {
		return entity.ProductCategory{}, errors.New("Title is required")
	}
	return s.repo.Update(ctx, id, c)
}

func (s *ProductCategoryService) DeleteCategory(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}