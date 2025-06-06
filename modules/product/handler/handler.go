package handler

import (
	"context"
	productModel "diploma/modules/product/model"
)

type IProductService interface {
	Product(ctx context.Context, query *productModel.ProductQuery) (*productModel.DetailedProduct, error)
	ProductList(ctx context.Context, query *productModel.ProductListQuery) (*productModel.ProductList, error)
	AddProduct(ctx context.Context, req *productModel.AddProductSupplier) error
	GetProductListBySupplier(ctx context.Context, supplierID int64, limit, offset int) (*productModel.ProductList, error)
}

type CatalogHandler struct {
	service IProductService
}

func NewHandler(service IProductService) *CatalogHandler {
	return &CatalogHandler{
		service: service,
	}
}
