package service

import (
	"context"
	"diploma/modules/cart/model"
	"errors"

	"go.uber.org/zap"
)

func (s *cartServ) AddProductToCard(ctx context.Context, query *model.PutCartQuery) error {
	s.LogInfo(ctx, "Adding product to cart",
		zap.Int64("customer_id", query.CustomerID),
		zap.Int64("product_id", query.ProductID),
		zap.Int64("supplier_id", query.SupplierID),
		zap.Int("quantity", query.Quantity),
	)

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		cart, errTx := s.cartRepo.Cart(ctx, query.CustomerID)

		if errTx != nil {
			if errors.Is(errTx, model.ErrNoRows) {
				s.LogInfo(ctx, "Cart not found, creating new cart", zap.Int64("customer_id", query.CustomerID))
				id, errTx := s.cartRepo.CreateCart(ctx, query.CustomerID)
				if errTx != nil {
					s.LogError(ctx, "Failed to create cart", errTx)
					return errTx
				}
				cart = &model.Cart{
					ID:    id,
					Total: 0,
				}
			} else {
				s.LogError(ctx, "Failed to get cart", errTx)
				return errTx
			}
		}
		query.CartID = cart.ID

		query.Price, errTx = s.productService.ProductPriceBySupplier(ctx, query.ProductID, query.SupplierID)
		if errTx != nil {
			s.LogError(ctx, "Failed to get product price", errTx,
				zap.Int64("product_id", query.ProductID),
				zap.Int64("supplier_id", query.SupplierID),
			)
			return errTx
		}

		itemQuantity, errTx := s.cartRepo.ItemQuantity(ctx, query.CartID, query.ProductID, query.SupplierID)
		if errTx != nil {
			if errors.Is(errTx, model.ErrNoRows) {
				s.LogInfo(ctx, "Adding new item to cart",
					zap.Int64("cart_id", query.CartID),
					zap.Int64("product_id", query.ProductID),
				)
				errTx = s.cartRepo.AddItem(ctx, query)
				if errTx != nil {
					s.LogError(ctx, "Failed to add item to cart", errTx)
					return errTx
				}
			} else {
				s.LogError(ctx, "Failed to get item quantity", errTx)
				return errTx
			}
		} else {
			s.LogInfo(ctx, "Updating existing item quantity",
				zap.Int64("cart_id", query.CartID),
				zap.Int64("product_id", query.ProductID),
				zap.Int("old_quantity", itemQuantity),
				zap.Int("new_quantity", itemQuantity+query.Quantity),
			)
			errTx = s.cartRepo.UpdateItemQuantity(ctx, query.CartID, query.ProductID, query.SupplierID, itemQuantity+query.Quantity)
			if errTx != nil {
				s.LogError(ctx, "Failed to update item quantity", errTx)
				return errTx
			}
		}

		cart.Total += query.Price * query.Quantity
		errTx = s.cartRepo.UpdateCartTotal(ctx, cart.ID, cart.Total)
		if errTx != nil {
			s.LogError(ctx, "Failed to update cart total", errTx)
			return errTx
		}
		return nil
	})
	if err != nil {
		return err
	}

	s.LogInfo(ctx, "Successfully added product to cart",
		zap.Int64("customer_id", query.CustomerID),
		zap.Int64("product_id", query.ProductID),
		zap.Int("quantity", query.Quantity),
	)
	return nil
}

