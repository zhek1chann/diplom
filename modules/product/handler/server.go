package handler

import (
	"context"
	"diploma/modules/product/model"
)

type CatalogHandler struct {
	service ICatalogService
}

func NewHandler(service ICatalogService) *CatalogHandler {
	return &CatalogHandler{service: service}
}

type ICatalogService interface {
	PageCount(ctx context.Context, query *model.PageCountQuery) (int, error)
	ProductList(ctx context.Context, query *model.ProductListQuery) (*model.ProductList, error)
	Product(ctx context.Context, query *model.ProductQuery) (*model.DetailedProduct, error)
}
