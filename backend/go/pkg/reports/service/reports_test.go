package service

import (
	"context"
	"errors"
	"testing"

	"github.com/shiven-lohia/interneers-lab/pkg/products/entity"
)

// MOCKS

type mockProductRepo struct {
	products []entity.Product
	err      error
}

func (m *mockProductRepo) Create(ctx context.Context, p entity.Product) (entity.Product, error) {
	return p, nil
}
func (m *mockProductRepo) GetAll(ctx context.Context, categoryID string) ([]entity.Product, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.products, nil
}
func (m *mockProductRepo) GetByID(ctx context.Context, id string) (entity.Product, error) {
	return entity.Product{}, nil
}
func (m *mockProductRepo) Update(ctx context.Context, id string, p entity.Product) (entity.Product, error) {
	return p, nil
}
func (m *mockProductRepo) Delete(ctx context.Context, id string) error { return nil }

type mockCategoryRepo struct {
	categories []entity.ProductCategory
	err        error
}

func (m *mockCategoryRepo) Create(ctx context.Context, c entity.ProductCategory) (entity.ProductCategory, error) {
	return c, nil
}
func (m *mockCategoryRepo) GetAll(ctx context.Context) ([]entity.ProductCategory, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.categories, nil
}
func (m *mockCategoryRepo) GetByID(ctx context.Context, id string) (entity.ProductCategory, error) {
	return entity.ProductCategory{}, nil
}
func (m *mockCategoryRepo) Update(ctx context.Context, id string, c entity.ProductCategory) (entity.ProductCategory, error) {
	return c, nil
}
func (m *mockCategoryRepo) Delete(ctx context.Context, id string) error { return nil }

func ptr(i int) *int { return &i }

// CATEGORY COUNTS TESTS

func TestCategoryCounts_Basic(t *testing.T) {
	svc := NewReportsService(
		&mockProductRepo{products: []entity.Product{
			{ID: "p1", Name: "Milk", CategoryID: "cat1", Price: 50, Quantity: 10, Brand: "A"},
			{ID: "p2", Name: "Bread", CategoryID: "cat1", Price: 30, Quantity: 5, Brand: "B"},
			{ID: "p3", Name: "Phone", CategoryID: "cat2", Price: 5000, Quantity: 2, Brand: "C"},
			{ID: "p4", Name: "Loose", CategoryID: "", Price: 10, Quantity: 1, Brand: "D"},
		}},
		&mockCategoryRepo{categories: []entity.ProductCategory{
			{ID: "cat1", Title: "Food"},
			{ID: "cat2", Title: "Electronics"},
		}},
	)

	report, err := svc.CategoryCounts(context.Background(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	if len(report.Categories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(report.Categories))
	}
	if report.UncategorizedCount != 1 {
		t.Fatalf("expected 1 uncategorized, got %d", report.UncategorizedCount)
	}
	// sorted descending by count
	if report.Categories[0].ProductCount != 2 {
		t.Errorf("expected first category count=2, got %d", report.Categories[0].ProductCount)
	}
}

func TestCategoryCounts_ZeroProductCategory(t *testing.T) {
	svc := NewReportsService(
		&mockProductRepo{products: []entity.Product{}},
		&mockCategoryRepo{categories: []entity.ProductCategory{
			{ID: "cat1", Title: "Food"},
		}},
	)

	report, err := svc.CategoryCounts(context.Background(), nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Categories) != 1 {
		t.Fatalf("expected 1 category (with 0 products), got %d", len(report.Categories))
	}
	if report.Categories[0].ProductCount != 0 {
		t.Errorf("expected count 0, got %d", report.Categories[0].ProductCount)
	}
}

func TestCategoryCounts_MinMaxFilter(t *testing.T) {
	svc := NewReportsService(
		&mockProductRepo{products: []entity.Product{
			{ID: "p1", CategoryID: "cat1", Price: 10, Quantity: 1, Name: "A", Brand: "X"},
			{ID: "p2", CategoryID: "cat1", Price: 10, Quantity: 1, Name: "B", Brand: "X"},
			{ID: "p3", CategoryID: "cat1", Price: 10, Quantity: 1, Name: "C", Brand: "X"},
			{ID: "p4", CategoryID: "cat2", Price: 10, Quantity: 1, Name: "D", Brand: "X"},
			{ID: "p5", CategoryID: "cat3", Price: 10, Quantity: 1, Name: "E", Brand: "X"},
			{ID: "p6", CategoryID: "cat3", Price: 10, Quantity: 1, Name: "F", Brand: "X"},
			{ID: "p7", CategoryID: "cat3", Price: 10, Quantity: 1, Name: "G", Brand: "X"},
			{ID: "p8", CategoryID: "cat3", Price: 10, Quantity: 1, Name: "H", Brand: "X"},
			{ID: "p9", CategoryID: "cat3", Price: 10, Quantity: 1, Name: "I", Brand: "X"},
		}},
		&mockCategoryRepo{categories: []entity.ProductCategory{
			{ID: "cat1", Title: "Food"},   // 3
			{ID: "cat2", Title: "Drinks"}, // 1
			{ID: "cat3", Title: "Tech"},   // 5
		}},
	)

	// min=2, max=4 → only cat1 (3) passes
	report, err := svc.CategoryCounts(context.Background(), ptr(2), ptr(4))
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Categories) != 1 {
		t.Fatalf("expected 1 category with min=2,max=4, got %d", len(report.Categories))
	}
	if report.Categories[0].CategoryID != "cat1" {
		t.Errorf("expected cat1, got %s", report.Categories[0].CategoryID)
	}
}

