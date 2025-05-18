package handler

import (
	"context"
	"diploma/modules/cart/model"
)

type CartHandler struct {
	service ICartService
}

func NewHandler(service ICartService) *CartHandler {
	return &CartHandler{service: service}
}

type ICartService interface {
	Checkout(ctx context.Context, userID int64) (bool, error)
	Cart(ctx context.Context, userID int64) (*model.Cart, error)
	AddProductToCard(ctx context.Context, input *model.PutCartQuery) error
	DeleteProductFromCart(ctx context.Context, input *model.PutCartQuery) error
	// DeleteProductFromCart(ctx context.Context, input *model.DeleteProductQuery) error
}
