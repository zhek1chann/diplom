package service

import (
	"context"
	"diploma/modules/cart/model"
	"diploma/pkg/client/db"
)

type cartServ struct {
	repo      ICartRepository
	txManager db.TxManager
}

func NewService(
	cartRepository ICartRepository,
	txManager db.TxManager,
) *cartServ {
	return &cartServ{
		repo:      cartRepository,
		txManager: txManager,
	}
}

type ICartRepository interface {
	GetCart(ctx context.Context, userID int64) (*model.Cart, error)
}
