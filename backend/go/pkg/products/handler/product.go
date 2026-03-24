package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"strconv"
	"encoding/csv"

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

func (h *ProductHandler) BulkCreateProductsHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		http.Error(w, "Error parsing form data", http.StatusBadRequest)
		return
	}

	file, _ ,err := r.FormFile("file")
	if(err!=nil) {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, "Invalid CSV", http.StatusBadRequest)
		return
	}

	var products []entity.Product
	
	for i, row := range records {
		if i == 0 {
			continue
		}
		if(len(row) < 5) {
			continue
		}

		price, err := strconv.ParseFloat(row[2], 64)
		if(err!=nil) {
			continue
		}

		qty, err := strconv.Atoi(row[3])
		if err != nil {
			continue
		}

		product := entity.Product{
			ID:          row[0],
			Name:        row[1],
			Price: 	     price,
			Quantity:    qty,
			Brand: 	     row[4],
		}
		products = append(products, product)
	}

	createdProducts, err := h.controller.BulkCreateProducts(r.Context(), products)
	if(err!=nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdProducts)
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
	products, _ := h.controller.GetAllProducts(r.Context())

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

	createdProduct, err := h.controller.CreateProduct(r.Context(), product)
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

	product, err := h.controller.GetProductById(r.Context(), id)
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

	updatedProduct, err := h.controller.UpdateProduct(r.Context(), id, product)
	if(err!=nil) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProduct)
}

func (h *ProductHandler) DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/products/")

	err := h.controller.DeleteProduct(r.Context(), id)
	if(err!=nil) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}