func (s *cartServ) Cart(ctx context.Context, userID int64) (*model.Cart, error) {
	s.LogInfo(ctx, "Fetching cart", zap.Int64("user_id", userID))

	var err error
	var cart *model.Cart
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		cart, errTx = s.cartRepo.Cart(ctx, userID)
		if errTx != nil {
			if errors.Is(errTx, model.ErrNoRows) {
				s.LogInfo(ctx, "No existing cart found, creating empty cart", zap.Int64("user_id", userID))
				cart = &model.Cart{
					ID:         0,
					CustomerID: userID,
					Total:      0,
					Suppliers:  []model.Supplier{},
				}
				return nil
			}
			s.LogError(ctx, "Failed to get cart", errTx)
			return errTx
		}

		cart.Suppliers, errTx = s.cartRepo.GetCartItems(ctx, cart.ID)
		if errTx != nil {
			if errors.Is(errTx, model.ErrNoRows) {
				s.LogInfo(ctx, "No items in cart", zap.Int64("cart_id", cart.ID))
				return model.ErrNoRows
			}
			s.LogError(ctx, "Failed to get cart items", errTx)
			return errTx
		}

		supplierIdList := make([]int64, 0, len(cart.Suppliers))
		for _, supplier := range cart.Suppliers {
			supplierIdList = append(supplierIdList, supplier.ID)
		}

		suppliers, errTx := s.supplierService.SupplierListByIDList(ctx, supplierIdList)
		if errTx != nil {
			s.LogError(ctx, "Failed to get supplier details", errTx)
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

	s.LogInfo(ctx, "Successfully fetched cart",
		zap.Int64("user_id", userID),
		zap.Int64("cart_id", cart.ID),
		zap.Int("total_amount", cart.Total),
		zap.Int("supplier_count", len(cart.Suppliers)),
	)
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
	s.LogInfo(ctx, "Deleting product from cart",
		zap.Int64("customer_id", query.CustomerID),
		zap.Int64("product_id", query.ProductID),
		zap.Int64("supplier_id", query.SupplierID),
		zap.Int("quantity", query.Quantity),
	)

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		cart, err := s.cartRepo.Cart(ctx, query.CustomerID)
		if err != nil {
			s.LogError(ctx, "Failed to get cart", err)
			return err
		}
		query.CartID = cart.ID

		existingQuantity, err := s.cartRepo.ItemQuantity(ctx, query.CartID, query.ProductID, query.SupplierID)
		if err != nil {
			s.LogError(ctx, "Failed to get item quantity", err)
			return err
		}

		price, err := s.productService.ProductPriceBySupplier(ctx, query.ProductID, query.SupplierID)
		if err != nil {
			s.LogError(ctx, "Failed to get product price", err)
			return err
		}
		query.Price = price

		if query.Quantity >= existingQuantity {
			s.LogInfo(ctx, "Removing item completely from cart",
				zap.Int64("cart_id", query.CartID),
				zap.Int64("product_id", query.ProductID),
			)
			itemTotal := price * existingQuantity
			err = s.cartRepo.DeleteItem(ctx, query.CartID, query.ProductID, query.SupplierID)
			if err != nil {
				s.LogError(ctx, "Failed to delete item", err)
				return err
			}
			cart.Total -= itemTotal
		} else {
			newQuantity := existingQuantity - query.Quantity
			s.LogInfo(ctx, "Updating item quantity",
				zap.Int64("cart_id", query.CartID),
				zap.Int64("product_id", query.ProductID),
				zap.Int("old_quantity", existingQuantity),
				zap.Int("new_quantity", newQuantity),
			)
			err = s.cartRepo.UpdateItemQuantity(ctx, query.CartID, query.ProductID, query.SupplierID, newQuantity)
			if err != nil {
				s.LogError(ctx, "Failed to update item quantity", err)
				return err
			}
			itemTotal := price * query.Quantity
			cart.Total -= itemTotal
		}

		if cart.Total < 0 {
			cart.Total = 0
		}
		err = s.cartRepo.UpdateCartTotal(ctx, cart.ID, cart.Total)
		if err != nil {
			s.LogError(ctx, "Failed to update cart total", err)
			return err
		}

		s.LogInfo(ctx, "Successfully updated cart total",
			zap.Int64("cart_id", cart.ID),
			zap.Int("new_total", cart.Total),
		)
		return nil
	})

	if err != nil {
		return err
	}

	s.LogInfo(ctx, "Successfully deleted product from cart",
		zap.Int64("customer_id", query.CustomerID),
		zap.Int64("product_id", query.ProductID),
		zap.Int("quantity", query.Quantity),
	)
	return nil
}

func (s *cartServ) ClearCart(ctx context.Context, userID int64) error {
	s.LogInfo(ctx, "Clearing cart", zap.Int64("user_id", userID))

	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		cart, err := s.cartRepo.Cart(ctx, userID)
		if err != nil {
			s.LogError(ctx, "Failed to get cart", err)
			return err
		}

		err = s.cartRepo.DeleteCartItems(ctx, cart.ID)
		if err != nil {
			s.LogError(ctx, "Failed to delete cart items", err)
			return err
		}

		s.LogInfo(ctx, "Successfully cleared cart items",
			zap.Int64("cart_id", cart.ID),
			zap.Int64("user_id", userID),
		)
		return nil
	})

	if err != nil {
		return err
	}

	s.LogInfo(ctx, "Successfully cleared cart", zap.Int64("user_id", userID))
	return nil
}
