package service

import (
	"context"
	"diploma/internal/testutils"
	"diploma/modules/product/model"
	"diploma/pkg/client/db"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockProductRepository struct {
	mock.Mock
}

func (m *mockProductRepository) GetProduct(ctx context.Context, id int64) (*model.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Product), args.Error(1)
}

func (m *mockProductRepository) GetSupplierProductListByProduct(ctx context.Context, id int64) ([]model.ProductSupplier, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.ProductSupplier), args.Error(1)
}

func (m *mockProductRepository) GetProductListByIDList(ctx context.Context, idList []int64) ([]*model.Product, error) {
	args := m.Called(ctx, idList)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Product), args.Error(1)
}

func (m *mockProductRepository) GetProductList(ctx context.Context, query *model.ProductListQuery) ([]model.Product, error) {
	args := m.Called(ctx, query)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Product), args.Error(1)
}

func (m *mockProductRepository) GetTotalProducts(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *mockProductRepository) GetProductPriceBySupplier(ctx context.Context, productID, supplierID int64) (int, error) {
	args := m.Called(ctx, productID, supplierID)
	return args.Int(0), args.Error(1)
}

type mockTxManager struct {
	mock.Mock
}

func (m *mockTxManager) ReadCommitted(ctx context.Context, h db.Handler) error {
	args := m.Called(ctx, h)
	return args.Error(0)
}

type ProductServiceTestSuite struct {
	suite.Suite
	service    *ProductService
	repository *mockProductRepository
	txManager  *mockTxManager
	helper     *testutils.AssertTestHelper
}

func TestProductService(t *testing.T) {
	suite.Run(t, new(ProductServiceTestSuite))
}

func (s *ProductServiceTestSuite) SetupTest() {
	s.repository = new(mockProductRepository)
	s.txManager = new(mockTxManager)
	s.helper = testutils.NewAssertTestHelper(s.T())

	s.service = NewService(s.repository, s.txManager)
}

func (s *ProductServiceTestSuite) TestGetProduct_Success() {
	now := time.Now()
	expectedProduct := &model.Product{
		ID:        1,
		GTIN:      123456789,
		Name:      "Test Product",
		ImageUrl:  "http://example.com/image.jpg",
		CreatedAt: now,
		UpdatedAt: now,
		LowestProductSupplier: model.ProductSupplier{
			Price:      100,
			SellAmount: 10,
			Supplier: model.Supplier{
				ID:                 1,
				Name:               "Test Supplier",
				OrderAmount:        1000,
				FreeDeliveryAmount: 2000,
				DeliveryFee:        500,
			},
		},
	}

	s.repository.On("GetProduct", mock.Anything, int64(1)).Return(expectedProduct, nil)

	product, err := s.service.productRepository.GetProduct(context.Background(), 1)

	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(product)
	s.helper.AssertEqual(expectedProduct.ID, product.ID)
	s.helper.AssertEqual(expectedProduct.Name, product.Name)
	s.helper.AssertEqual(expectedProduct.GTIN, product.GTIN)
	s.helper.AssertEqual(expectedProduct.ImageUrl, product.ImageUrl)
}

func (s *ProductServiceTestSuite) TestGetProductList_Success() {
	query := &model.ProductListQuery{
		Offset: 0,
		Limit:  10,
	}

	now := time.Now()
	expectedProducts := []model.Product{
		{
			ID:        1,
			GTIN:      123456789,
			Name:      "Product 1",
			ImageUrl:  "http://example.com/image1.jpg",
			CreatedAt: now,
			UpdatedAt: now,
			LowestProductSupplier: model.ProductSupplier{
				Price:      100,
				SellAmount: 10,
				Supplier: model.Supplier{
					ID:                 1,
					Name:               "Supplier 1",
					OrderAmount:        1000,
					FreeDeliveryAmount: 2000,
					DeliveryFee:        500,
				},
			},
		},
		{
			ID:        2,
			GTIN:      987654321,
			Name:      "Product 2",
			ImageUrl:  "http://example.com/image2.jpg",
			CreatedAt: now,
			UpdatedAt: now,
			LowestProductSupplier: model.ProductSupplier{
				Price:      200,
				SellAmount: 20,
				Supplier: model.Supplier{
					ID:                 2,
					Name:               "Supplier 2",
					OrderAmount:        2000,
					FreeDeliveryAmount: 3000,
					DeliveryFee:        600,
				},
			},
		},
	}

	s.repository.On("GetProductList", mock.Anything, query).Return(expectedProducts, nil)

	products, err := s.service.productRepository.GetProductList(context.Background(), query)

	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(products)
	s.helper.AssertEqual(len(expectedProducts), len(products))

	for i, product := range products {
		s.helper.AssertEqual(expectedProducts[i].ID, product.ID)
		s.helper.AssertEqual(expectedProducts[i].Name, product.Name)
		s.helper.AssertEqual(expectedProducts[i].GTIN, product.GTIN)
		s.helper.AssertEqual(expectedProducts[i].ImageUrl, product.ImageUrl)
	}
}

func (s *ProductServiceTestSuite) TestGetProductPriceBySupplier_Success() {
	expectedPrice := 100

	s.repository.On("GetProductPriceBySupplier", mock.Anything, int64(1), int64(1)).Return(expectedPrice, nil)

	price, err := s.service.productRepository.GetProductPriceBySupplier(context.Background(), 1, 1)

	s.helper.AssertNoError(err)
	s.helper.AssertEqual(expectedPrice, price)
}

func (s *ProductServiceTestSuite) TestGetTotalProducts_Success() {
	expectedTotal := 100

	s.repository.On("GetTotalProducts", mock.Anything).Return(expectedTotal, nil)

	total, err := s.service.productRepository.GetTotalProducts(context.Background())

	s.helper.AssertNoError(err)
	s.helper.AssertEqual(expectedTotal, total)
}
