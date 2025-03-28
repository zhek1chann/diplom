package supplier

import (
	"context"
	"diploma/modules/cart/model"
	supplierModel "diploma/modules/supplier/model"
)

type SupplierAdapter struct {
	supplierService ISupplierService
}

func NewAdapter(supplierService ISupplierService) *SupplierAdapter {
	return &SupplierAdapter{supplierService: supplierService}
}

type ISupplierService interface {
	SupplierListByIDList(ctx context.Context, idList []int64) ([]supplierModel.Supplier, error)
}

func (a *SupplierAdapter) SupplierListByIDList(ctx context.Context, IDList []int64) ([]model.Supplier, error) {

	suppliers, err := a.supplierService.SupplierListByIDList(ctx, IDList)
	if err != nil {
		return nil, err
	}
	var res []model.Supplier
	for _, supplier := range suppliers {
		res = append(res, model.Supplier{
			ID:                 supplier.ID,
			Name:               supplier.Name,
			OrderAmount:        supplier.OrderAmount,
			DeliveryFee:        supplier.DeliveryFee,
			FreeDeliveryAmount: supplier.FreeDeliveryAmount,
		})
	}
	return res, nil
}
