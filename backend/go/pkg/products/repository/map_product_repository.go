package repository

import (
	"errors"

	"github.com/shiven-lohia/interneers-lab/pkg/products/entity"
)

type MapProductRepository struct {
	store map[string]entity.Product
}

func NewMapProductRepository() *MapProductRepository {
	return &MapProductRepository{
		store: make(map[string]entity.Product),
	}
}

func (r *MapProductRepository) Create(product entity.Product) (entity.Product, error) {
	r.store[product.ID] = product
	return product, nil
}

func (r *MapProductRepository) GetAll() ([]entity.Product, error) {
	products := []entity.Product{}

	for _, p := range r.store {
		products = append(products, p)
	}

	return products, nil
}

func (r *MapProductRepository) GetByID(id string) (entity.Product, error) {
	product, exists := r.store[id]

	if !exists {
		return entity.Product{}, errors.New("product not found")
	}

	return product, nil
}

func (r *MapProductRepository) Update(id string, product entity.Product) (entity.Product, error) {
	_, exists := r.store[id]

	if !exists {
		return entity.Product{}, errors.New("product not found")
	}

	product.ID = id
	r.store[id] = product

	return product, nil
}

func (r *MapProductRepository) Delete(id string) error {
	_, exists := r.store[id]

	if !exists {
		return errors.New("product not found")
	}

	delete(r.store, id)
	return nil
}