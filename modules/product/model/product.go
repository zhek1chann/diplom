package model

import "time"

type ProductListQuery struct {
	Offset        int
	Limit         int
	CategoryID    *int // Optional category filter
	SubcategoryID *int // Optional subcategory filter
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
	CategoryID            *int   `json:"category_id"`
	SubcategoryID         *int   `json:"subcategory_id"`
	CategoryName          string `json:"category_name,omitempty"`
	SubcategoryName       string `json:"subcategory_name,omitempty"`
}

type ProductSupplier struct {
	Price      *int
	SellAmount *int
	Supplier   Supplier
}

type Supplier struct {
	ID                 int64
	Name               *string
	OrderAmount        *int
	FreeDeliveryAmount *int
	DeliveryFee        *int
}

type AddProductSupplier struct {
	GTIN          string
	Price         float64
	SupplierID    int64
	CategoryID    *int
	SubcategoryID *int
}

type UpdateProductRequest struct {
	Name          string `json:"name,omitempty"`
	ImageUrl      string `json:"image_url,omitempty"`
	CategoryID    *int   `json:"category_id,omitempty"`
	SubcategoryID *int   `json:"subcategory_id,omitempty"`
}
