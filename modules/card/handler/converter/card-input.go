package converter

import (
	modelApi "diploma/modules/card/handler/model"
	"diploma/modules/card/model"
)

func ToServiceCardInputFromAPI(input *modelApi.CardInput) *model.CardInput {
	return &model.CardInput{
		UserID:     input.UserID,
		ProductID:  input.ProductID,
		Quantity:   input.Quantity,
		Price:      input.Price,
		SupplierID: input.SupplierID,
	}
}
