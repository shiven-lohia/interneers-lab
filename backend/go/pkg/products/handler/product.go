package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/shiven-lohia/interneers-lab/pkg/products/entity"
	"github.com/shiven-lohia/interneers-lab/pkg/products/controller"
)

type ProductHandler struct {
	controller *controller.ProductController
}

func NewProductHandler(controller *controller.ProductController) *ProductHandler {
	return &ProductHandler{
		controller: controller,
	}
}

func (h *ProductHandler) ProductsHandler(w http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(r.URL.Path, "/products/") {

		switch r.Method {

		case http.MethodGet:
			h.GetProductByIDHandler(w, r)

		case http.MethodPut:
			h.UpdateProductHandler(w, r)

		case http.MethodDelete:
			h.DeleteProductHandler(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		return
	}

	switch r.Method {

	case http.MethodGet:
		h.GetProductHandler(w, r)

	case http.MethodPost:
		h.CreateProductHandler(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductHandler) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	products, _ := h.controller.GetAllProducts()

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var product entity.Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if(err!=nil) {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	createdProduct, err := h.controller.CreateProduct(product)
	if(err!=nil) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProduct)
}

func (h *ProductHandler) GetProductByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/products/")

	product, err := h.controller.GetProductById(id)
	if(err!=nil) {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/products/")

	var product entity.Product

	err := json.NewDecoder(r.Body).Decode(&product)
	if(err!=nil) {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	updatedProduct, err := h.controller.UpdateProduct(id, product)
	if(err!=nil) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProduct)
}

func (h *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/products/")

	err := h.controller.DeleteProduct(id)
	if(err!=nil) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}