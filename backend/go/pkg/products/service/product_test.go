package service

import (
	"context"
	"errors"
	"testing"

	"github.com/shiven-lohia/interneers-lab/pkg/products/entity"
)

// MOCKS

type MockProductRepo struct{
	shouldFailUpdate bool
}

func (m *MockProductRepo) Create(ctx context.Context, p entity.Product) (entity.Product, error) {
	return p, nil
}

func (m *MockProductRepo) GetByID(ctx context.Context, id string) (entity.Product, error) {
	return entity.Product{}, nil
}

func (m *MockProductRepo) GetAll(ctx context.Context, categoryID string) ([]entity.Product, error) {
	return []entity.Product{}, nil
}

func (m *MockProductRepo) Update(ctx context.Context, id string, p entity.Product) (entity.Product, error) {
	if m.shouldFailUpdate {
		return entity.Product{}, errors.New("not found")
	}
	return p, nil
}

func (m *MockProductRepo) Delete(ctx context.Context, id string) error {
	return nil
}


type MockCategoryRepo struct{
	shouldFail bool
}

func (m *MockCategoryRepo) GetByID(ctx context.Context, id string) (entity.ProductCategory, error) {
	if m.shouldFail {
		return entity.ProductCategory{}, errors.New("not found")
	}
	return entity.ProductCategory{ID: id}, nil
}

func (m *MockCategoryRepo) Create(ctx context.Context, c entity.ProductCategory) (entity.ProductCategory, error) {
	return c, nil
}

func (m *MockCategoryRepo) GetAll(ctx context.Context) ([]entity.ProductCategory, error) {
	return []entity.ProductCategory{}, nil
}

func (m *MockCategoryRepo) Update(ctx context.Context, id string, c entity.ProductCategory) (entity.ProductCategory, error) {
	if m.shouldFail {
		return entity.ProductCategory{}, errors.New("not found")
	}
	return c, nil
}

func (m *MockCategoryRepo) Delete(ctx context.Context, id string) error {
	if m.shouldFail {
		return errors.New("not found")
	}
	return nil
}

// TESTS

func TestCreateProduct(t *testing.T) {

	tests := []struct {
		name        string
		product     entity.Product
		categoryErr bool
		expectError bool
	}{
		{
			name: "valid product",
			product: entity.Product{
				ID: "1",
				Name: "Milk",
				Price: 50,
				Quantity: 10,
				Brand: "Amul",
				CategoryID: "cat1",
			},
			categoryErr: false,
			expectError: false,
		},
		{
			name: "missing name",
			product: entity.Product{
				ID: "1",
				Price: 50,
				Quantity: 10,
				Brand: "Amul",
			},
			expectError: true,
		},
		{
			name: "invalid price",
			product: entity.Product{
				ID: "1",
				Name: "Milk",
				Price: 0,
				Quantity: 10,
				Brand: "Amul",
			},
			expectError: true,
		},
		{
			name: "invalid quantity",
			product: entity.Product{
				ID: "1",
				Name: "Milk",
				Price: 50,
				Quantity: -1,
				Brand: "Amul",
			},
			expectError: true,
		},
		{
			name: "missing brand",
			product: entity.Product{
				ID: "1",
				Name: "Milk",
				Price: 50,
				Quantity: 10,
			},
			expectError: true,
		},
		{
			name: "invalid category",
			product: entity.Product{
				ID: "1",
				Name: "Milk",
				Price: 50,
				Quantity: 10,
				Brand: "Amul",
				CategoryID: "bad",
			},
			categoryErr: true,
			expectError: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &MockProductRepo{}
			mockCatRepo := &MockCategoryRepo{shouldFail: tt.categoryErr}

			service := NewProductService(mockRepo, mockCatRepo)

			_, err := service.CreateProduct(context.Background(), tt.product)

			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestUpdateProduct(t *testing.T) {

	tests := []struct {
		name              string
		product           entity.Product
		categoryErr       bool
		updateShouldFail  bool
		expectError       bool
	}{
		{
			name: "valid update",
			product: entity.Product{
				Name: "Milk",
				Price: 50,
				Quantity: 10,
				Brand: "Amul",
				CategoryID: "cat1",
			},
			expectError: false,
		},
		{
			name: "missing name",
			product: entity.Product{
				Price: 50,
				Quantity: 10,
				Brand: "Amul",
			},
			expectError: true,
		},
		{
			name: "invalid price",
			product: entity.Product{
				Name: "Milk",
				Price: 0,
				Quantity: 10,
				Brand: "Amul",
			},
			expectError: true,
		},
		{
			name: "invalid quantity",
			product: entity.Product{
				Name: "Milk",
				Price: 50,
				Quantity: -1,
				Brand: "Amul",
			},
			expectError: true,
		},
		{
			name: "missing brand",
			product: entity.Product{
				Name: "Milk",
				Price: 50,
				Quantity: 10,
			},
			expectError: true,
		},
		{
			name: "invalid category",
			product: entity.Product{
				Name: "Milk",
				Price: 50,
				Quantity: 10,
				Brand: "Amul",
				CategoryID: "bad",
			},
			categoryErr: true,
			expectError: true,
		},
		{
			name: "product not found",
			product: entity.Product{
				Name: "Milk",
				Price: 50,
				Quantity: 10,
				Brand: "Amul",
			},
			updateShouldFail: true,
			expectError: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &MockProductRepo{
				shouldFailUpdate: tt.updateShouldFail,
			}

			mockCatRepo := &MockCategoryRepo{
				shouldFail: tt.categoryErr,
			}

			service := NewProductService(mockRepo, mockCatRepo)

			_, err := service.UpdateProduct(context.Background(), "1", tt.product)

			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestCreateCategory(t *testing.T) {

	tests := []struct {
		name        string
		category    entity.ProductCategory
		expectError bool
	}{
		{
			name: "valid category",
			category: entity.ProductCategory{
				ID: "cat1",
				Title: "Food",
			},
			expectError: false,
		},
		{
			name: "missing title",
			category: entity.ProductCategory{
				ID: "cat1",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &MockCategoryRepo{}

			service := NewProductCategoryService(mockRepo)

			_, err := service.CreateCategory(context.Background(), tt.category)

			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}

func TestUpdateCategory(t *testing.T) {

	tests := []struct {
		name        string
		category    entity.ProductCategory
		shouldFail  bool
		expectError bool
	}{
		{
			name: "valid update",
			category: entity.ProductCategory{
				Title: "Updated",
			},
			expectError: false,
		},
		{
			name: "missing title",
			category: entity.ProductCategory{},
			expectError: true,
		},
		{
			name: "not found",
			category: entity.ProductCategory{
				Title: "Updated",
			},
			shouldFail: true,
			expectError: true,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			mockRepo := &MockCategoryRepo{
				shouldFail: tt.shouldFail,
			}

			service := NewProductCategoryService(mockRepo)

			_, err := service.UpdateCategory(context.Background(), "cat1", tt.category)

			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
		})
	}
}