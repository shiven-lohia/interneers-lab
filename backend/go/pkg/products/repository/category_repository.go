package repository

import (
	"context"

	"github.com/shiven-lohia/interneers-lab/pkg/products/entity"
)

type ProductCategoryRepository interface {
	Create(ctx context.Context, category entity.ProductCategory) (entity.ProductCategory, error)
	GetByID(ctx context.Context, id string) (entity.ProductCategory, error)
	GetAll(ctx context.Context) ([]entity.ProductCategory, error)
	Update(ctx context.Context, id string, category entity.ProductCategory) (entity.ProductCategory, error)
	Delete(ctx context.Context, id string) error
}