package converter

import (
	modelApi "diploma/modules/product/handler/model"
	"diploma/modules/product/model"
)

func ToServiceProductListQueryFromAPI(input *modelApi.ProductListInput) *model.ProductListQuery {
	return &model.ProductListQuery{
		Limit:         input.Limit,
		Offset:        input.Offset,
		CategoryID:    input.CategoryID,
		SubcategoryID: input.SubcategoryID,
	}
}

func ToProductListResponeFromService(producList *model.ProductList) *modelApi.ProductListResponse {
	return &modelApi.ProductListResponse{
		ProductList: ToProductsFromService(producList.Products),
		Total:       producList.Total,
	}
}

func ToProductsFromService(products []model.Product) []modelApi.Product {
	res := make([]modelApi.Product, 0, len(products))
	for _, e := range products {
		res = append(res, *ToAPIProductFromService(&e))
	}

	return res
}
