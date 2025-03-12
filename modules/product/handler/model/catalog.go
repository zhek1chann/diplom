package model

type ProductListInput struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type ProductListResponse struct {
	ProductList []Product `json:"product_list"`
	Total    int       `json:"total"`
}



type ProductSupplierInfo struct {
	SupplierID                int64   `json:"supplier_id"`
	Name                      string  `json:"name"`
	MinimumFreeDeliveryAmount float64 `json:"minimum_free_delivery_amount"`
	DeliveryFee               float64 `json:"delivery_fee"`
}

