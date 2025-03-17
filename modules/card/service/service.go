package service

import (
	"context"
	"diploma/modules/product/model"
	"diploma/pkg/client/db"
)

type productServ struct {
	productRepository IProductRepository
	txManager         db.TxManager
}

func NewService(
	productRepository IProductRepository,
	txManager db.TxManager,
) *productServ {
	return &productServ{
		productRepository: productRepository,
		txManager:         txManager,
	}
}

type IProductRepository interface {
	GetProduct(ctx context.Context, id int64) (*model.Product, error)
	GetSupplierInfoListByProduct(ctx context.Context, id int64) ([]model.ProductSupplierInfo, error)
	GetProductList(ctx context.Context, query *model.ProductListQuery) ([]model.Product, error)
	GetTotalProducts(ctx context.Context) (int, error)
}
