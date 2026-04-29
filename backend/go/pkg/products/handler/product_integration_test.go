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
)

func setupTest(t *testing.T) (*ProductHandler, *ProductCategoryHandler, *repository.MongoProductRepository) {

	client, err := repository.ConnectMongo()
	if err != nil {
		t.Fatal(err)
	}

	productRepo := repository.NewMongoProductRepository(client)

	categoryCollection := client.Database("inventory").Collection("categories")
	categoryRepo := repository.NewMongoCategoryRepository(categoryCollection)

	productService := service.NewProductService(productRepo, categoryRepo)
	categoryService := service.NewProductCategoryService(categoryRepo)

	productHandler := NewProductHandler(productService)
	categoryHandler := NewProductCategoryHandler(categoryService)

	return productHandler, categoryHandler, productRepo
}

func TestIntegration_CreateAndGetProduct(t *testing.T) {

	pHandler, cHandler, productRepo := setupTest(t)

	ctx := context.Background()

	categoryCollection := productRepo.Collection().Database().Collection("categories")

	_, err := productRepo.Collection().DeleteMany(ctx, map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	_, err = categoryCollection.DeleteMany(ctx, map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		productRepo.Collection().DeleteMany(ctx, map[string]interface{}{})
		categoryCollection.DeleteMany(ctx, map[string]interface{}{})
	})

	// ---- CREATE CATEGORY ----

	category := entity.ProductCategory{
		ID:    "cat1",
		Title: "Food",
	}

	body, _ := json.Marshal(category)

	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	cHandler.CategoriesHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	// ---- CREATE PRODUCT ----

	product := entity.Product{
		ID:         "p1",
		Name:       "Milk",
		Price:      50,
		Quantity:   10,
		Brand:      "Amul",
		CategoryID: "cat1",
	}

	body, _ = json.Marshal(product)

	req = httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	w = httptest.NewRecorder()

	pHandler.ProductsHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	// ---- GET PRODUCTS ----

	req = httptest.NewRequest(http.MethodGet, "/products", nil)
	w = httptest.NewRecorder()

	pHandler.ProductsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var products []entity.Product
	err = json.Unmarshal(w.Body.Bytes(), &products)
	if err != nil {
		t.Fatal(err)
	}

	if len(products) != 1 {
		t.Fatalf("expected 1 product, got %d", len(products))
	}

	if products[0].Name != "Milk" {
		t.Fatalf("expected product name Milk, got %s", products[0].Name)
	}
}

func TestIntegration_FilterProductsByCategory(t *testing.T) {

	pHandler, cHandler, productRepo := setupTest(t)

	ctx := context.Background()

	categoryCollection := productRepo.Collection().Database().Collection("categories")

	productRepo.Collection().DeleteMany(ctx, map[string]interface{}{})
	categoryCollection.DeleteMany(ctx, map[string]interface{}{})

	t.Cleanup(func() {
		productRepo.Collection().DeleteMany(ctx, map[string]interface{}{})
		categoryCollection.DeleteMany(ctx, map[string]interface{}{})
	})

	// ---- CREATE CATEGORIES ----

	categories := []entity.ProductCategory{
		{ID: "cat1", Title: "Food"},
		{ID: "cat2", Title: "Electronics"},
	}

	for _, cat := range categories {
		body, _ := json.Marshal(cat)
		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		cHandler.CategoriesHandler(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("failed to create category %s", cat.ID)
		}
	}

	// ---- CREATE PRODUCTS ----

	products := []entity.Product{
		{
			ID: "p1", Name: "Milk", Price: 50, Quantity: 10, Brand: "Amul", CategoryID: "cat1",
		},
		{
			ID: "p2", Name: "Bread", Price: 30, Quantity: 5, Brand: "Britannia", CategoryID: "cat1",
		},
		{
			ID: "p3", Name: "Phone", Price: 5000, Quantity: 2, Brand: "Samsung", CategoryID: "cat2",
		},
	}

	for _, p := range products {
		body, _ := json.Marshal(p)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		pHandler.ProductsHandler(w, req)

		if w.Code != http.StatusCreated {
			t.Fatalf("failed to create product %s", p.ID)
		}
	}

	// ---- FILTER PRODUCTS ----

	req := httptest.NewRequest(http.MethodGet, "/products?category_id=cat1", nil)
	w := httptest.NewRecorder()

	pHandler.ProductsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var result []entity.Product
	err := json.Unmarshal(w.Body.Bytes(), &result)
	if err != nil {
		t.Fatal(err)
	}

	// ---- ASSERT ----

	if len(result) != 2 {
		t.Fatalf("expected 2 products, got %d", len(result))
	}

	for _, p := range result {
		if p.CategoryID != "cat1" {
			t.Fatalf("expected only cat1 products, got category %s", p.CategoryID)
		}
	}
}

func TestIntegration_CreateProduct_InvalidCategory(t *testing.T) {

	pHandler, _, productRepo := setupTest(t)

	ctx := context.Background()

	categoryCollection := productRepo.Collection().Database().Collection("categories")

	productRepo.Collection().DeleteMany(ctx, map[string]interface{}{})
	categoryCollection.DeleteMany(ctx, map[string]interface{}{})

	t.Cleanup(func() {
		productRepo.Collection().DeleteMany(ctx, map[string]interface{}{})
		categoryCollection.DeleteMany(ctx, map[string]interface{}{})
	})

	// ---- CREATE PRODUCT WITH INVALID CATEGORY ----

	product := entity.Product{
		ID:         "p1",
		Name:       "Milk",
		Price:      50,
		Quantity:   10,
		Brand:      "Amul",
		CategoryID: "non-existent",
	}

	body, _ := json.Marshal(product)

	req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	pHandler.ProductsHandler(w, req)

	// ---- ASSERT RESPONSE ----

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}

	// ---- VERIFY DB: SHOULD BE EMPTY ----

	req = httptest.NewRequest(http.MethodGet, "/products", nil)
	w = httptest.NewRecorder()

	pHandler.ProductsHandler(w, req)

	var products []entity.Product
	err := json.Unmarshal(w.Body.Bytes(), &products)
	if err != nil {
		t.Fatal(err)
	}

	if len(products) != 0 {
		t.Fatalf("expected 0 products, got %d", len(products))
	}
}