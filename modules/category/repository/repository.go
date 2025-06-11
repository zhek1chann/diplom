package repository

import (
	"context"
	"fmt"
	"time"

	"diploma/modules/category/model"
	"diploma/pkg/client/db"
)

type Repository interface {
	CreateCategory(ctx context.Context, req model.CreateCategoryRequest) (*model.Category, error)
	UpdateCategory(ctx context.Context, id int, req model.UpdateCategoryRequest) (*model.Category, error)
	DeleteCategory(ctx context.Context, id int) error
	GetCategory(ctx context.Context, id int) (*model.Category, error)
	ListCategories(ctx context.Context) ([]model.Category, error)
	GetCategoriesTree(ctx context.Context) ([]model.CategoryTree, error)

	CreateSubcategory(ctx context.Context, req model.CreateSubcategoryRequest) (*model.Subcategory, error)
	UpdateSubcategory(ctx context.Context, id int, req model.UpdateSubcategoryRequest) (*model.Subcategory, error)
	DeleteSubcategory(ctx context.Context, id int) error
	GetSubcategory(ctx context.Context, id int) (*model.Subcategory, error)
	ListSubcategories(ctx context.Context, categoryID int) ([]model.Subcategory, error)
}

type repository struct {
	db db.Client
}

func NewRepository(db db.Client) Repository {
	return &repository{db: db}
}

