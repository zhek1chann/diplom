package handler

import (
	"context"
	"diploma/modules/product/model"
)

type CatalogHandler struct {
	service IProductService
}

func NewHandler(service IProductService) *CatalogHandler {
	return &CatalogHandler{service: service}
}

type IProductService interface {
	ProductList(ctx context.Context, query *model.ProductListQuery) (*model.ProductList, error)
	Product(ctx context.Context, query *model.ProductQuery) (*model.DetailedProduct, error)
}
