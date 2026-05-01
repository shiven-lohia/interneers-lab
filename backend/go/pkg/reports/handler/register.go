package handler

import "net/http"

func RegisterRoutes(mux *http.ServeMux, h *ReportsHandler) {
	mux.HandleFunc("/reports/category-counts", h.CategoryCountsHandler)
	mux.HandleFunc("/reports/price-distribution", h.PriceDistributionHandler)
	mux.HandleFunc("/reports/low-stock", h.LowStockHandler)
}
