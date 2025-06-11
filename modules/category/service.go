package category

import (
	"context"
)

type Service interface {
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

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateCategory(ctx context.Context, req CreateCategoryRequest) (*Category, error) {
	return s.repo.CreateCategory(ctx, req)
}

func (s *service) UpdateCategory(ctx context.Context, id int, req UpdateCategoryRequest) (*Category, error) {
	return s.repo.UpdateCategory(ctx, id, req)
}

func (s *service) DeleteCategory(ctx context.Context, id int) error {
	return s.repo.DeleteCategory(ctx, id)
}

func (s *service) GetCategory(ctx context.Context, id int) (*Category, error) {
	return s.repo.GetCategory(ctx, id)
}

func (s *service) ListCategories(ctx context.Context) ([]Category, error) {
	return s.repo.ListCategories(ctx)
}

func (s *service) CreateSubcategory(ctx context.Context, req CreateSubcategoryRequest) (*Subcategory, error) {
	// Verify that the category exists
	_, err := s.GetCategory(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	return s.repo.CreateSubcategory(ctx, req)
}

func (s *service) UpdateSubcategory(ctx context.Context, id int, req UpdateSubcategoryRequest) (*Subcategory, error) {
	return s.repo.UpdateSubcategory(ctx, id, req)
}

func (s *service) DeleteSubcategory(ctx context.Context, id int) error {
	return s.repo.DeleteSubcategory(ctx, id)
}

func (s *service) GetSubcategory(ctx context.Context, id int) (*Subcategory, error) {
	return s.repo.GetSubcategory(ctx, id)
}

func (s *service) ListSubcategories(ctx context.Context, categoryID int) ([]Subcategory, error) {
	return s.repo.ListSubcategories(ctx, categoryID)
} 