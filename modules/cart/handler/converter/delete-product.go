package converter

import (
	modelApi "diploma/modules/cart/handler/model"
	"diploma/modules/cart/model"
)

func ToServiceDeleleProductFromApi(input *modelApi.DeleteProductFromCartInput) *model.DeleteProductQuery {
	return &model.DeleteProductQuery{
		CustomerID: input.CustomerID,
		ProductID:  input.ProductID,
		SupplierID: input.SupplierID,
		Quantity:   input.Quantity,
	}
}
