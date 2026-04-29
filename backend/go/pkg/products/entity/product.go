package entity

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	CategoryID  string  `json:"category_id" bson:"category_id"`
	Category    *ProductCategory `json:"category,omitempty"`
	Price       float64 `json:"price"`
	Brand       string  `json:"brand"`
	Quantity    int     `json:"quantity"`
}
