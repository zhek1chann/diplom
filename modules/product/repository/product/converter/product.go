package converter

import (
	"diploma/modules/product/model"
	repoModel "diploma/modules/product/repository/product/model"
)

// ToProductFromRepo converts a repository-level Product to a domain-level Product.
func ToProductFromRepo(p *repoModel.Product) *model.Product {
	return &model.Product{
		ID:           p.ID,
		Name:         p.Name,
		MinPrice:     p.MinPrice,
		ImageURL:     p.ImageUrl,
		GTIN:         p.GTIN,
		SupplierInfo: ToProductSupplierInfoFromRepo(&p.SupplierInfo),
	}
}

// ToProductSupplierInfoFromRepo converts a repository-level ProductSupplierInfo
// to a domain-level ProductSupplierInfo.
func ToProductSupplierInfoFromRepo(info *repoModel.ProductSupplierInfo) model.ProductSupplierInfo {
	return model.ProductSupplierInfo{
		SupplierID:                info.SupplierID,
		Name:                      info.SupplierName,
		MinimumFreeDeliveryAmount: info.MinimumFreeDeliveryAmount,
		DeliveryFee:               info.DeliveryFee,
	}
}
