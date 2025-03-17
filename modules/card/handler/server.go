package handler

import (
	"context"
	"diploma/modules/card/model"
)

type CardHandler struct {
	service ICardService
}

func NewHandler(service ICardService) *CardHandler {
	return &CardHandler{service: service}
}

type ICardService interface {
	PutToCard(ctx context.Context, input *model.CardInput) error
}
