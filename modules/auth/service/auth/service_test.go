package auth

import (
	"context"
	"diploma/internal/testutils"
	"diploma/modules/auth/model"
	"diploma/pkg/client/db"
	"errors"
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
	if h != nil {
		_ = h(ctx)
	}
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

func (s *AuthServiceTestSuite) TestRegister_Success() {
	newUser := &model.AuthUser{
		Info: &model.UserInfo{
			Name:        "testuser",
			PhoneNumber: "+1234567890",
			Role:        model.CustomerRole,
		},
		Password: "password123",
	}

	expectedID := int64(1)
	expectedUser := &model.User{
		ID:        expectedID,
		Info:      newUser.Info,
		CreatedAt: time.Now(),
	}

	s.authRepository.On("Create", mock.Anything, mock.MatchedBy(func(u *model.AuthUser) bool {
		return u.Info.Name == newUser.Info.Name && u.Info.PhoneNumber == newUser.Info.PhoneNumber
	})).Return(expectedID, nil).Once()

	s.authRepository.On("GetById", mock.Anything, expectedID).Return(expectedUser, nil).Once()

	s.txManager.On("ReadCommitted", mock.Anything, mock.AnythingOfType("db.Handler")).Return(nil).Once()

	id, err := s.service.Register(context.Background(), newUser)

	s.helper.AssertNoError(err)
	s.helper.AssertEqual(expectedID, id)
	s.authRepository.AssertExpectations(s.T())
	s.txManager.AssertExpectations(s.T())
}

func (s *AuthServiceTestSuite) TestRegister_AdminRole_Failure() {
	newUser := &model.AuthUser{
		Info: &model.UserInfo{
			Name:        "admin",
			PhoneNumber: "+1234567890",
			Role:        model.AdminRole,
		},
		Password: "password123",
	}

	_, err := s.service.Register(context.Background(), newUser)

	s.helper.AssertError(err)
	s.helper.AssertEqual("admin role is not allowed", err.Error())
	s.authRepository.AssertNotCalled(s.T(), "Create")
}

func (s *AuthServiceTestSuite) TestLogin_Success() {
	phoneNumber := "+1234567890"
	password := "password123"
	hashedPassword := "$2a$10$abcdefghijklmnopqrstuvwxyz" // This would be a real bcrypt hash
	userID := int64(1)
	userName := "testuser"
	userRole := model.CustomerRole

	storedUser := &model.AuthUser{
		ID: userID,
		Info: &model.UserInfo{
			Name:        userName,
			PhoneNumber: phoneNumber,
			Role:        userRole,
		},
		HashedPassword: hashedPassword,
	}

	expectedAccessToken := "access.token.123"
	expectedRefreshToken := "refresh.token.123"

	s.authRepository.On("GetByPhoneNumber", mock.Anything, phoneNumber).Return(storedUser, nil).Once()
	s.jwt.On("GenerateJSONWebTokens", userID, userName, userRole).Return(expectedAccessToken, expectedRefreshToken, nil).Once()

	accessToken, refreshToken, err := s.service.Login(context.Background(), phoneNumber, password)

	s.helper.AssertNoError(err)
	s.helper.AssertEqual(expectedAccessToken, accessToken)
	s.helper.AssertEqual(expectedRefreshToken, refreshToken)
	s.authRepository.AssertExpectations(s.T())
	s.jwt.AssertExpectations(s.T())
}

func (s *AuthServiceTestSuite) TestLogin_InvalidCredentials() {
	phoneNumber := "+1234567890"
	password := "wrongpassword"

	s.authRepository.On("GetByPhoneNumber", mock.Anything, phoneNumber).Return(nil, model.ErrNoRows).Once()

	accessToken, refreshToken, err := s.service.Login(context.Background(), phoneNumber, password)

	s.helper.AssertError(err)
	s.helper.AssertEqual(model.ErrInvalidCredentials, err)
	s.helper.AssertEqual("", accessToken)
	s.helper.AssertEqual("", refreshToken)
	s.authRepository.AssertExpectations(s.T())
	s.jwt.AssertNotCalled(s.T(), "GenerateJSONWebTokens")
}

func (s *AuthServiceTestSuite) TestRegister_RepositoryError() {
	newUser := &model.AuthUser{
		Info: &model.UserInfo{
			Name:        "testuser",
			PhoneNumber: "+1234567890",
			Role:        model.CustomerRole,
		},
		Password: "password123",
	}

	expectedError := errors.New("database error")
	s.authRepository.On("Create", mock.Anything, mock.MatchedBy(func(u *model.AuthUser) bool {
		return u.Info.Name == newUser.Info.Name && u.Info.PhoneNumber == newUser.Info.PhoneNumber
	})).Return(int64(0), expectedError).Once()

	s.txManager.On("ReadCommitted", mock.Anything, mock.AnythingOfType("db.Handler")).Return(expectedError).Once()

	id, err := s.service.Register(context.Background(), newUser)

	s.helper.AssertError(err)
	s.helper.AssertEqual(int64(0), id)
	s.helper.AssertEqual(expectedError, err)
	s.authRepository.AssertExpectations(s.T())
	s.txManager.AssertExpectations(s.T())
}
