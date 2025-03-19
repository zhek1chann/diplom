package service

import (
	"context"
	"diploma/modules/cart/model"
)

type ICardService interface {
	GetCart(ctx context.Context, userID int64) (*model.Cart, error)
	AddProductToCard(ctx context.Context, input *model.PutCartQuery) error
	DeleteProductFromCart(ctx context.Context, input *model.DeleteProductQuery) error
}

func (s *cartServ) GetCart(ctx context.Context, userID int64) (*model.Cart, error) {
	return s.repo.GetCart(ctx, userID)
}
