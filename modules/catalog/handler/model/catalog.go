package model

type GetProductsInput struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

type GetProductsResponse struct {
	Products []Product `json:"products"`
	Total    int       `json:"total"`
}

type GetProductInput struct {
	ID int64 `json:"id"`
}

type GetProductResponse struct {
	Product DetailedProduct `json:"product"`
}

type Product struct {
	ID           int64               `json:"id"`
	Name         string              `json:"name"`
	MinPrice     int                 `json:"min_price"`
	ImageURL     string              `json:"image_url"`
	GTIN         int64               `json:"gtin"`
	SupplierInfo ProductSupplierInfo `json:"supplier_info"`
}

type DetailedProduct struct {
	Product
	Suppliers []ProductSupplierInfo `json:"suppliers"`
}

type ProductSupplierInfo struct {
	SupplierID                int64   `json:"supplier_id"`
	Name                      string  `json:"name"`
	MinimumFreeDeliveryAmount float64 `json:"minimum_free_delivery_amount"`
	DeliveryFee               float64 `json:"delivery_fee"`
}

type PageCountInput struct {
	PageSize int `json:"pageSize"`
}

type PageCountResponse struct {
	PageCount int `json:"pageCount"`
}
