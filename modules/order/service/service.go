package service

import (
	"context"
	"diploma/modules/order/model"
	"diploma/pkg/client/db"
)

type OrderService struct {
	orderRepo      IOrderRepository
	supplierClient ISupplierClient
	productClient  IProductClient
	txManager      db.TxManager
}

func NewService(repo IOrderRepository, supplierClient ISupplierClient, productClient IProductClient, tx db.TxManager) *OrderService {
	return &OrderService{
		supplierClient: supplierClient,
		productClient:  productClient,
		orderRepo:      repo,
		txManager:      tx,
	}

}

type IOrderRepository interface {
	ICreateOrderRepo
	IOrderRepo
	UpdateOrderStatus(ctx context.Context, orderID int64, newStatus int) error
	GetOrderByID(ctx context.Context, orderID int64) (*model.Order, error)
}

type ISupplierClient interface {
	Supplier(ctx context.Context, id int64) (*model.Supplier, error)
}

type IProductClient interface {
	Product(ctx context.Context, id int64) (*model.Product, error)
}
