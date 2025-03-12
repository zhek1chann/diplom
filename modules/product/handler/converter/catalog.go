package converter

import (
	modelApi "diploma/modules/product/handler/model"
	"diploma/modules/product/model"
)

func ToServieProductQueryFromApi(input modelApi.ProductInput) *model.ProductQuery {
	return &model.ProductQuery{
		ID: input.ID,
	}
}

func ToApiProductResponeFromService(product *model.DetailedProduct) *modelApi.ProductResponse {
	return &modelApi.ProductResponse{
		DetailedProduct: ToApiDetailedProductFromSerivce(product),
	}
}

func ToApiDetailedProductFromSerivce(dProudct *model.DetailedProduct) *modelApi.DetailedProduct {
	return &modelApi.DetailedProduct{
		Product:      ToAPIProductFromService(dProudct.Product),
		SupplierList: ToApiSupplierListFromService(dProudct.SupplierList),
	}
}

func ToApiSupplierListFromService(supplierList []model.ProductSupplierInfo) []modelApi.ProductSupplierInfo {
	res := make([]modelApi.ProductSupplierInfo, 0, len(supplierList))

	for _, s := range supplierList {
		res = append(res, ToAPISupplierInfoFromService(s))
	}

	return res
}

func ToServiceProductFromApi(apiProduct *modelApi.Product) *model.Product {
	return &model.Product{
		ID:           apiProduct.ID,
		Name:         apiProduct.Name,
		MinPrice:     apiProduct.MinPrice,
		ImageURL:     apiProduct.ImageURL,
		GTIN:         apiProduct.GTIN,
		SupplierInfo: ToServiceSupplierInfoFromAPI(apiProduct.SupplierInfo),
	}
}

func ToAPIProductFromService(serviceProduct *model.Product) *modelApi.Product {
	return &modelApi.Product{
		ID:           serviceProduct.ID,
		Name:         serviceProduct.Name,
		MinPrice:     serviceProduct.MinPrice,
		ImageURL:     serviceProduct.ImageURL,
		GTIN:         serviceProduct.GTIN,
		SupplierInfo: ToAPISupplierInfoFromService(serviceProduct.SupplierInfo),
	}
}

// ConvertAPIToServiceSupplierInfo преобразует информацию о поставщике из API в сервис.
func ToServiceSupplierInfoFromAPI(apiSupplierInfo modelApi.ProductSupplierInfo) model.ProductSupplierInfo {
	return model.ProductSupplierInfo{
		SupplierID:                apiSupplierInfo.SupplierID,
		Name:                      apiSupplierInfo.Name,
		MinimumFreeDeliveryAmount: apiSupplierInfo.MinimumFreeDeliveryAmount,
		DeliveryFee:               apiSupplierInfo.DeliveryFee,
	}
}

// ConvertServiceToAPISuppSlierInfo преобразует информацию о поставщике из сервиса в API.
func ToAPISupplierInfoFromService(serviceSupplierInfo model.ProductSupplierInfo) modelApi.ProductSupplierInfo {
	return modelApi.ProductSupplierInfo{
		SupplierID:                serviceSupplierInfo.SupplierID,
		Name:                      serviceSupplierInfo.Name,
		MinimumFreeDeliveryAmount: serviceSupplierInfo.MinimumFreeDeliveryAmount,
		DeliveryFee:               serviceSupplierInfo.DeliveryFee,
	}
}
