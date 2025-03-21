package service

import (
	"context"
	"diploma/modules/product/model"
)

func (s *productServ) Product(ctx context.Context, query *model.ProductQuery) (*model.DetailedProduct, error) {

	product, err := s.productRepository.GetProduct(ctx, query.ID)
	if err != nil {
		return nil, err
	}

	productSupplierList, err := s.productRepository.GetSupplierProductListByProduct(ctx, query.ID)

	if err != nil {
		return nil, err
	}

	return &model.DetailedProduct{
		Product:             product,
		ProductSupplierList: productSupplierList,
	}, err

}
