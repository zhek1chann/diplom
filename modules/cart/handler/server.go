package handler

import (
	"context"
	"diploma/modules/cart/model"
)

type CardHandler struct {
	service ICardService
}

func NewHandler(service ICardService) *CardHandler {
	return &CardHandler{service: service}
}

type ICardService interface {
	GetCart(ctx context.Context, userID int64) (*model.Cart, error)
	AddProductToCard(ctx context.Context, input *model.PutCartQuery) error
	DeleteProductFromCart(ctx context.Context, input *model.DeleteProductQuery) error
}
