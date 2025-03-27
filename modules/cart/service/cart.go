package service

import (
	"context"
	"diploma/modules/cart/model"
	"errors"
)

func (s *cartServ) AddProductToCard(ctx context.Context, query *model.PutCartQuery) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		cart, errTx := s.cartRepo.GetCart(ctx, query.CustomerID)

		if errTx != nil {
			if errors.Is(errTx, model.ErrNoRows) {
				id, errTx := s.cartRepo.CreateCart(ctx, query.CustomerID)
				if errTx != nil {
					return errTx
				}
				cart = &model.Cart{
					ID:    id,
					Total: 0,
				}
			} else {
				return errTx

			}
			query.CartID = cart.CustomerID
		}
		query.Price, errTx = s.productService.GetProductPriceBySupplier(ctx, query.ProductID, query.SupplierID)
		if errTx != nil {
			return errTx
		}

		itemQuantity, errTx := s.cartRepo.ItemQuantity(ctx, query.CartID, query.ProductID, query.SupplierID)

		if errTx != nil {
			if errors.Is(errTx, model.ErrNoRows) {
				errTx = s.cartRepo.AddItem(ctx, query)
				if errTx != nil {
					return errTx
				}
			} else {
				return errTx
			}
		} else {

			errTx = s.cartRepo.UpdateItemQuantity(ctx, query.CartID, query.ProductID, query.SupplierID, itemQuantity+query.Quantity)
			if errTx != nil {
				return errTx
			}
		}

		cart.Total += query.Price * query.Quantity
		// fmt.Println(cart)
		// errTx = s.cartRepo.UpdateCartTotal(ctx, cart.ID, cart.Total)
		// if errTx != nil {
		// 	return errTx
		// }
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
