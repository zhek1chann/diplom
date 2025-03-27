package service

import (
	"context"
	"diploma/modules/cart/model"
	"diploma/pkg/client/db"
)

type cartServ struct {
	cartRepo       ICartRepository
	productService IProductService
	txManager      db.TxManager
}

func NewService(
	cartRepository ICartRepository,
	productService IProductService,
	txManager db.TxManager,
) *cartServ {
	return &cartServ{
		cartRepo:       cartRepository,
		txManager:      txManager,
		productService: productService,
	}
}

type IProductService interface {
	GetProductPriceBySupplier(ctx context.Context, productID, supplierID int64) (int, error)
}

type ICartRepository interface {
	UpdateItemQuantity(ctx context.Context, cartId, productId, supplierId int64, quantity int) error
	ItemQuantity(ctx context.Context, cartId, productId, supplierId int64) (int, error)
	GetCart(ctx context.Context, userID int64) (*model.Cart, error)
	CreateCart(ctx context.Context, userID int64) (int64, error)
	AddItem(ctx context.Context, input *model.PutCartQuery) error
	UpdateCartTotal(ctx context.Context, cartID int64, total int) error
}