func (r *repository) CreateCategory(ctx context.Context, req model.CreateCategoryRequest) (*model.Category, error) {
	query := `
		INSERT INTO categories (name, description)
		VALUES ($1, $2)
		RETURNING id, name, description, created_at, updated_at`

	q := db.Query{
		Name:     "category_repository.CreateCategory",
		QueryRaw: query,
	}

	var category model.Category
	err := r.db.DB().ScanOneContext(ctx, &category, q, req.Name, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return &category, nil
}

func (r *repository) UpdateCategory(ctx context.Context, id int, req model.UpdateCategoryRequest) (*model.Category, error) {
	query := `
		UPDATE categories
		SET name = COALESCE($1, name),
			description = COALESCE($2, description),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, name, description, created_at, updated_at`

	q := db.Query{
		Name:     "category_repository.UpdateCategory",
		QueryRaw: query,
	}

	var category model.Category
	err := r.db.DB().ScanOneContext(ctx, &category, q, req.Name, req.Description, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return &category, nil
}

func (r *repository) DeleteCategory(ctx context.Context, id int) error {
	query := `DELETE FROM categories WHERE id = $1`
	q := db.Query{
		Name:     "category_repository.DeleteCategory",
		QueryRaw: query,
	}

	result, err := r.db.DB().ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	rows := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

func (r *repository) GetCategory(ctx context.Context, id int) (*model.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		WHERE id = $1`

	q := db.Query{
		Name:     "category_repository.GetCategory",
		QueryRaw: query,
	}

	var category model.Category
	err := r.db.DB().ScanOneContext(ctx, &category, q, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &category, nil
}

func (r *repository) ListCategories(ctx context.Context) ([]model.Category, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM categories
		ORDER BY name`

	q := db.Query{
		Name:     "category_repository.ListCategories",
		QueryRaw: query,
	}

	var categories []model.Category
	err := r.db.DB().ScanAllContext(ctx, &categories, q)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return categories, nil
}

func (r *repository) CreateSubcategory(ctx context.Context, req model.CreateSubcategoryRequest) (*model.Subcategory, error) {
	query := `
		INSERT INTO subcategories (category_id, name, description)
		VALUES ($1, $2, $3)
		RETURNING id, category_id, name, description, created_at, updated_at`

	q := db.Query{
		Name:     "category_repository.CreateSubcategory",
		QueryRaw: query,
	}

	var subcategory model.Subcategory
	err := r.db.DB().ScanOneContext(ctx, &subcategory, q, req.CategoryID, req.Name, req.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create subcategory: %w", err)
	}

	return &subcategory, nil
}

func (r *repository) UpdateSubcategory(ctx context.Context, id int, req model.UpdateSubcategoryRequest) (*model.Subcategory, error) {
	query := `
		UPDATE subcategories
		SET name = COALESCE($1, name),
			description = COALESCE($2, description),
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
		RETURNING id, category_id, name, description, created_at, updated_at`

	q := db.Query{
		Name:     "category_repository.UpdateSubcategory",
		QueryRaw: query,
	}

	var subcategory model.Subcategory
	err := r.db.DB().ScanOneContext(ctx, &subcategory, q, req.Name, req.Description, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update subcategory: %w", err)
	}

	return &subcategory, nil
}

func (r *repository) DeleteSubcategory(ctx context.Context, id int) error {
	query := `DELETE FROM subcategories WHERE id = $1`
	q := db.Query{
		Name:     "category_repository.DeleteSubcategory",
		QueryRaw: query,
	}

	result, err := r.db.DB().ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("failed to delete subcategory: %w", err)
	}

	rows := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("subcategory not found")
	}

	return nil
}

func (r *repository) GetSubcategory(ctx context.Context, id int) (*model.Subcategory, error) {
	query := `
		SELECT id, category_id, name, description, created_at, updated_at
		FROM subcategories
		WHERE id = $1`

	q := db.Query{
		Name:     "category_repository.GetSubcategory",
		QueryRaw: query,
	}

	var subcategory model.Subcategory
	err := r.db.DB().ScanOneContext(ctx, &subcategory, q, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get subcategory: %w", err)
	}

	return &subcategory, nil
}

func (r *repository) ListSubcategories(ctx context.Context, categoryID int) ([]model.Subcategory, error) {
	query := `
		SELECT id, category_id, name, description, created_at, updated_at
		FROM subcategories
		WHERE category_id = $1
		ORDER BY name`

	q := db.Query{
		Name:     "category_repository.ListSubcategories",
		QueryRaw: query,
	}

	var subcategories []model.Subcategory
	err := r.db.DB().ScanAllContext(ctx, &subcategories, q, categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to list subcategories: %w", err)
	}

	return subcategories, nil
}

func (r *repository) GetCategoriesTree(ctx context.Context) ([]model.CategoryTree, error) {
	// Query to get all categories with their subcategories using LEFT JOIN
	query := `
		SELECT 
			c.id as category_id,
			c.name as category_name,
			c.description as category_description,
			c.created_at as category_created_at,
			c.updated_at as category_updated_at,
			COALESCE(s.id, 0) as subcategory_id,
			COALESCE(s.category_id, 0) as subcategory_category_id,
			COALESCE(s.name, '') as subcategory_name,
			COALESCE(s.description, '') as subcategory_description,
			COALESCE(s.created_at, c.created_at) as subcategory_created_at,
			COALESCE(s.updated_at, c.updated_at) as subcategory_updated_at
		FROM categories c
		LEFT JOIN subcategories s ON c.id = s.category_id
		ORDER BY c.name, s.name`

	q := db.Query{
		Name:     "category_repository.GetCategoriesTree",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories tree: %w", err)
	}
	defer rows.Close()

	// Map to store categories by their ID
	categoryMap := make(map[int]*model.CategoryTree)

	for rows.Next() {
		var categoryID int
		var categoryName, categoryDescription string
		var categoryCreatedAt, categoryUpdatedAt time.Time
		var subcategoryID, subcategoryCategoryID int
		var subcategoryName, subcategoryDescription string
		var subcategoryCreatedAt, subcategoryUpdatedAt time.Time

		err := rows.Scan(
			&categoryID,
			&categoryName,
			&categoryDescription,
			&categoryCreatedAt,
			&categoryUpdatedAt,
			&subcategoryID,
			&subcategoryCategoryID,
			&subcategoryName,
			&subcategoryDescription,
			&subcategoryCreatedAt,
			&subcategoryUpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Get or create category in map
		category, exists := categoryMap[categoryID]
		if !exists {
			category = &model.CategoryTree{
				ID:            categoryID,
				Name:          categoryName,
				Description:   categoryDescription,
				CreatedAt:     categoryCreatedAt,
				UpdatedAt:     categoryUpdatedAt,
				Subcategories: []model.Subcategory{},
			}
			categoryMap[categoryID] = category
		}

		// Add subcategory if it exists (subcategoryID > 0 means it's not NULL)
		if subcategoryID > 0 {
			subcategory := model.Subcategory{
				ID:          subcategoryID,
				CategoryID:  subcategoryCategoryID,
				Name:        subcategoryName,
				Description: subcategoryDescription,
				CreatedAt:   subcategoryCreatedAt,
				UpdatedAt:   subcategoryUpdatedAt,
			}
			category.Subcategories = append(category.Subcategories, subcategory)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	// Convert map to slice
	var categoryTrees []model.CategoryTree
	for _, category := range categoryMap {
		categoryTrees = append(categoryTrees, *category)
	}

	return categoryTrees, nil
}
