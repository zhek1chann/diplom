package service

import (
	"context"
	"diploma/modules/category/model"
	"diploma/modules/category/repository"
)

type Service interface {
	CreateCategory(ctx context.Context, req model.CreateCategoryRequest) (*model.Category, error)
	UpdateCategory(ctx context.Context, id int, req model.UpdateCategoryRequest) (*model.Category, error)
	DeleteCategory(ctx context.Context, id int) error
	GetCategory(ctx context.Context, id int) (*model.Category, error)
	ListCategories(ctx context.Context) ([]model.Category, error)
	GetCategoriesTree(ctx context.Context) (*model.CategoriesTreeResponse, error)

	CreateSubcategory(ctx context.Context, req model.CreateSubcategoryRequest) (*model.Subcategory, error)
	UpdateSubcategory(ctx context.Context, id int, req model.UpdateSubcategoryRequest) (*model.Subcategory, error)
	DeleteSubcategory(ctx context.Context, id int) error
	GetSubcategory(ctx context.Context, id int) (*model.Subcategory, error)
	ListSubcategories(ctx context.Context, categoryID int) ([]model.Subcategory, error)
}

type service struct {
	repo repository.Repository
}

func NewService(repo repository.Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateCategory(ctx context.Context, req model.CreateCategoryRequest) (*model.Category, error) {
	return s.repo.CreateCategory(ctx, req)
}

func (s *service) UpdateCategory(ctx context.Context, id int, req model.UpdateCategoryRequest) (*model.Category, error) {
	return s.repo.UpdateCategory(ctx, id, req)
}

func (s *service) DeleteCategory(ctx context.Context, id int) error {
	return s.repo.DeleteCategory(ctx, id)
}

func (s *service) GetCategory(ctx context.Context, id int) (*model.Category, error) {
	return s.repo.GetCategory(ctx, id)
}

func (s *service) ListCategories(ctx context.Context) ([]model.Category, error) {
	return s.repo.ListCategories(ctx)
}

func (s *service) CreateSubcategory(ctx context.Context, req model.CreateSubcategoryRequest) (*model.Subcategory, error) {
	// Verify that the category exists
	_, err := s.GetCategory(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	return s.repo.CreateSubcategory(ctx, req)
}

func (s *service) UpdateSubcategory(ctx context.Context, id int, req model.UpdateSubcategoryRequest) (*model.Subcategory, error) {
	return s.repo.UpdateSubcategory(ctx, id, req)
}

func (s *service) DeleteSubcategory(ctx context.Context, id int) error {
	return s.repo.DeleteSubcategory(ctx, id)
}

func (s *service) GetSubcategory(ctx context.Context, id int) (*model.Subcategory, error) {
	return s.repo.GetSubcategory(ctx, id)
}

func (s *service) ListSubcategories(ctx context.Context, categoryID int) ([]model.Subcategory, error) {
	return s.repo.ListSubcategories(ctx, categoryID)
}

func (s *service) GetCategoriesTree(ctx context.Context) (*model.CategoriesTreeResponse, error) {
	categories, err := s.repo.GetCategoriesTree(ctx)
	if err != nil {
		return nil, err
	}

	return &model.CategoriesTreeResponse{
		Categories: categories,
		Total:      len(categories),
	}, nil
}
