package model

type AddProductToCartInput struct {
	Quantity   int   `json:"quantity"`
	Price      int   `json:"price"`
	ProductID  int64 `json:"product_id"`
	SupplierID int64 `json:"supplier_id"`
	UserID     int64 `json:"user_id"`
}

type GetCardInput struct {
	UserID int64 `json:"user_id"`
}

type GetCartResponse struct {
	Total     int        `json:"total"`
	UserID    int64      `json:"user_id"`
	Suppliers []Supplier `json:"suppliers"`
}

type Supplier struct {
	MinOrderAmount     int       `json:"min_order_amount`
	TotalAmount        int       `json:"total_amount"`
	FreeDeliveryAmount int       `json:"free_delivery_amount"`
	DeliveryFee        int       `json:"delivery_fee"`
	ID                 int64     `json:"id"`
	Name               string    `json:"name"`
	ProductList        []Product `json:"product_list"`
}

type Product struct {
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Image    string `json:"image"`
}
