package repository

import "github.com/shiven-lohia/interneers-lab/pkg/products/entity"

type ProductRepository interface {
	Create(product entity.Product) (entity.Product, error)
	GetAll() ([]entity.Product, error)
	GetByID(id string) (entity.Product, error)
	Update(id string, product entity.Product) (entity.Product, error)
	Delete(id string) error
}