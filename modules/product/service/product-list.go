package service

import (
	"context"
	"diploma/modules/product/model"
)

func (s *productServ) ProductList(ctx context.Context, query *model.ProductListQuery) (*model.ProductList, error) {
	productList, err := s.productRepository.GetProductList(ctx, query)
	if err != nil {
		return nil, err
	}

	return &model.ProductList{
		Products: productList,
	}, err
}

func (s *productServ) PageCount(ctx context.Context, query *model.PageCountQuery) (*model.PageCount, error) {
	total, err := s.productRepository.GetTotalProducts(ctx)

	if err != nil {
		return nil, err
	}
	pages := total / query.PageSize
	return &model.PageCount{Pages: pages}, nil
}
