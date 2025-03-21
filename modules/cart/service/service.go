package service

import (
	"context"
	"diploma/modules/cart/model"
	"diploma/pkg/client/db"
)

type cartServ struct {
	cartRepo  ICartRepository
	txManager db.TxManager
}

func NewService(
	cartRepository ICartRepository,
	txManager db.TxManager,
) *cartServ {
	return &cartServ{
		cartRepo:  cartRepository,
		txManager: txManager,
	}
}

type ICartRepository interface {
	GetCart(ctx context.Context, userID int64) (*model.Cart, error)
	CreateCart(ctx context.Context, userID int64) (int64, error)
	AddItem(ctx context.Context, input *model.PutCartQuery) error
	UpdateCartTotal(ctx context.Context, cartID int64, total int) error
}
