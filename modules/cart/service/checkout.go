package service

import (
	"context"
	"diploma/modules/cart/model"
)

func (s *cartServ) Checkout(ctx context.Context, userID int64) (bool, error) {

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		cart, errTx := s.Cart(ctx, userID)
		if errTx != nil {
			return errTx
		}
		cart.CustomerID = userID
		if !checkCartForCheckout(cart) {
			return model.ErrInvalidCart
		}

		errTx = s.orderService.CreateOrder(ctx, cart)
		if errTx != nil {
			return errTx
		}
		errTx = s.cartRepo.DeleteCart(ctx, cart.ID)
		if errTx != nil {
			return errTx
		}

		return nil

	})

	if err != nil {
		return false, err
	}
	return false, nil
}

func checkCartForCheckout(cart *model.Cart) bool {
	for _, supplier := range cart.Suppliers {
		sum := 0
		for _, product := range supplier.ProductList {
			if product.Quantity <= 0 {
				return false
			} else if product.Price <= 0 {
				return false
			}
			sum += product.Price * product.Quantity
		}
		if sum < supplier.OrderAmount {
			return false
		}
	}
	return true
}
