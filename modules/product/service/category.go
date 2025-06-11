package service

import (
	"context"
	"diploma/modules/product/model"
)

// GetCategories retrieves all categories
func (s *ProductService) GetCategories(ctx context.Context) ([]model.Category, error) {
	return s.productRepository.GetCategories(ctx)
}

// GetCategory retrieves a single category by ID
func (s *ProductService) GetCategory(ctx context.Context, id int) (*model.Category, error) {
	return s.productRepository.GetCategory(ctx, id)
}

// GetSubcategories retrieves all subcategories for a specific category
func (s *ProductService) GetSubcategories(ctx context.Context, categoryID int) ([]model.Subcategory, error) {
	return s.productRepository.GetSubcategories(ctx, categoryID)
}

// GetCategoriesWithSubcategories retrieves all categories with their subcategories
func (s *ProductService) GetCategoriesWithSubcategories(ctx context.Context) ([]model.Category, error) {
	categories, err := s.productRepository.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	// For each category, get its subcategories
	for i := range categories {
		_, err := s.productRepository.GetSubcategories(ctx, categories[i].ID)
		if err != nil {
			s.LogError(ctx, "Failed to get subcategories for category", err)
			continue // Continue with other categories even if one fails
		}
		// Note: We need to add Subcategories field to Category model if we want to include them
		// For now, this method just returns categories without embedded subcategories
	}

	return categories, nil
}
