package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/shiven-lohia/interneers-lab/pkg/reports/service"
)

type ReportsHandler struct {
	service *service.ReportsService
}

func NewReportsHandler(s *service.ReportsService) *ReportsHandler {
	return &ReportsHandler{service: s}
}

func (h *ReportsHandler) CategoryCountsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	minCount, err := parseOptionalInt(r, "min_count")
	if err != nil {
		http.Error(w, "Invalid min_count: must be a non-negative integer", http.StatusBadRequest)
		return
	}
	maxCount, err := parseOptionalInt(r, "max_count")
	if err != nil {
		http.Error(w, "Invalid max_count: must be a non-negative integer", http.StatusBadRequest)
		return
	}
	if minCount != nil && maxCount != nil && *minCount > *maxCount {
		http.Error(w, "min_count must not be greater than max_count", http.StatusBadRequest)
		return
	}

	report, err := h.service.CategoryCounts(r.Context(), minCount, maxCount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *ReportsHandler) PriceDistributionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	edges, err := parseEdges(r.URL.Query().Get("buckets"))
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid buckets: %s", err.Error()), http.StatusBadRequest)
		return
	}

	report, err := h.service.PriceDistribution(r.Context(), edges)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func (h *ReportsHandler) LowStockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	threshold := 80
	if raw := r.URL.Query().Get("threshold"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v < 0 {
			http.Error(w, "Invalid threshold: must be a non-negative integer", http.StatusBadRequest)
			return
		}
		threshold = v
	}

	report, err := h.service.LowStock(r.Context(), threshold)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(report)
}

func parseOptionalInt(r *http.Request, key string) (*int, error) {
	raw := r.URL.Query().Get(key)
	if raw == "" {
		return nil, nil
	}
	v, err := strconv.Atoi(raw)
	if err != nil || v < 0 {
		return nil, fmt.Errorf("invalid value for %s", key)
	}
	return &v, nil
}

func parseEdges(raw string) ([]float64, error) {
	if raw == "" {
		return nil, nil
	}
	parts := strings.Split(raw, ",")
	edges := make([]float64, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		v, err := strconv.ParseFloat(p, 64)
		if err != nil {
			return nil, fmt.Errorf("each bucket edge must be a number, got %q", p)
		}
		edges = append(edges, v)
	}
	return edges, nil
}