func TestCategoryCounts_MinFilterBoundaryInclusive(t *testing.T) {
	svc := NewReportsService(
		&mockProductRepo{products: []entity.Product{
			{ID: "p1", CategoryID: "cat1", Price: 10, Quantity: 1, Name: "A", Brand: "X"},
			{ID: "p2", CategoryID: "cat1", Price: 10, Quantity: 1, Name: "B", Brand: "X"},
		}},
		&mockCategoryRepo{categories: []entity.ProductCategory{
			{ID: "cat1", Title: "Food"},
		}},
	)

	// min=2 → cat1 (2 products) should be included (boundary inclusive)
	report, err := svc.CategoryCounts(context.Background(), ptr(2), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Categories) != 1 {
		t.Fatalf("expected 1 category at min boundary, got %d", len(report.Categories))
	}
}

// PRICE DISTRIBUTION TESTS

func TestPriceDistribution_DefaultEdges(t *testing.T) {
	svc := NewReportsService(
		&mockProductRepo{products: []entity.Product{
			{ID: "p1", CategoryID: "cat1", Price: 50, Quantity: 1, Name: "A", Brand: "X"},
			{ID: "p2", CategoryID: "cat1", Price: 200, Quantity: 1, Name: "B", Brand: "X"},
		}},
		&mockCategoryRepo{categories: []entity.ProductCategory{
			{ID: "cat1", Title: "Food"},
		}},
	)

	report, err := svc.PriceDistribution(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	// default 4 edges → 5 buckets
	if len(report.Buckets) != 5 {
		t.Fatalf("expected 5 buckets, got %d", len(report.Buckets))
	}
	if len(report.Categories) != 1 {
		t.Fatalf("expected 1 category, got %d", len(report.Categories))
	}
	counts := report.Categories[0].Counts
	if counts[0] != 1 { // 50 → bucket 0 (0-100)
		t.Errorf("expected 1 in bucket 0, got %d", counts[0])
	}
	if counts[1] != 1 { // 200 → bucket 1 (100-500)
		t.Errorf("expected 1 in bucket 1, got %d", counts[1])
	}
}

func TestPriceDistribution_EdgeAtBoundary(t *testing.T) {
	svc := NewReportsService(
		&mockProductRepo{products: []entity.Product{
			// price == 100: should go into bucket 1 (100-500), not bucket 0 (0-100)
			{ID: "p1", CategoryID: "cat1", Price: 100, Quantity: 1, Name: "A", Brand: "X"},
		}},
		&mockCategoryRepo{categories: []entity.ProductCategory{
			{ID: "cat1", Title: "Food"},
		}},
	)

	report, err := svc.PriceDistribution(context.Background(), []float64{100, 500})
	if err != nil {
		t.Fatal(err)
	}
	counts := report.Categories[0].Counts
	if counts[0] != 0 {
		t.Errorf("expected 0 in bucket 0, got %d (price=100 should not be in [0,100))", counts[0])
	}
	if counts[1] != 1 {
		t.Errorf("expected 1 in bucket 1 (100-500), got %d", counts[1])
	}
}

func TestPriceDistribution_InvalidEdges(t *testing.T) {
	svc := NewReportsService(&mockProductRepo{}, &mockCategoryRepo{})

	tests := []struct {
		name  string
		edges []float64
	}{
		{"unsorted", []float64{500, 100}},
		{"duplicate", []float64{100, 100, 500}},
		{"negative", []float64{-1, 100}},
		{"zero", []float64{0, 100}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.PriceDistribution(context.Background(), tt.edges)
			if err == nil {
				t.Errorf("expected error for edges %v, got nil", tt.edges)
			}
		})
	}
}

