package service

import (
	"context"
	"diploma/modules/contract/model"
	"time"
)

type Repository interface {
	Create(ctx context.Context, contract *model.Contract) (int64, error)
	SignByParty(ctx context.Context, contractID int64, role int, signature string) error
	GetByID(ctx context.Context, id int64) (*model.Contract, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateContract(ctx context.Context, orderID, supplierID, customerID int64, content string) (int64, error) {
	contract := &model.Contract{
		OrderID:    orderID,
		SupplierID: supplierID,
		CustomerID: customerID,
		Content:    content,
		Status:     model.StatusCreated,
		CreatedAt:  time.Now(),
	}
	return s.repo.Create(ctx, contract)
}

func (s *Service) SignContract(ctx context.Context, contractID int64, role int, signature string) error {
	return s.repo.SignByParty(ctx, contractID, role, signature)
}

func (s *Service) GetContract(ctx context.Context, id int64) (*model.Contract, error) {
	return s.repo.GetByID(ctx, id)
}
