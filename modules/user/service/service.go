package auth

import (
	"diploma/pkg/client/db"
)

type userServ struct {
	userRepository IUserRepository
	txManager      db.TxManager
}

func NewService(
	userRepository IUserRepository,
	txManager db.TxManager,
) *userServ {
	return &userServ{
		userRepository: userRepository,
		txManager:      txManager,
	}
}

type IUserRepository interface {
	IAddressRepo
}
