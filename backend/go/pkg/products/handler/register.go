package handler

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *ProductHandler) {

	mux.HandleFunc("/products", h.ProductsHandler)
	mux.HandleFunc("/products/", h.ProductsHandler)

}