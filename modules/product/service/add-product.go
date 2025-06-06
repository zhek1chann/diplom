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

	// First, try to find the product in our database by GTIN
	gtin, err := strconv.ParseInt(req.GTIN, 10, 64)
	if err != nil {
		s.LogError(ctx, "Failed to parse GTIN", err)

		return fmt.Errorf("invalid GTIN format: %w", err)
	}

	//
	// Try to find the product in our database by GTIN
	product, err := s.productRepository.GetProductByGTIN(ctx, gtin)

	if err == nil {
		s.LogInfo(ctx, "Product already exists in the database",
			zap.Int64("product_id", product.ID),
			zap.Int64("gtin", product.GTIN),
		)

		err := s.productRepository.CreateProductSupplier(ctx, req.SupplierID, product.ID, int(req.Price))
		if err != nil {
			s.LogError(ctx, "Failed to create product supplier", err)
			return fmt.Errorf("failed to create product supplier: %w", err)
		}
	}

	// // If product not found in our database, try to fetch it from NCT
	productParsed, err := s.nctParser.ParseProductByGTIN(req.GTIN)
	if err != nil {
		s.LogError(ctx, "Failed to parse product from NCT", err)
		return fmt.Errorf("failed to parse product from NCT: %w", err)
	}
	if len(productParsed) == 0 {
		s.LogError(ctx, "Failed to parse product from NCT", err)
		return fmt.Errorf("failed to parse product from NCT: %w", err)
	}
	product = &productParsed[0]

	if err != nil {
		s.LogError(ctx, "Failed to fetch product from NCT", err)
		return fmt.Errorf("failed to fetch product from NCT: %w", err)
	}
	s.LogInfo(ctx, "Fetched product from NCT",
		zap.Int64("gtin", product.GTIN),
		zap.String("name", product.Name))
	// Convert NCT product to our domain model

	id, err := s.productRepository.CreateProduct(ctx, product)

	if err != nil {
		s.LogError(ctx, "Failed to convert NCT product", err)
		return fmt.Errorf("failed to convert NCT product: %w", err)
	}
	return s.productRepository.CreateProductSupplier(ctx, req.SupplierID, id, int(req.Price))
}
