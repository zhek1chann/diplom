package service

import (
	"context"
	"diploma/internal/testutils"
	"diploma/modules/cart/model"
	"diploma/pkg/client/db"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockCartRepository struct {
	mock.Mock
}

func (m *mockCartRepository) UpdateItemQuantity(ctx context.Context, cartId, productId, supplierId int64, quantity int) error {
	args := m.Called(ctx, cartId, productId, supplierId, quantity)
	return args.Error(0)
}

func (m *mockCartRepository) ItemQuantity(ctx context.Context, cartId, productId, supplierId int64) (int, error) {
	args := m.Called(ctx, cartId, productId, supplierId)
	return args.Int(0), args.Error(1)
}

func (m *mockCartRepository) Cart(ctx context.Context, userID int64) (*model.Cart, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Cart), args.Error(1)
}

func (m *mockCartRepository) CreateCart(ctx context.Context, userID int64) (int64, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockCartRepository) AddItem(ctx context.Context, input *model.PutCartQuery) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *mockCartRepository) UpdateCartTotal(ctx context.Context, cartID int64, total int) error {
	args := m.Called(ctx, cartID, total)
	return args.Error(0)
}

func (m *mockCartRepository) DeleteCart(ctx context.Context, cartID int64) error {
	args := m.Called(ctx, cartID)
	return args.Error(0)
}

func (m *mockCartRepository) GetCartItems(ctx context.Context, cartID int64) ([]model.Supplier, error) {
	args := m.Called(ctx, cartID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Supplier), args.Error(1)
}

func (m *mockCartRepository) DeleteCartItems(ctx context.Context, cartID int64) error {
	args := m.Called(ctx, cartID)
	return args.Error(0)
}

func (m *mockCartRepository) DeleteItem(ctx context.Context, cartID, productId, supplierId int64) error {
	args := m.Called(ctx, cartID, productId, supplierId)
	return args.Error(0)
}

type mockProductService struct {
	mock.Mock
}

func (m *mockProductService) ProductPriceBySupplier(ctx context.Context, productID, supplierID int64) (int, error) {
	args := m.Called(ctx, productID, supplierID)
	return args.Int(0), args.Error(1)
}

type mockSupplierClient struct {
	mock.Mock
}

func (m *mockSupplierClient) SupplierListByIDList(ctx context.Context, IDList []int64) ([]model.Supplier, error) {
	args := m.Called(ctx, IDList)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Supplier), args.Error(1)
}

type mockOrderClient struct {
	mock.Mock
}

func (m *mockOrderClient) CreateOrder(ctx context.Context, cart *model.Cart) error {
	args := m.Called(ctx, cart)
	return args.Error(0)
}

type mockPaymentClient struct {
	mock.Mock
}

func (m *mockPaymentClient) PaymentRequest(orderID, amount, currency, description string) (model.CheckoutResponse, error) {
	args := m.Called(orderID, amount, currency, description)
	return args.Get(0).(model.CheckoutResponse), args.Error(1)
}

type mockRedis struct {
	mock.Mock
}

func (m *mockRedis) SavePaymentOrder(ctx context.Context, paymentOrder model.PaymentOrder) error {
	args := m.Called(ctx, paymentOrder)
	return args.Error(0)
}

func (m *mockRedis) PaymentOrder(ctx context.Context, orderID string) (model.PaymentOrder, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return model.PaymentOrder{}, args.Error(1)
	}
	return args.Get(0).(model.PaymentOrder), args.Error(1)
}

type mockTxManager struct {
	mock.Mock
}

func (m *mockTxManager) ReadCommitted(ctx context.Context, h db.Handler) error {
	args := m.Called(ctx, h)
	return args.Error(0)
}

type CartServiceTestSuite struct {
	suite.Suite
	service         *cartServ
	cartRepo        *mockCartRepository
	productService  *mockProductService
	supplierService *mockSupplierClient
	orderService    *mockOrderClient
	paymentClient   *mockPaymentClient
	redis           *mockRedis
	txManager       *mockTxManager
	helper          *testutils.AssertTestHelper
}

func TestCartService(t *testing.T) {
	suite.Run(t, new(CartServiceTestSuite))
}

func (s *CartServiceTestSuite) SetupTest() {
	s.cartRepo = new(mockCartRepository)
	s.productService = new(mockProductService)
	s.supplierService = new(mockSupplierClient)
	s.orderService = new(mockOrderClient)
	s.paymentClient = new(mockPaymentClient)
	s.redis = new(mockRedis)
	s.txManager = new(mockTxManager)
	s.helper = testutils.NewAssertTestHelper(s.T())

	s.service = NewService(
		s.cartRepo,
		s.productService,
		s.supplierService,
		s.orderService,
		s.paymentClient,
		s.redis,
		s.txManager,
	)
}

func (s *CartServiceTestSuite) TestCreateCart_Success() {
	userID := int64(1)
	expectedCartID := int64(1)

	s.cartRepo.On("CreateCart", mock.Anything, userID).Return(expectedCartID, nil)

	cartID, err := s.service.cartRepo.CreateCart(context.Background(), userID)

	s.helper.AssertNoError(err)
	s.helper.AssertEqual(expectedCartID, cartID)
}

func (s *CartServiceTestSuite) TestAddItem_Success() {
	input := &model.PutCartQuery{
		CartID:     1,
		ProductID:  1,
		SupplierID: 1,
		CustomerID: 1,
		Quantity:   2,
		Price:      100,
	}

	s.cartRepo.On("AddItem", mock.Anything, input).Return(nil)

	err := s.service.cartRepo.AddItem(context.Background(), input)

	s.helper.AssertNoError(err)
}

func (s *CartServiceTestSuite) TestGetCart_Success() {
	userID := int64(1)
	expectedCart := &model.Cart{
		ID:         1,
		CustomerID: userID,
		Total:      100,
		Suppliers: []model.Supplier{
			{
				ID:                 1,
				Name:               "Test Supplier",
				OrderAmount:        1000,
				TotalAmount:        100,
				FreeDeliveryAmount: 2000,
				DeliveryFee:        500,
				ProductList: []model.Product{
					{
						ID:       1,
						Name:     "Test Product",
						Price:    100,
						Quantity: 1,
						ImageUrl: "http://example.com/image.jpg",
					},
				},
			},
		},
	}

	s.cartRepo.On("Cart", mock.Anything, userID).Return(expectedCart, nil)

	cart, err := s.service.cartRepo.Cart(context.Background(), userID)

	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(cart)
	s.helper.AssertEqual(expectedCart.ID, cart.ID)
	s.helper.AssertEqual(expectedCart.CustomerID, cart.CustomerID)
	s.helper.AssertEqual(expectedCart.Total, cart.Total)
	s.helper.AssertEqual(len(expectedCart.Suppliers), len(cart.Suppliers))
}

func (s *CartServiceTestSuite) TestDeleteCart_Success() {
	cartID := int64(1)

	s.cartRepo.On("DeleteCart", mock.Anything, cartID).Return(nil)

	err := s.service.cartRepo.DeleteCart(context.Background(), cartID)

	s.helper.AssertNoError(err)
}
