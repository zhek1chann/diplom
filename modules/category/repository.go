package category

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateCategory(ctx context.Context, req CreateCategoryRequest) (*Category, error)
	UpdateCategory(ctx context.Context, id int, req UpdateCategoryRequest) (*Category, error)
	DeleteCategory(ctx context.Context, id int) error
	GetCategory(ctx context.Context, id int) (*Category, error)
	ListCategories(ctx context.Context) ([]Category, error)

	CreateSubcategory(ctx context.Context, req CreateSubcategoryRequest) (*Subcategory, error)
	UpdateSubcategory(ctx context.Context, id int, req UpdateSubcategoryRequest) (*Subcategory, error)
	DeleteSubcategory(ctx context.Context, id int) error
	GetSubcategory(ctx context.Context, id int) (*Subcategory, error)
	ListSubcategories(ctx context.Context, categoryID int) ([]Subcategory, error)
}

type repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateCategory(ctx context.Context, req CreateCategoryRequest) (*Category, error) {
	query := `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at`

	var category Category
	err := r.db.GetContext(ctx, &category, query, req.Name, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return &category, nil
}

func (r *repository) UpdateCategory(ctx context.Context, id int, req UpdateCategoryRequest) (*Category, error) {
	query := `
		UPDATE categories
		SET name = COALESCE($1, name),
			description = COALESCE($2, description),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, name, description, created_at, updated_at`

	var category Category
	err := r.db.GetContext(ctx, &category, query, req.Name, req.Description, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return &category, nil
}

func (r *repository) DeleteCategory(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

func (r *repository) GetCategory(ctx context.Context, id int) (*Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		WHERE id = $1`

	var category Category
	err := r.db.GetContext(ctx, &category, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &category, nil
}

func (r *repository) ListCategories(ctx context.Context) ([]Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		ORDER BY name`

	var categories []Category
	err := r.db.SelectContext(ctx, &categories, query)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return categories, nil
}

func (r *repository) CreateSubcategory(ctx context.Context, req CreateSubcategoryRequest) (*Subcategory, error) {
	query := `
		INSERT INTO subcategories (category_id, name, description)
		VALUES ($1, $2, $3)
		RETURNING id, category_id, name, description, created_at, updated_at`

	var subcategory Subcategory
	err := r.db.GetContext(ctx, &subcategory, query, req.CategoryID, req.Name, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create subcategory: %w", err)
	}

	return &subcategory, nil
}

func (r *repository) UpdateSubcategory(ctx context.Context, id int, req UpdateSubcategoryRequest) (*Subcategory, error) {
	query := `
		UPDATE subcategories
		SET name = COALESCE($1, name),
			description = COALESCE($2, description),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, category_id, name, description, created_at, updated_at`

	var subcategory Subcategory
	err := r.db.GetContext(ctx, &subcategory, query, req.Name, req.Description, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update subcategory: %w", err)
	}

	return &subcategory, nil
}

func (r *repository) DeleteSubcategory(ctx context.Context, id int) error {
	query := `DELETE FROM subcategories WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subcategory: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("subcategory not found")
	}

	return nil
}

func (r *repository) GetSubcategory(ctx context.Context, id int) (*Subcategory, error) {
	query := `
		SELECT id, category_id, name, description, created_at, updated_at
		FROM subcategories
		WHERE id = $1`

	var subcategory Subcategory
	err := r.db.GetContext(ctx, &subcategory, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subcategory: %w", err)
	}

	return &subcategory, nil
}

func (r *repository) ListSubcategories(ctx context.Context, categoryID int) ([]Subcategory, error) {
	query := `
		SELECT id, category_id, name, description, created_at, updated_at
		FROM subcategories
		WHERE category_id = $1
		ORDER BY name`

	var subcategories []Subcategory
	err := r.db.SelectContext(ctx, &subcategories, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to list subcategories: %w", err)
	}

	return subcategories, nil
}
