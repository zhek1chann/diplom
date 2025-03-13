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
	PageCount(ctx context.Context, query *model.PageCountQuery) (*model.PageCount, error)
	ProductList(ctx context.Context, query *model.ProductListQuery) (*model.ProductList, error)
	Product(ctx context.Context, query *model.ProductQuery) (*model.DetailedProduct, error)
}
