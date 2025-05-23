package service

import (
	"context"
	"diploma/modules/cart/model"
	"errors"
)

func (s *cartServ) AddProductToCard(ctx context.Context, query *model.PutCartQuery) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		cart, errTx := s.cartRepo.Cart(ctx, query.CustomerID)

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
		}
		query.CartID = cart.ID

		query.Price, errTx = s.productService.ProductPriceBySupplier(ctx, query.ProductID, query.SupplierID)
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
		errTx = s.cartRepo.UpdateCartTotal(ctx, cart.ID, cart.Total)
		if errTx != nil {
			return errTx
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (s *cartServ) Cart(ctx context.Context, userID int64) (*model.Cart, error) {
	var err error
	var cart *model.Cart
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		cart, errTx = s.cartRepo.Cart(ctx, userID)
		if errTx != nil {
			if errors.Is(errTx, model.ErrNoRows) {
				cart = &model.Cart{
					ID:         0,
					CustomerID: userID,
					Total:      0,
					Suppliers:  []model.Supplier{},
				}
				return nil
			}
			return errTx
		}

		cart.Suppliers, errTx = s.cartRepo.GetCartItems(ctx, cart.ID)
		if errTx != nil {
			if errors.Is(errTx, model.ErrNoRows) {
				return model.ErrNoRows
			}
			return errTx
		}

		supplierIdList := make([]int64, 0, len(cart.Suppliers))
		for _, supplier := range cart.Suppliers {
			supplierIdList = append(supplierIdList, supplier.ID)
		}

		suppliers, errTx := s.supplierService.SupplierListByIDList(ctx, supplierIdList)
		if errTx != nil {
			return errTx
		}

		for i, supplier := range suppliers {
			cart.Suppliers[i].Name = supplier.Name
			cart.Suppliers[i].OrderAmount = supplier.OrderAmount
			cart.Suppliers[i].FreeDeliveryAmount = supplier.FreeDeliveryAmount
			cart.Suppliers[i].DeliveryFee = supplier.DeliveryFee
			cart.Suppliers[i].TotalAmount = getTotalSupplier(ctx, cart.Suppliers[i].ProductList, cart.Suppliers[i])
		}

		total := 0
		for _, supplier := range cart.Suppliers {
			total += supplier.TotalAmount
		}
		cart.Total = total

		return nil
	})
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func getTotalSupplier(ctx context.Context, products []model.Product, supplier model.Supplier) int {
	total := 0
	for _, product := range products {
		total += product.Price * product.Quantity
	}

	if total < supplier.FreeDeliveryAmount {
		total += supplier.DeliveryFee
	}

	return total
}

func (s *cartServ) DeleteProductFromCart(ctx context.Context, query *model.PutCartQuery) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		cart, err := s.cartRepo.Cart(ctx, query.CustomerID)
		if err != nil {
			return err
		}
		query.CartID = cart.ID

		existingQuantity, err := s.cartRepo.ItemQuantity(ctx, query.CartID, query.ProductID, query.SupplierID)
		if err != nil {
			return err
		}

		price, err := s.productService.ProductPriceBySupplier(ctx, query.ProductID, query.SupplierID)
		if err != nil {
			return err
		}
		query.Price = price

		if query.Quantity >= existingQuantity {
			itemTotal := price * existingQuantity
			err = s.cartRepo.DeleteItem(ctx, query.CartID, query.ProductID, query.SupplierID)
			if err != nil {
				return err
			}
			cart.Total -= itemTotal
		} else {
			newQuantity := existingQuantity - query.Quantity
			err = s.cartRepo.UpdateItemQuantity(ctx, query.CartID, query.ProductID, query.SupplierID, newQuantity)
			if err != nil {
				return err
			}
			itemTotal := price * query.Quantity
			cart.Total -= itemTotal
		}

		if cart.Total < 0 {
			cart.Total = 0
		}
		return s.cartRepo.UpdateCartTotal(ctx, cart.ID, cart.Total)
	})

}

func (s *cartServ) ClearCart(ctx context.Context, userID int64) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		cart, err := s.cartRepo.Cart(ctx, userID)
		if err != nil {
			return err
		}
		err = s.cartRepo.DeleteCartItems(ctx, cart.ID)
		if err != nil {
			return err
		}
		return nil
	})
}
