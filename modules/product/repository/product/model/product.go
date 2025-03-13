package model

import "time"

type Product struct {
	ID           int64
	Name         string
	MinPrice     int
	ImageUrl     string
	GTIN         int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	SupplierInfo ProductSupplierInfo
}

type ProductSupplierInfo struct {
	SupplierID                int64
	SupplierName              string
	Price                     int
	MinSellAmount             int
	MinimumFreeDeliveryAmount float64
	DeliveryFee               float64
}

type ProductListQuery struct {
	Offset int
	Limit  int
}
