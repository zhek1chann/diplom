package model

type ProductListInput struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ProductListResponse struct {
	ProductList []Product `json:"product_list"`
	Total       int       `json:"total"`
}
