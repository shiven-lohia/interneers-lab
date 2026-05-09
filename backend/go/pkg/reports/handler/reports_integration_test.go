package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/shiven-lohia/interneers-lab/pkg/products/entity"
	"github.com/shiven-lohia/interneers-lab/pkg/products/repository"
	"github.com/shiven-lohia/interneers-lab/pkg/products/service"
	reportsEntity "github.com/shiven-lohia/interneers-lab/pkg/reports/entity"
	reportsService "github.com/shiven-lohia/interneers-lab/pkg/reports/service"
)

func setupTest(t *testing.T) (*ReportsHandler, *repository.MongoProductRepository) {
	t.Helper()

	client, err := repository.ConnectMongo()
	if err != nil {
		t.Fatal(err)
	}

	productRepo := repository.NewMongoProductRepository(client)
	categoryCollection := client.Database("inventory").Collection("categories")
	categoryRepo := repository.NewMongoCategoryRepository(categoryCollection)

	rService := reportsService.NewReportsService(productRepo, categoryRepo)
	rHandler := NewReportsHandler(rService)

	return rHandler, productRepo
}

func seedData(t *testing.T, productRepo *repository.MongoProductRepository) {
	t.Helper()
	ctx := context.Background()

	categoryCollection := productRepo.Collection().Database().Collection("categories")

	productRepo.Collection().DeleteMany(ctx, map[string]interface{}{})
	categoryCollection.DeleteMany(ctx, map[string]interface{}{})

	t.Cleanup(func() {
		productRepo.Collection().DeleteMany(ctx, map[string]interface{}{})
		categoryCollection.DeleteMany(ctx, map[string]interface{}{})
	})

	// Create categories via product service so they appear in the DB
	catColl := productRepo.Collection().Database().Collection("categories")
	cats := []entity.ProductCategory{
		{ID: "cat1", Title: "Food"},
		{ID: "cat2", Title: "Electronics"},
	}
	for _, cat := range cats {
		catColl.InsertOne(ctx, cat)
	}

	// Create products via product service
	pService := service.NewProductService(productRepo, repository.NewMongoCategoryRepository(catColl))
	products := []entity.Product{
		{Name: "Milk", CategoryID: "cat1", Price: 50, Quantity: 3, Brand: "Amul"},
		{Name: "Bread", CategoryID: "cat1", Price: 30, Quantity: 15, Brand: "Britannia"},
		{Name: "Phone", CategoryID: "cat2", Price: 5000, Quantity: 2, Brand: "Samsung"},
		{Name: "Laptop", CategoryID: "cat2", Price: 50000, Quantity: 20, Brand: "Dell"},
	}
	for _, p := range products {
		_, err := pService.CreateProduct(ctx, p)
		if err != nil {
			t.Fatalf("failed to seed product %q: %v", p.Name, err)
		}
	}
}

func TestIntegration_CategoryCounts(t *testing.T) {
	rHandler, productRepo := setupTest(t)
	seedData(t, productRepo)

	req := httptest.NewRequest(http.MethodGet, "/reports/category-counts", nil)
	w := httptest.NewRecorder()

	rHandler.CategoryCountsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var report reportsEntity.CategoryCountsReport
	if err := json.Unmarshal(w.Body.Bytes(), &report); err != nil {
		t.Fatal(err)
	}

	if len(report.Categories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(report.Categories))
	}
}

func TestIntegration_CategoryCounts_MinFilter(t *testing.T) {
	rHandler, productRepo := setupTest(t)
	seedData(t, productRepo)

	// cat1 has 2 products, cat2 has 2 products — min_count=2, max_count=2 → both included
	req := httptest.NewRequest(http.MethodGet, "/reports/category-counts?min_count=2&max_count=2", nil)
	w := httptest.NewRecorder()

	rHandler.CategoryCountsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var report reportsEntity.CategoryCountsReport
	json.Unmarshal(w.Body.Bytes(), &report)

	if len(report.Categories) != 2 {
		t.Fatalf("expected 2 categories with min=max=2, got %d", len(report.Categories))
	}
}

func TestIntegration_PriceDistribution(t *testing.T) {
	rHandler, productRepo := setupTest(t)
	seedData(t, productRepo)

	req := httptest.NewRequest(http.MethodGet, "/reports/price-distribution?buckets=100,1000,10000", nil)
	w := httptest.NewRecorder()

	rHandler.PriceDistributionHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var report reportsEntity.PriceDistributionReport
	if err := json.Unmarshal(w.Body.Bytes(), &report); err != nil {
		t.Fatal(err)
	}

	// 3 edges → 4 buckets
	if len(report.Buckets) != 4 {
		t.Fatalf("expected 4 buckets, got %d", len(report.Buckets))
	}
	if len(report.Categories) != 2 {
		t.Fatalf("expected 2 categories, got %d", len(report.Categories))
	}
}

func TestIntegration_PriceDistribution_InvalidBuckets(t *testing.T) {
	rHandler, productRepo := setupTest(t)
	seedData(t, productRepo)

	req := httptest.NewRequest(http.MethodGet, "/reports/price-distribution?buckets=500,100", nil)
	w := httptest.NewRecorder()

	rHandler.PriceDistributionHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for unsorted buckets, got %d", w.Code)
	}
}

func TestIntegration_LowStock(t *testing.T) {
	rHandler, productRepo := setupTest(t)
	seedData(t, productRepo)

	// Milk(qty=3) and Phone(qty=2) are below threshold=10
	req := httptest.NewRequest(http.MethodGet, "/reports/low-stock?threshold=10", nil)
	w := httptest.NewRecorder()

	rHandler.LowStockHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var report reportsEntity.LowStockReport
	if err := json.Unmarshal(w.Body.Bytes(), &report); err != nil {
		t.Fatal(err)
	}

	if len(report.Products) != 2 {
		t.Fatalf("expected 2 low-stock products, got %d", len(report.Products))
	}
	if report.Products[0].Quantity > report.Products[1].Quantity {
		t.Errorf("expected products sorted ascending by quantity")
	}
}

func TestIntegration_LowStock_InvalidThreshold(t *testing.T) {
	rHandler, productRepo := setupTest(t)
	seedData(t, productRepo)

	req := httptest.NewRequest(http.MethodGet, "/reports/low-stock?threshold=abc", nil)
	w := httptest.NewRecorder()

	rHandler.LowStockHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid threshold, got %d", w.Code)
	}
}

func TestIntegration_MethodNotAllowed(t *testing.T) {
	rHandler, productRepo := setupTest(t)
	seedData(t, productRepo)

	req := httptest.NewRequest(http.MethodPost, "/reports/category-counts", bytes.NewBuffer([]byte("{}")))
	w := httptest.NewRecorder()

	rHandler.CategoryCountsHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", w.Code)
	}
}
