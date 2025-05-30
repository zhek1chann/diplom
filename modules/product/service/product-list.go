package service

import (
	"context"
	"diploma/modules/product/model"

	"go.uber.org/zap"
)

func (s *ProductService) ProductList(ctx context.Context, query *model.ProductListQuery) (*model.ProductList, error) {
	s.LogInfo(ctx, "Fetching product list",
		zap.Int("offset", query.Offset),
		zap.Int("limit", query.Limit),
	)

	productList, err := s.productRepository.GetProductList(ctx, query)
	if err != nil {
		s.LogError(ctx, "Failed to get product list", err)
		return nil, err
	}

	total, err := s.productRepository.GetTotalProducts(ctx)
	if err != nil {
		s.LogError(ctx, "Failed to get total products count", err)
		return nil, err
	}

	result := &model.ProductList{
		Products: productList,
		Total:    total,
	}

	s.LogInfo(ctx, "Successfully fetched product list",
		zap.Int("total_products", total),
		zap.Int("returned_products", len(productList)),
	)

	return result, nil
}

func (s *ProductService) ProductListByIDList(ctx context.Context, idList []int64) ([]*model.Product, error) {
	s.LogInfo(ctx, "Fetching products by ID list",
		zap.Int("product_count", len(idList)),
	)

	productList, err := s.productRepository.GetProductListByIDList(ctx, idList)
	if err != nil {
		s.LogError(ctx, "Failed to get products by ID list", err,
			zap.Any("product_ids", idList),
		)
		return nil, err
	}

	s.LogInfo(ctx, "Successfully fetched products by ID list",
		zap.Int("requested_count", len(idList)),
		zap.Int("found_count", len(productList)),
	)

	return productList, nil
}
