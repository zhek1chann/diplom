package service

import (
	"context"
	"diploma/modules/cart/model"
	"diploma/pkg/client/db"
)

type cartServ struct {
	cartRepo        ICartRepository
	productService  IProductService
	supplierService ISupplierClient
	orderService    IOrderClient
	txManager       db.TxManager
}

func NewService(
	cartRepository ICartRepository,
	productService IProductService,
	supplierService ISupplierClient,
	OrderClient IOrderClient,
	txManager db.TxManager,
) *cartServ {
	return &cartServ{
		cartRepo:        cartRepository,
		supplierService: supplierService,
		orderService:    OrderClient,
		txManager:       txManager,
		productService:  productService,
	}
}

type IProductService interface {
	ProductPriceBySupplier(ctx context.Context, productID, supplierID int64) (int, error)
}

type ISupplierClient interface {
	SupplierListByIDList(ctx context.Context, IDList []int64) ([]model.Supplier, error)
}

type IOrderClient interface {
	CreateOrder(ctx context.Context, cart *model.Cart) error
}

type ICartRepository interface {
	UpdateItemQuantity(ctx context.Context, cartId, productId, supplierId int64, quantity int) error
	ItemQuantity(ctx context.Context, cartId, productId, supplierId int64) (int, error)
	Cart(ctx context.Context, userID int64) (*model.Cart, error)
	CreateCart(ctx context.Context, userID int64) (int64, error)
	AddItem(ctx context.Context, input *model.PutCartQuery) error
	UpdateCartTotal(ctx context.Context, cartID int64, total int) error
	DeleteCart(ctx context.Context, cartID int64) error
	GetCartItems(ctx context.Context, cartID int64) ([]model.Supplier, error)
	DeleteCartItems(ctx context.Context, cartID int64) error
	DeleteItem(ctx context.Context, cartID, productId, supplierId int64) error
}