// LOW STOCK TESTS

func TestLowStock_Basic(t *testing.T) {
	svc := NewReportsService(
		&mockProductRepo{products: []entity.Product{
			{ID: "p1", Name: "Milk", CategoryID: "cat1", Quantity: 3, Price: 10, Brand: "X"},
			{ID: "p2", Name: "Bread", CategoryID: "cat1", Quantity: 15, Price: 10, Brand: "X"},
			{ID: "p3", Name: "Phone", CategoryID: "cat2", Quantity: 2, Price: 10, Brand: "X"},
		}},
		&mockCategoryRepo{categories: []entity.ProductCategory{
			{ID: "cat1", Title: "Food"},
			{ID: "cat2", Title: "Electronics"},
		}},
	)

	report, err := svc.LowStock(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Products) != 2 {
		t.Fatalf("expected 2 low-stock products, got %d", len(report.Products))
	}
	// sorted ascending by quantity: Phone(2) < Milk(3)
	if report.Products[0].Name != "Phone" {
		t.Errorf("expected Phone first (qty=2), got %s", report.Products[0].Name)
	}
}

func TestLowStock_ThresholdBoundary(t *testing.T) {
	svc := NewReportsService(
		&mockProductRepo{products: []entity.Product{
			{ID: "p1", Name: "Milk", CategoryID: "cat1", Quantity: 10, Price: 10, Brand: "X"},
		}},
		&mockCategoryRepo{categories: []entity.ProductCategory{
			{ID: "cat1", Title: "Food"},
		}},
	)

	// qty=10 with threshold=10 → NOT low (< not <=)
	report, err := svc.LowStock(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Products) != 0 {
		t.Errorf("expected 0 low-stock products at boundary (qty==threshold), got %d", len(report.Products))
	}
}

func TestLowStock_CategoryThresholdStrict(t *testing.T) {
	svc := NewReportsService(
		&mockProductRepo{products: []entity.Product{
			// cat1: 10 products, exactly 1 low → 10%, should NOT appear (must be >10%)
			{ID: "p1", CategoryID: "cat1", Quantity: 5, Price: 10, Name: "A", Brand: "X"},
			{ID: "p2", CategoryID: "cat1", Quantity: 20, Price: 10, Name: "B", Brand: "X"},
			{ID: "p3", CategoryID: "cat1", Quantity: 20, Price: 10, Name: "C", Brand: "X"},
			{ID: "p4", CategoryID: "cat1", Quantity: 20, Price: 10, Name: "D", Brand: "X"},
			{ID: "p5", CategoryID: "cat1", Quantity: 20, Price: 10, Name: "E", Brand: "X"},
			{ID: "p6", CategoryID: "cat1", Quantity: 20, Price: 10, Name: "F", Brand: "X"},
			{ID: "p7", CategoryID: "cat1", Quantity: 20, Price: 10, Name: "G", Brand: "X"},
			{ID: "p8", CategoryID: "cat1", Quantity: 20, Price: 10, Name: "H", Brand: "X"},
			{ID: "p9", CategoryID: "cat1", Quantity: 20, Price: 10, Name: "I", Brand: "X"},
			{ID: "p10", CategoryID: "cat1", Quantity: 20, Price: 10, Name: "J", Brand: "X"},
			// cat2: 2 products, both low → 100%, should appear
			{ID: "p11", CategoryID: "cat2", Quantity: 1, Price: 10, Name: "K", Brand: "X"},
			{ID: "p12", CategoryID: "cat2", Quantity: 2, Price: 10, Name: "L", Brand: "X"},
		}},
		&mockCategoryRepo{categories: []entity.ProductCategory{
			{ID: "cat1", Title: "Food"},
			{ID: "cat2", Title: "Electronics"},
		}},
	)

	report, err := svc.LowStock(context.Background(), 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Categories) != 1 {
		t.Fatalf("expected 1 low-stock category (>10%%), got %d", len(report.Categories))
	}
	if report.Categories[0].CategoryID != "cat2" {
		t.Errorf("expected cat2, got %s", report.Categories[0].CategoryID)
	}
}

func TestLowStock_RepoError(t *testing.T) {
	svc := NewReportsService(
		&mockProductRepo{err: errors.New("db error")},
		&mockCategoryRepo{},
	)
	_, err := svc.LowStock(context.Background(), 10)
	if err == nil {
		t.Error("expected error from repo, got nil")
	}
}
