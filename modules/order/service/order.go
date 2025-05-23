package service

import (
	"context"
	"diploma/modules/order/model"
	"fmt"
	"time"
)

type IOrderRepo interface {
	OrdersByUserID(ctx context.Context, userID int64) ([]*model.Order, error)
	OrdersBySupplierID(ctx context.Context, supplierID int64) ([]*model.Order, error)
	OrderProducts(ctx context.Context, orderID int64) ([]*model.OrderProduct, error)
}

func (s *OrderService) Orders(ctx context.Context, userID int64, role int) ([]*model.Order, error) {
	// Retrieve the orders for the user
	var orders []*model.Order
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		switch role {
		case model.CustomerRole:
			orders, errTx = s.orderRepo.OrdersByUserID(ctx, userID)
			if errTx != nil {
				return errTx
			}
			for _, order := range orders {
				order.Supplier, errTx = s.supplierClient.Supplier(ctx, order.SupplierID)

				if errTx != nil {
					return errTx
				}
				order.Supplier.ID = order.SupplierID
			}

		case model.SupplierRole:
			orders, errTx = s.orderRepo.OrdersBySupplierID(ctx, userID)
			if errTx != nil {
				return errTx
			}
			for _, order := range orders {
				order.Supplier, errTx = s.supplierClient.Supplier(ctx, order.SupplierID)
				if errTx != nil {
					return errTx
				}
			}
		}
		for _, order := range orders {
			orderProducts, errTx := s.orderRepo.OrderProducts(ctx, order.ID)
			if errTx != nil {
				return errTx
			}
			order.ProductList = orderProducts
			for _, orderProduct := range orderProducts {
				product, errTx := s.productClient.Product(ctx, orderProduct.ProductID)
				if errTx != nil {
					return errTx
				}
				orderProduct.Product = product
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return orders, nil
}

type ICreateOrderRepo interface {
	CreateOrder(ctx context.Context, order *model.Order) (int64, error)
	CreateOrderProduct(ctx context.Context, orderProduct *model.OrderProduct) error
}

func (s *OrderService) CreateOrder(ctx context.Context, orders []*model.Order) error {

	// Convert the cart to a list of orders
	ordersID := make([]int64, len(orders))
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		for _, order := range orders {
			order.OrderDate = time.Now().Add(72 * time.Hour)
			order.StatusID = 1
			id, errTx := s.orderRepo.CreateOrder(ctx, order)
			if errTx != nil {
				return errTx
			}
			for _, op := range order.ProductList {
				op.OrderID = id
				errTx = s.orderRepo.CreateOrderProduct(ctx, op)
				fmt.Println(errTx)
				if errTx != nil {
					return errTx
				}
			}
			if errTx != nil {
				return errTx
			}
			ordersID = append(ordersID, id)
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) UpdateOrderStatusBySupplier(ctx context.Context, supplierID int64, orderID int64, newStatusID int) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		order, err := s.orderRepo.GetOrderByID(ctx, orderID)
		if err != nil {
			return err
		}

		if order.SupplierID != supplierID {
			return fmt.Errorf("supplier %d does not own order %d", supplierID, orderID)
		}

		// Проверка допустимости перехода статуса
		valid := false
		switch order.StatusID {
		case model.Pending:
			if newStatusID == model.InProgress {
				valid = true
			}
		case model.InProgress:
			if newStatusID == model.Completed || newStatusID == model.Cancelled {
				valid = true
			}
		}

		if !valid {
			return fmt.Errorf("invalid status transition from %d to %d", order.StatusID, newStatusID)
		}

		// Обновляем статус
		err = s.orderRepo.UpdateOrderStatus(ctx, orderID, newStatusID)
		if err != nil {
			return err
		}

		// Если переходим в InProgress — создаём контракт
		if order.StatusID == model.Pending && newStatusID == model.InProgress {
			_, err := s.contractClient.CreateContract(
				ctx,
				order.ID,
				order.SupplierID,
				order.CustomerID,
				fmt.Sprintf("Контракт для заказа #%d", order.ID),
			)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *OrderService) GetOrderByID(ctx context.Context, userID int64, role int, orderID int64) (*model.Order, error) {
	var order *model.Order
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		order, errTx = s.orderRepo.GetOrderByID(ctx, orderID)
		if errTx != nil {
			return errTx
		}

		switch role {
		case model.CustomerRole:
			if order.CustomerID != userID {
				return fmt.Errorf("unauthorized: customer can only view their own orders")
			}
		case model.SupplierRole:
			if order.SupplierID != userID {
				return fmt.Errorf("unauthorized: supplier can only view their own orders")
			}
		case model.AdminRole:
		default:
			return fmt.Errorf("unauthorized role")
		}

		order.Supplier, errTx = s.supplierClient.Supplier(ctx, order.SupplierID)
		if errTx != nil {
			return errTx
		}
		orderProducts, errTx := s.orderRepo.OrderProducts(ctx, order.ID)
		if errTx != nil {
			return errTx
		}
		order.ProductList = orderProducts
		for _, op := range order.ProductList {
			op.Product, errTx = s.productClient.Product(ctx, op.ProductID)
			if errTx != nil {
				return errTx
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *OrderService) CancelOrderByCustomer(ctx context.Context, customerID int64, orderID int64) error {
	return s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		order, err := s.orderRepo.GetOrderByID(ctx, orderID)
		if err != nil {
			return err
		}

		if order.CustomerID != customerID {
			return fmt.Errorf("customer %d does not own order %d", customerID, orderID)
		}

		if order.StatusID != model.Pending {
			return fmt.Errorf("only orders in Pending status can be cancelled")
		}

		return s.orderRepo.UpdateOrderStatus(ctx, orderID, model.Cancelled)
	})
}
