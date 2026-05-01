package service

import (
	"context"
	"errors"
	"fmt"
	"sort"

	reportEntity "github.com/shiven-lohia/interneers-lab/pkg/reports/entity"
	"github.com/shiven-lohia/interneers-lab/pkg/products/entity"
	"github.com/shiven-lohia/interneers-lab/pkg/products/repository"
)

var defaultBucketEdges = []float64{100, 500, 1000, 5000}

type ReportsService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.ProductCategoryRepository
}

func NewReportsService(p repository.ProductRepository, c repository.ProductCategoryRepository) *ReportsService {
	return &ReportsService{productRepo: p, categoryRepo: c}
}

func (s *ReportsService) CategoryCounts(ctx context.Context, minCount, maxCount *int) (reportEntity.CategoryCountsReport, error) {
	products, err := s.productRepo.GetAll(ctx, "")
	if err != nil {
		return reportEntity.CategoryCountsReport{}, err
	}
	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		return reportEntity.CategoryCountsReport{}, err
	}

	counts := make(map[string]int)
	uncategorized := 0
	for _, p := range products {
		if p.CategoryID == "" {
			uncategorized++
		} else {
			counts[p.CategoryID]++
		}
	}

	catMap := make(map[string]entity.ProductCategory, len(categories))
	for _, c := range categories {
		catMap[c.ID] = c
	}

	var result []reportEntity.CategoryCount
	for _, cat := range categories {
		count := counts[cat.ID]
		if minCount != nil && count < *minCount {
			continue
		}
		if maxCount != nil && count > *maxCount {
			continue
		}
		result = append(result, reportEntity.CategoryCount{
			CategoryID:    cat.ID,
			CategoryTitle: cat.Title,
			ProductCount:  count,
		})
	}
	if result == nil {
		result = []reportEntity.CategoryCount{}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ProductCount > result[j].ProductCount
	})

	return reportEntity.CategoryCountsReport{
		Categories:         result,
		UncategorizedCount: uncategorized,
		Filters:            reportEntity.CountFilters{MinCount: minCount, MaxCount: maxCount},
	}, nil
}

func (s *ReportsService) PriceDistribution(ctx context.Context, edges []float64) (reportEntity.PriceDistributionReport, error) {
	if len(edges) == 0 {
		edges = defaultBucketEdges
	}
	if err := validateEdges(edges); err != nil {
		return reportEntity.PriceDistributionReport{}, err
	}

	products, err := s.productRepo.GetAll(ctx, "")
	if err != nil {
		return reportEntity.PriceDistributionReport{}, err
	}
	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		return reportEntity.PriceDistributionReport{}, err
	}

	buckets := buildBuckets(edges)

	catCounts := make(map[string][]int)
	for _, cat := range categories {
		catCounts[cat.ID] = make([]int, len(buckets))
	}
	// track uncategorized too but not included in report
	for _, p := range products {
		if p.CategoryID == "" {
			continue
		}
		if _, ok := catCounts[p.CategoryID]; !ok {
			catCounts[p.CategoryID] = make([]int, len(buckets))
		}
		idx := bucketIndex(p.Price, edges)
		catCounts[p.CategoryID][idx]++
	}

	catMap := make(map[string]entity.ProductCategory, len(categories))
	for _, c := range categories {
		catMap[c.ID] = c
	}

	var result []reportEntity.CategoryPriceCounts
	for _, cat := range categories {
		result = append(result, reportEntity.CategoryPriceCounts{
			CategoryID:    cat.ID,
			CategoryTitle: cat.Title,
			Counts:        catCounts[cat.ID],
		})
	}
	if result == nil {
		result = []reportEntity.CategoryPriceCounts{}
	}

	return reportEntity.PriceDistributionReport{
		Buckets:    buckets,
		Categories: result,
	}, nil
}

func (s *ReportsService) LowStock(ctx context.Context, threshold int) (reportEntity.LowStockReport, error) {
	products, err := s.productRepo.GetAll(ctx, "")
	if err != nil {
		return reportEntity.LowStockReport{}, err
	}
	categories, err := s.categoryRepo.GetAll(ctx)
	if err != nil {
		return reportEntity.LowStockReport{}, err
	}

	catMap := make(map[string]entity.ProductCategory, len(categories))
	for _, c := range categories {
		catMap[c.ID] = c
	}

	type catStats struct {
		total    int
		lowStock int
	}
	stats := make(map[string]*catStats)
	for _, cat := range categories {
		stats[cat.ID] = &catStats{}
	}

	var lowProducts []reportEntity.LowStockProduct
	for _, p := range products {
		if p.CategoryID != "" {
			if _, ok := stats[p.CategoryID]; !ok {
				stats[p.CategoryID] = &catStats{}
			}
			stats[p.CategoryID].total++
		}
		if p.Quantity < threshold {
			catTitle := ""
			if cat, ok := catMap[p.CategoryID]; ok {
				catTitle = cat.Title
			}
			lowProducts = append(lowProducts, reportEntity.LowStockProduct{
				ID:            p.ID,
				Name:          p.Name,
				Quantity:      p.Quantity,
				CategoryID:    p.CategoryID,
				CategoryTitle: catTitle,
			})
			if p.CategoryID != "" {
				stats[p.CategoryID].lowStock++
			}
		}
	}
	if lowProducts == nil {
		lowProducts = []reportEntity.LowStockProduct{}
	}

	sort.Slice(lowProducts, func(i, j int) bool {
		return lowProducts[i].Quantity < lowProducts[j].Quantity
	})

	var lowCats []reportEntity.LowStockCategory
	for _, cat := range categories {
		st := stats[cat.ID]
		if st == nil || st.total == 0 {
			continue
		}
		pct := float64(st.lowStock) / float64(st.total) * 100
		if pct <= 10 {
			continue
		}
		lowCats = append(lowCats, reportEntity.LowStockCategory{
			CategoryID:    cat.ID,
			CategoryTitle: cat.Title,
			LowStockCount: st.lowStock,
			TotalCount:    st.total,
			Percentage:    roundTo1(pct),
		})
	}
	if lowCats == nil {
		lowCats = []reportEntity.LowStockCategory{}
	}

	sort.Slice(lowCats, func(i, j int) bool {
		return lowCats[i].Percentage > lowCats[j].Percentage
	})

	return reportEntity.LowStockReport{
		Threshold:  threshold,
		Products:   lowProducts,
		Categories: lowCats,
	}, nil
}

func validateEdges(edges []float64) error {
	for i, e := range edges {
		if e <= 0 {
			return errors.New("bucket edges must be positive")
		}
		if i > 0 && e <= edges[i-1] {
			return errors.New("bucket edges must be strictly increasing")
		}
	}
	return nil
}

func buildBuckets(edges []float64) []reportEntity.PriceBucket {
	buckets := make([]reportEntity.PriceBucket, 0, len(edges)+1)
	prev := 0.0
	for _, e := range edges {
		eCopy := e
		buckets = append(buckets, reportEntity.PriceBucket{
			Label: fmt.Sprintf("%.0f-%.0f", prev, e),
			Min:   prev,
			Max:   &eCopy,
		})
		prev = e
	}
	buckets = append(buckets, reportEntity.PriceBucket{
		Label: fmt.Sprintf("%.0f+", prev),
		Min:   prev,
		Max:   nil,
	})
	return buckets
}

func bucketIndex(price float64, edges []float64) int {
	for i, e := range edges {
		if price < e {
			return i
		}
	}
	return len(edges)
}

func roundTo1(f float64) float64 {
	return float64(int(f*10+0.5)) / 10
}
