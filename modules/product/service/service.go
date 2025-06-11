package service

import (
	"context"
	"diploma/modules/product/client/nct"
	"diploma/modules/product/model"
	"diploma/pkg/client/db"
	"diploma/pkg/service"
)

type ProductService struct {
	service.BaseService
	productRepository IProductRepository
	txManager         db.TxManager
	nctParser         *nct.NCTParser
}

func NewService(
	productRepository IProductRepository,
	txManager db.TxManager,
	nctParser *nct.NCTParser,
) *ProductService {
	return &ProductService{
		BaseService:       service.NewBaseService("product"),
		productRepository: productRepository,
		txManager:         txManager,
		nctParser:         nctParser,
	}
}

type IProductRepository interface {
	GetProduct(ctx context.Context, id int64) (*model.Product, error)
	GetProductByGTIN(ctx context.Context, gtin int64) (*model.Product, error)
	GetSupplierProductListByProduct(ctx context.Context, id int64) ([]model.ProductSupplier, error)
	GetProductListByIDList(ctx context.Context, idList []int64) ([]*model.Product, error)
	GetProductList(ctx context.Context, query *model.ProductListQuery) ([]model.Product, error)
	GetTotalProducts(ctx context.Context) (int, error)
	GetProductPriceBySupplier(ctx context.Context, productID, supplierID int64) (int, error)
	CreateProduct(ctx context.Context, product *model.Product) (int64, error)

	CreateProductSupplier(ctx context.Context, supplierID, productID int64, price int) error
	GetProductListBySupplier(ctx context.Context, supplierID int64, limit, offset int) ([]model.Product, error)
	GetTotalProductsBySupplier(ctx context.Context, supplierID int64) (int, error)

	// Category management
	FindOrCreateCategory(ctx context.Context, name string) (int, error)
	FindOrCreateSubcategory(ctx context.Context, name string, categoryID int) (int, error)
}
