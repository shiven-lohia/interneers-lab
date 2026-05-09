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

// var ProductStore = map[string]entity.Product{} (deprecated, using repository instead)

const bulkWorkers = 10

type BulkCreateError struct {
	Index  int    `json:"index"`
	Reason string `json:"reason"`
}

type BulkCreateResult struct {
	Created []entity.Product  `json:"created"`
	Errors  []BulkCreateError `json:"errors"`
}

func (s *ProductService) BulkCreateProducts(ctx context.Context, products []entity.Product) (BulkCreateResult, error) {
	sem := make(chan struct{}, bulkWorkers)
	var wg sync.WaitGroup
	var mu sync.Mutex

	result := BulkCreateResult{
		Created: []entity.Product{},
		Errors:  []BulkCreateError{},
	}

	for i, p := range products {
		wg.Add(1)
		sem <- struct{}{}

		go func(idx int, prod entity.Product) {
			defer wg.Done()
			defer func() { <-sem }()

			created, err := s.CreateProduct(ctx, prod)

			mu.Lock()
			defer mu.Unlock()

			if err != nil {
				result.Errors = append(result.Errors, BulkCreateError{Index: idx, Reason: err.Error()})
			} else {
				result.Created = append(result.Created, created)
			}
		}(i, p)
	}

	wg.Wait()
	return result, nil
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
	products, err := s.repo.GetAll(ctx, categoryID)
	if err != nil {
		return nil, err
	}

	for i, p := range products {
		if p.CategoryID != "" {
			category, err := s.categoryRepo.GetByID(ctx, p.CategoryID)
			if err == nil {
				products[i].Category = &category
			}
		}
	}

	return products, nil
}

func (s *ProductService) GetProductById(ctx context.Context, id string) (entity.Product, error) {
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return entity.Product{}, err
	}

	if product.CategoryID != "" {
		category, err := s.categoryRepo.GetByID(ctx, product.CategoryID)
		if err == nil {
			product.Category = &category
		}
	}

	return product, nil
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
