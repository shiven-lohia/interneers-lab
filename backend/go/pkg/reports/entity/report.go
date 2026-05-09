package entity

type CategoryCountsReport struct {
	Categories         []CategoryCount `json:"categories"`
	UncategorizedCount int             `json:"uncategorized_count"`
	Filters            CountFilters    `json:"filters"`
}

type CategoryCount struct {
	CategoryID    string `json:"category_id"`
	CategoryTitle string `json:"category_title"`
	ProductCount  int    `json:"product_count"`
}

type CountFilters struct {
	MinCount *int `json:"min_count"`
	MaxCount *int `json:"max_count"`
}

type PriceDistributionReport struct {
	Buckets    []PriceBucket         `json:"buckets"`
	Categories []CategoryPriceCounts `json:"categories"`
}

type PriceBucket struct {
	Label string   `json:"label"`
	Min   float64  `json:"min"`
	Max   *float64 `json:"max"`
}

type CategoryPriceCounts struct {
	CategoryID    string `json:"category_id"`
	CategoryTitle string `json:"category_title"`
	Counts        []int  `json:"counts"`
}

type LowStockReport struct {
	Threshold  int               `json:"threshold"`
	Products   []LowStockProduct `json:"products"`
	Categories []LowStockCategory `json:"categories"`
}

type LowStockProduct struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	Quantity      int    `json:"quantity"`
	CategoryID    string `json:"category_id"`
	CategoryTitle string `json:"category_title"`
}

type LowStockCategory struct {
	CategoryID    string  `json:"category_id"`
	CategoryTitle string  `json:"category_title"`
	LowStockCount int     `json:"low_stock_count"`
	TotalCount    int     `json:"total_count"`
	Percentage    float64 `json:"percentage"`
}
