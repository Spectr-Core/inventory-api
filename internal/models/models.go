package models

type Product struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Category string  `json:"category"`
	Stock    int     `json:"stock"`
}
type ProductStats struct {
	Sum   int     `json:"sum"`
	Count int     `json:"count"`
	Avg   float64 `json:"avg"`
	Max   float64 `json:"max"`
	Min   float64 `json:"min"`
}
type UserOrder struct {
	UserName string `json:"user_name"`
	OrderID  int    `json:"order_id"`
	Quantity int    `json:"quantity"`
}
type Order struct {
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
}
type AllOrders struct {
	UserName    string `json:"user_name"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	OrderStatus string `json:"order_status"`
}
type CategoryCount struct {
	CategoryName  string `json:"category_name"`
	CategoryCount int    `json:"category_count"`
}
