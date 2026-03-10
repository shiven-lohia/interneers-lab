package controller

import (
	"errors"
	"github.com/shiven-lohia/interneers-lab/pkg/products/entity"
	"github.com/shiven-lohia/interneers-lab/pkg/products/repository"
)

type ProductController struct {
	repo repository.ProductRepository
}

func NewProductController(repo repository.ProductRepository) *ProductController {
	return &ProductController{
		repo: repo,
	}
}

// var ProductStore = map[string]entity.Product{}

func (s *ProductController) CreateProduct(p entity.Product) (entity.Product, error) {

	if p.Name == "" {
		return entity.Product{}, errors.New("Name is required")
	}

	if p.Price <= 0 {
		return entity.Product{}, errors.New("Price must be greater than 0")
	}

	if p.Quantity < 0 {
		return entity.Product{}, errors.New("Quantity cannot be negative")
	}

	return s.repo.Create(p)
}

func (s *ProductController) GetAllProducts() ([]entity.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductController) GetProductById(id string) (entity.Product, error) {
	return s.repo.GetByID(id)
}

func (s *ProductController) UpdateProduct(id string, p entity.Product) (entity.Product, error) {

	if p.Name == "" {
		return entity.Product{}, errors.New("name is required")
	}

	if p.Price <= 0 {
		return entity.Product{}, errors.New("price must be greater than 0")
	}

	if p.Quantity < 0 {
		return entity.Product{}, errors.New("quantity cannot be negative")
	}

	return s.repo.Update(id, p)
}

func (s *ProductController) DeleteProduct(id string) error {
	return s.repo.Delete(id)
}