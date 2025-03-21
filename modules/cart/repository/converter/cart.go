package converter

import (
	"diploma/modules/cart/model"
	modelRepo "diploma/modules/cart/repository/model"
)

func ToServiceCartFromRepo(input *modelRepo.Cart) *model.Cart {
	return &model.Cart{
		ID:         input.ID,
		CustomerID: input.CustomerID,
		Total:      input.Total,
	}
}
