package service

import (
	"context"
	"diploma/modules/product/model"
	"diploma/pkg/client/db"
	"diploma/pkg/service"
)

type ProductService struct {
	service.BaseService
	productRepository IProductRepository
	txManager         db.TxManager
}

func NewService(
	productRepository IProductRepository,
	txManager db.TxManager,
) *ProductService {
	return &ProductService{
		BaseService:       service.NewBaseService("product"),
		productRepository: productRepository,
		txManager:         txManager,
	}
}

type IProductRepository interface {
	GetProduct(ctx context.Context, id int64) (*model.Product, error)
	GetSupplierProductListByProduct(ctx context.Context, id int64) ([]model.ProductSupplier, error)
	GetProductListByIDList(ctx context.Context, idList []int64) ([]*model.Product, error)
	GetProductList(ctx context.Context, query *model.ProductListQuery) ([]model.Product, error)
	GetTotalProducts(ctx context.Context) (int, error)
	GetProductPriceBySupplier(ctx context.Context, productID, supplierID int64) (int, error)
}
