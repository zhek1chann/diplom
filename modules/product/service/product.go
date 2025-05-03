package service

import (
	"context"
	"diploma/modules/product/model"
)

func (s *ProductService) Product(ctx context.Context, query *model.ProductQuery) (*model.DetailedProduct, error) {

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

func (s *ProductService) ProductPriceBySupplier(ctx context.Context, productID, supplierID int64) (int, error) {
	return s.productRepository.GetProductPriceBySupplier(ctx, productID, supplierID)
}

func (s *ProductService) ProductInfo(ctx context.Context, id int64) (*model.Product, error) {
	return s.productRepository.GetProduct(ctx, id)
}
