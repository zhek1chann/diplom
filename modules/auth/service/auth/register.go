package auth

import (
	"context"
	"errors"

	"diploma/modules/auth/model"
	passwordutil "diploma/pkg/password-util"
)

func (s *authServ) Register(ctx context.Context, user *model.AuthUser) (int64, error) {
	var id int64
	var err error
	user.HashedPassword, err = passwordutil.HashPassword(user.Password)
	if err != nil {
		return 0, err
	}
	if user.Info.Role == model.AdminRole {
		return 0, errors.New("admin role is not allowed")
	}
	err = s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		var errTx error
		id, errTx = s.authRepository.Create(ctx, user)
		if errTx != nil {
			return errTx
		}

		_, errTx = s.authRepository.GetById(ctx, id)
		if errTx != nil {
			return errTx
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *authServ) Login(ctx context.Context, phoneNumber string, password string) (accessToken string, refreshToken string, err error) {
	authUser, err := s.authRepository.GetByPhoneNumber(ctx, phoneNumber)

	if err != nil {
		if errors.Is(err, model.ErrNoRows) {
			return "", "", model.ErrInvalidCredentials
		}
	}

	if !passwordutil.CheckPasswordHash(password, authUser.HashedPassword) {
		return "", "", model.ErrInvalidCredentials
	}
	return s.jwt.GenerateJSONWebTokens(authUser.ID, authUser.Info.Name, authUser.Info.Role)

}
