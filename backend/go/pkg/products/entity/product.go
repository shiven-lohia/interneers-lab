package entity

type Product struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Category string `json:"category"`
	Price float64 `json:"price"`
	Brand string `json:"brand"`
	Quantity int `json:"quantity"`
}