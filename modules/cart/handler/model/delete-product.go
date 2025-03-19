package model

type DeleteProductFromCartInput struct {
	UserID     int64 `json:"user_id"`
	ProductID  int64 `json:"product_id"`
	SupplierID int64 `json:"supplier_id"`
	Quantity   int   `json:"quantity"`
}
