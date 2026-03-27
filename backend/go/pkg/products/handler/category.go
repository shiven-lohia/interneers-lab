package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/shiven-lohia/interneers-lab/pkg/products/entity"
	"github.com/shiven-lohia/interneers-lab/pkg/products/service"
)

type ProductCategoryHandler struct {
	service *service.ProductCategoryService
}

func NewProductCategoryHandler(service *service.ProductCategoryService) *ProductCategoryHandler {
	return &ProductCategoryHandler{
		service: service,
	}
}

func RegisterCategoryRoutes(mux *http.ServeMux, h *ProductCategoryHandler) {
	mux.HandleFunc("/categories", h.CategoriesHandler)
	mux.HandleFunc("/categories/", h.CategoriesHandler)
}

func (h *ProductCategoryHandler) CategoriesHandler(w http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(r.URL.Path, "/categories/") {

		switch r.Method {

		case http.MethodGet:
			h.GetCategoryByIDHandler(w, r)

		case http.MethodPut:
			h.UpdateCategoryHandler(w, r)

		case http.MethodDelete:
			h.DeleteCategoryHandler(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		return
	}

	switch r.Method {

	case http.MethodGet:
		h.GetCategoriesHandler(w, r)

	case http.MethodPost:
		h.CreateCategoryHandler(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ProductCategoryHandler) GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.GetAllCategories(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

func (h *ProductCategoryHandler) CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var category entity.ProductCategory

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	createdCategory, err := h.service.CreateCategory(r.Context(), category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdCategory)
}

func (h *ProductCategoryHandler) GetCategoryByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/categories/")

	category, err := h.service.GetCategoryByID(r.Context(), id)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func (h *ProductCategoryHandler) UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/categories/")

	var category entity.ProductCategory
	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		http.Error(w, "Invalid Request Body", http.StatusBadRequest)
		return
	}

	updatedCategory, err := h.service.UpdateCategory(r.Context(), id, category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCategory)
}

func (h *ProductCategoryHandler) DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/categories/")

	err := h.service.DeleteCategory(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}