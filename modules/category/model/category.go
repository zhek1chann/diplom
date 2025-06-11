package model

import (
	"time"
)

type Category struct {
	ID          int       `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

type Subcategory struct {
	ID          int       `json:"id" db:"id"`
	CategoryID  int       `json:"category_id" db:"category_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// CategoryTree represents a category with its nested subcategories
type CategoryTree struct {
	ID            int           `json:"id"`
	Name          string        `json:"name"`
	Description   string        `json:"description,omitempty"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	Subcategories []Subcategory `json:"subcategories"`
}

// CategoriesTreeResponse represents the response for the categories tree endpoint
type CategoriesTreeResponse struct {
	Categories []CategoryTree `json:"categories"`
	Total      int            `json:"total"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description,omitempty"`
}

type CreateSubcategoryRequest struct {
	CategoryID  int    `json:"category_id" validate:"required"`
	Name        string `json:"name" validate:"required,min=2,max=100"`
	Description string `json:"description,omitempty"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description string `json:"description,omitempty"`
}

type UpdateSubcategoryRequest struct {
	Name        string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description string `json:"description,omitempty"`
}
