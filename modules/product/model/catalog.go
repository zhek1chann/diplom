package model

type ProductListQuery struct {
	Page     int
	PageSize int
}

type ProductList struct {
	Products []Product
	Total    int
}

type ProductQuery struct {
	ID int64
}

type Product struct {
	ID           int64
	Name         string
	MinPrice     int
	ImageURL     string
	GTIN         int64
	SupplierInfo ProductSupplierInfo
}

type DetailedProduct struct {
	*Product
	SupplierList []ProductSupplierInfo
}

type ProductSupplierInfo struct {
	SupplierID                int64
	Name                      string
	MinimumFreeDeliveryAmount float64
	DeliveryFee               float64
}
