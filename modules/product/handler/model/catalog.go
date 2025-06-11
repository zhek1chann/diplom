package model

type ProductListInput struct {
	Limit         int  `json:"limit"`
	Offset        int  `json:"offset"`
	CategoryID    *int `json:"category_id,omitempty"`
	SubcategoryID *int `json:"subcategory_id,omitempty"`
}

type ProductListResponse struct {
	ProductList []Product `json:"product_list"`
	Total       int       `json:"total"`
}
