package service

import (
	"context"
	"diploma/modules/cart/model"
	"errors"
)

type ICardService interface {
	GetCart(ctx context.Context, userID int64) (*model.Cart, error)
	AddProductToCard(ctx context.Context, input *model.PutCartQuery) error
	DeleteProductFromCart(ctx context.Context, input *model.DeleteProductQuery) error
}

func (s *cartServ) AddProductToCard(ctx context.Context, input *model.PutCartQuery) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		cart, errTx := s.cartRepo.GetCart(ctx, input.CustomerID)
		if errTx != nil {
			if errors.Is(errTx, model.ErrNoRows) {
				id, errTx := s.cartRepo.CreateCart(ctx, input.CustomerID)
				if errTx != nil {
					return errTx
				}
				cart = &model.Cart{
					ID:    id,
					Total: 0}
			} else {
				return errTx

			}
			input.CartID = cart.CustomerID
			errTx = s.cartRepo.AddItem(ctx, input)
			if errTx != nil {
				return errTx
			}

			cart.Total += input.Quantity * input.Price
			errTx = s.cartRepo.UpdateCartTotal(ctx, cart.ID, cart.Total)
			if errTx != nil {
				return errTx
			}
			return nil
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *cartServ) GetCart(ctx context.Context, userID int64) (*model.Cart, error) {
	var err error
	var cart *model.Cart
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		cart, errTx = s.cartRepo.GetCart(ctx, userID)
		if errTx != nil {
			return errTx
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return cart, nil
}
