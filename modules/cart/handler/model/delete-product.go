package model

type DeleteProductFromCartInput struct {
	CustomerID int64 `json:"customer_id"`
	ProductID  int64 `json:"product_id"`
	SupplierID int64 `json:"supplier_id"`
	Quantity   int   `json:"quantity"`
}
