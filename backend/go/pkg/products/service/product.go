package service

import (
	"context"
	"errors"
	"sync"

	"github.com/shiven-lohia/interneers-lab/pkg/products/entity"
	"github.com/shiven-lohia/interneers-lab/pkg/products/repository"
)

type ProductService struct {
	repo repository.ProductRepository
	categoryRepo repository.ProductCategoryRepository
}

func NewProductService(
	repo repository.ProductRepository,
	categoryRepo repository.ProductCategoryRepository,
) *ProductService {
	return &ProductService{
		repo: repo,
		categoryRepo: categoryRepo,
	}
}

// var ProductStore = map[string]entity.Product{}

func (s *ProductService) BulkCreateProducts(ctx context.Context, products []entity.Product) ([]entity.Product, error) {

	var wg sync.WaitGroup
	var mu sync.Mutex

	var createdProducts []entity.Product

	for _, p := range products {
		wg.Add(1)

		go func(prod entity.Product) {
			defer wg.Done()
			
			createdProduct, err := s.CreateProduct(ctx, prod)
			if(err!=nil) {
				return
			}

			mu.Lock()
			createdProducts = append(createdProducts, createdProduct)
			mu.Unlock()
		} (p)
	}

	wg.Wait()

	return createdProducts, nil
}

func (s *ProductService) CreateProduct(ctx context.Context, p entity.Product) (entity.Product, error) {

	if p.Name == "" {
		return entity.Product{}, errors.New("Name is required")
	}

	if p.Price <= 0 {
		return entity.Product{}, errors.New("Price must be greater than 0")
	}

	if p.Quantity < 0 {
		return entity.Product{}, errors.New("Quantity cannot be negative")
	}

	if p.Brand == "" {
		return entity.Product{}, errors.New("Brand is required")
	}

	if p.CategoryID != "" {
		_, err := s.categoryRepo.GetByID(ctx, p.CategoryID)
		if err != nil {
			return entity.Product{}, errors.New("Invalid category ID")
		}
	}

	return s.repo.Create(ctx, p)
}

func (s *ProductService) GetAllProducts(ctx context.Context, categoryID string) ([]entity.Product, error) {
	return s.repo.GetAll(ctx, categoryID)
}

func (s *ProductService) GetProductById(ctx context.Context, id string) (entity.Product, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *ProductService) UpdateProduct(ctx context.Context, id string, p entity.Product) (entity.Product, error) {

	if p.Name == "" {
		return entity.Product{}, errors.New("Name is required")
	}

	if p.Price <= 0 {
		return entity.Product{}, errors.New("Price must be greater than 0")
	}

	if p.Quantity < 0 {
		return entity.Product{}, errors.New("Quantity cannot be negative")
	}

	if p.Brand == "" {
		return entity.Product{}, errors.New("Brand is required")
	}

	if p.CategoryID != "" {
		_, err := s.categoryRepo.GetByID(ctx, p.CategoryID)
		if err != nil {
			return entity.Product{}, errors.New("Invalid category ID")
		}
	}

	return s.repo.Update(ctx, id, p)
}

func (s *ProductService) DeleteProduct(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
