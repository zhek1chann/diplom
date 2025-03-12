package model

type ProductInput struct {
	ID int64 `json:"id"`
}

type ProductResponse struct {
	DetailedProduct *DetailedProduct `json:"product"`
}

type Product struct {
	ID           int64                `json:"id"`
	Name         string               `json:"name"`
	MinPrice     int                  `json:"min_price"`
	ImageURL     string               `json:"image_url"`
	GTIN         int64                `json:"gtin"`
	SupplierInfo ProductSupplierInfo `json:"supplier_info"`
}

type DetailedProduct struct {
	*Product
	SupplierList []ProductSupplierInfo `json:"suppliers"`
}
