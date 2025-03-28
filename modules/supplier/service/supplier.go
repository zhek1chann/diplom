package service

import (
	"context"
	"diploma/modules/supplier/model"
)

func (s *SupplierService) SupplierListByIDList(ctx context.Context, id []int64) ([]model.Supplier, error) {
	return s.supplierRepo.SupplierListByIDList(ctx, id)
}
