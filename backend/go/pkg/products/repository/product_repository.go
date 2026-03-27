package repository

import (
	"context"

	"github.com/shiven-lohia/interneers-lab/pkg/products/entity"
)

type ProductRepository interface {
	Create(ctx context.Context, product entity.Product) (entity.Product, error)
	GetAll(ctx context.Context, categoryID string) ([]entity.Product, error)
	GetByID(ctx context.Context, id string) (entity.Product, error)
	Update(ctx context.Context, id string, product entity.Product) (entity.Product, error)
	Delete(ctx context.Context, id string) error
}