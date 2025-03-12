package converter

import (
	modelApi "diploma/modules/product/handler/model"
	"diploma/modules/product/model"
)

func ToServicePageCountFromAPI(input *modelApi.PageCountInput) *model.PageCountQuery {
	return &model.PageCountQuery{
		PageSize: input.PageSize,
	}
}
