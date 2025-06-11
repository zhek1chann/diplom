package service

import (
	"context"
	"diploma/modules/product/model"
	"fmt"
	"strconv"

	"go.uber.org/zap"
)

// AddProduct adds a new product to the system
func (s *ProductService) AddProduct(ctx context.Context, req *model.AddProductSupplier) error {
	s.LogInfo(ctx, "Adding new product",
		zap.String("gtin", req.GTIN),
		zap.Float64("price", req.Price),
	)

	// Use transaction for the entire operation
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		// First, try to find the product in our database by GTIN
		gtin, err := strconv.ParseInt(req.GTIN, 10, 64)
		if err != nil {
			s.LogError(ctx, "Failed to parse GTIN", err)
			return fmt.Errorf("invalid GTIN format: %w", err)
		}

		// Try to find the product in our database by GTIN
		product, err := s.productRepository.GetProductByGTIN(ctx, gtin)

		if err == nil {
			s.LogInfo(ctx, "Product already exists in the database",
				zap.Int64("product_id", product.ID),
				zap.Int64("gtin", product.GTIN),
			)

			// Product exists, just create the product-supplier relationship
			err := s.productRepository.CreateProductSupplier(ctx, req.SupplierID, product.ID, int(req.Price))
			if err != nil {
				s.LogError(ctx, "Failed to create product supplier", err)
				return fmt.Errorf("failed to create product supplier: %w", err)
			}
			return nil
		}

		// If product not found in our database, try to fetch it from NCT
		s.LogInfo(ctx, "Product not found in database, fetching from NCT")
		productParsed, err := s.nctParser.ParseProductByGTIN(req.GTIN)
		if err != nil {
			s.LogError(ctx, "Failed to parse product from NCT", err)
			return fmt.Errorf("failed to parse product from NCT: %w", err)
		}
		if len(productParsed) == 0 {
			s.LogError(ctx, "Product not found in NCT", err)
			return fmt.Errorf("product not found in NCT")
		}

		product = &productParsed[0]

		// Handle categories intelligently
		err = s.handleCategories(ctx, product, req)
		if err != nil {
			s.LogError(ctx, "Failed to handle categories", err)
			return fmt.Errorf("failed to handle categories: %w", err)
		}

		s.LogInfo(ctx, "Fetched product from NCT",
			zap.Int64("gtin", product.GTIN),
			zap.String("name", product.Name),
			zap.Any("category_id", product.CategoryID),
			zap.Any("subcategory_id", product.SubcategoryID))

		// Create the product in our database
		id, err := s.productRepository.CreateProduct(ctx, product)
		if err != nil {
			s.LogError(ctx, "Failed to create product", err)
			return fmt.Errorf("failed to create product: %w", err)
		}

		// Create the product-supplier relationship
		err = s.productRepository.CreateProductSupplier(ctx, req.SupplierID, id, int(req.Price))
		if err != nil {
			s.LogError(ctx, "Failed to create product supplier", err)
			return fmt.Errorf("failed to create product supplier: %w", err)
		}

		return nil
	})
}

// handleCategories handles category resolution - either from request or from NCT data
func (s *ProductService) handleCategories(ctx context.Context, product *model.Product, req *model.AddProductSupplier) error {
	// If supplier provided category IDs, use them (override NCT data)
	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
		if req.SubcategoryID != nil {
			product.SubcategoryID = req.SubcategoryID
		}
		s.LogInfo(ctx, "Using supplier-provided categories",
			zap.Any("category_id", product.CategoryID),
			zap.Any("subcategory_id", product.SubcategoryID))
		return nil
	}

	// No category provided by supplier, try to use NCT categories
	if product.CategoryName != "" {
		s.LogInfo(ctx, "Processing NCT category",
			zap.String("category_name", product.CategoryName),
			zap.String("subcategory_name", product.SubcategoryName))

		// Find or create category
		categoryID, err := s.productRepository.FindOrCreateCategory(ctx, product.CategoryName)
		if err != nil {
			return fmt.Errorf("failed to find or create category '%s': %w", product.CategoryName, err)
		}
		product.CategoryID = &categoryID

		// Handle subcategory if provided
		if product.SubcategoryName != "" {
			subcategoryID, err := s.productRepository.FindOrCreateSubcategory(ctx, product.SubcategoryName, categoryID)
			if err != nil {
				return fmt.Errorf("failed to find or create subcategory '%s': %w", product.SubcategoryName, err)
			}
			product.SubcategoryID = &subcategoryID
		}

		s.LogInfo(ctx, "Successfully resolved NCT categories",
			zap.Int("category_id", categoryID),
			zap.Any("subcategory_id", product.SubcategoryID))
	} else {
		s.LogInfo(ctx, "No category information available from NCT or supplier")
	}

	return nil
}
