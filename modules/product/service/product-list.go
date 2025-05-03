package service

import (
	"context"
	"diploma/modules/product/model"
)

func (s *ProductService) ProductList(ctx context.Context, query *model.ProductListQuery) (*model.ProductList, error) {
	productList, err := s.productRepository.GetProductList(ctx, query)
	if err != nil {
		return nil, err
	}

	total, err := s.productRepository.GetTotalProducts(ctx)

	if err != nil {
		return nil, err
	}
	return &model.ProductList{
		Products: productList,
		Total:    total,
	}, err
}

func (s *ProductService) ProductListByIDList(ctx context.Context, idList []int64) ([]*model.Product, error) {
	productList, err := s.productRepository.GetProductListByIDList(ctx, idList)
	if err != nil {
		return nil, err
	}
	return productList, nil
}
