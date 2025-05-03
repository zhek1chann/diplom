package model

import "errors"

var (
	ErrUnauthorized = errors.New("api: unauthorized")
)

type ErrorResponse struct {
	Err string `json:"error"`
}

type GetOrdersResponse struct {
	Orders []Order `json:"orders"`
}

type Order struct {
	ID          int64      `json:"id"`
	Status      string     `json:"status"`
	OrderDate   string     `json:"order_date"`
	Supplier    *Supplier  `json:"supplier"`
	ProductList []*Product `json:"product_list"`
}
type Product struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	ImageUrl string `json:"image"`
	Quantity int    `json:"quantity"`
	Price    int    `json:"price"`
}
type Supplier struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
