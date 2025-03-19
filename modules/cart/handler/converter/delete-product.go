package converter

import (
	modelApi "diploma/modules/cart/handler/model"
	"diploma/modules/cart/model"
)

func ToServiceDeleleProductFromApi(input *modelApi.DeleteProductFromCartInput) *model.DeleteProductQuery {
	return &model.DeleteProductQuery{
		UserID:     input.UserID,
		ProductID:  input.ProductID,
		SupplierID: input.SupplierID,
		Quantity:   input.Quantity,
	}
}
