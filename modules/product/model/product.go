package model

import "time"

type ProductListQuery struct {
	Offset int
	Limit  int
}

type ProductList struct {
	Products []Product
	Total    int
}

type ProductQuery struct {
	ID int64
}

type DetailedProduct struct {
	*Product
	ProductSupplierList []ProductSupplier
}

type Product struct {
	ID                    int64
	GTIN                  int64
	Name                  string
	ImageUrl              string
	CreatedAt             time.Time
	UpdatedAt             time.Time
	LowestProductSupplier ProductSupplier
	Category              string
	Subcategory           string
}

type ProductSupplier struct {
	Price      int
	SellAmount int
	Supplier   Supplier
}

type Supplier struct {
	ID                 int64
	Name               string
	OrderAmount        int
	FreeDeliveryAmount int
	DeliveryFee        int
}

type AddProductSupplier struct{
	GTIN      string 
	Price     float64
	SupplierID int64 
}