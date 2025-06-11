package model

type Category struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description,omitempty"`
	Subcategories []Subcategory `json:"subcategories,omitempty"`
}

type Subcategory struct {
	ID          int    `json:"id"`
	CategoryID  int    `json:"category_id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

type CategoriesResponse struct {
	Categories []Category `json:"categories"`
	Total      int        `json:"total"`
}

type CategoryResponse struct {
	Category *Category `json:"category"`
}

type SubcategoriesResponse struct {
	Subcategories []Subcategory `json:"subcategories"`
	Total         int           `json:"total"`
}
