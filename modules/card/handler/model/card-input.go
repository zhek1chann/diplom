package model

type CardInput struct {
	Quantity   int   `json:"quantity"`
	Price      int   `json:"price"`
	ProductID  int64 `json:"product_id"`
	SupplierID int64 `json:"supplier_id"`
	UserID     int64 `json:"user_id"`
}
