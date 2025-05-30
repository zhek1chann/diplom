package auth

import (
	"context"
	"diploma/internal/testutils"
	"diploma/modules/auth/model"
	"diploma/pkg/client/db"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type mockAuthRepository struct {
	mock.Mock
}

func (m *mockAuthRepository) Create(ctx context.Context, user *model.AuthUser) (int64, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockAuthRepository) GetById(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *mockAuthRepository) GetByPhoneNumber(ctx context.Context, phoneNumber string) (*model.AuthUser, error) {
	args := m.Called(ctx, phoneNumber)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.AuthUser), args.Error(1)
}

type mockTxManager struct {
	mock.Mock
}

func (m *mockTxManager) ReadCommitted(ctx context.Context, h db.Handler) error {
	args := m.Called(ctx, h)
	return args.Error(0)
}

type mockJWT struct {
	mock.Mock
}

func (m *mockJWT) GenerateJSONWebTokens(id int64, username string, role int) (string, string, error) {
	args := m.Called(id, username, role)
	return args.String(0), args.String(1), args.Error(2)
}

func (m *mockJWT) RefreshAccessToken(refreshToken string) (string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.Error(1)
}

type AuthServiceTestSuite struct {
	suite.Suite
	service        *authServ
	authRepository *mockAuthRepository
	jwt            *mockJWT
	txManager      *mockTxManager
	helper         *testutils.AssertTestHelper
}

func TestAuthService(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

func (s *AuthServiceTestSuite) SetupTest() {
	s.authRepository = new(mockAuthRepository)
	s.jwt = new(mockJWT)
	s.txManager = new(mockTxManager)
	s.helper = testutils.NewAssertTestHelper(s.T())
	s.service = NewService(s.authRepository, s.jwt, s.txManager)
}

func (s *AuthServiceTestSuite) TestGetUserById_Success() {
	expectedUser := &model.User{
		ID: 1,
		Info: &model.UserInfo{
			Name:        "testuser",
			PhoneNumber: "+1234567890",
			Role:        model.CustomerRole,
		},
		CreatedAt: time.Now(),
	}

	s.authRepository.On("GetById", mock.Anything, int64(1)).Return(expectedUser, nil)

	user, err := s.service.authRepository.GetById(context.Background(), 1)

	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(user)
	s.helper.AssertEqual(expectedUser.ID, user.ID)
	s.helper.AssertEqual(expectedUser.Info.Name, user.Info.Name)
	s.helper.AssertEqual(expectedUser.Info.Role, user.Info.Role)
}

func (s *AuthServiceTestSuite) TestGetByPhoneNumber_Success() {
	expectedUser := &model.AuthUser{
		ID: 1,
		Info: &model.UserInfo{
			Name:        "testuser",
			PhoneNumber: "+1234567890",
			Role:        model.CustomerRole,
		},
		Password: "hashedpassword",
	}

	s.authRepository.On("GetByPhoneNumber", mock.Anything, "+1234567890").Return(expectedUser, nil)

	user, err := s.service.authRepository.GetByPhoneNumber(context.Background(), "+1234567890")

	s.helper.AssertNoError(err)
	s.helper.AssertNotNil(user)
	s.helper.AssertEqual(expectedUser.ID, user.ID)
	s.helper.AssertEqual(expectedUser.Info.PhoneNumber, user.Info.PhoneNumber)
	s.helper.AssertEqual(expectedUser.Password, user.Password)
}

func (s *AuthServiceTestSuite) TestCreate_Success() {
	newUser := &model.AuthUser{
		Info: &model.UserInfo{
			Name:        "testuser",
			PhoneNumber: "+1234567890",
			Role:        model.CustomerRole,
		},
		Password: "hashedpassword",
	}

	expectedID := int64(1)
	s.authRepository.On("Create", mock.Anything, newUser).Return(expectedID, nil)

	id, err := s.service.authRepository.Create(context.Background(), newUser)

	s.helper.AssertNoError(err)
	s.helper.AssertEqual(expectedID, id)
}
