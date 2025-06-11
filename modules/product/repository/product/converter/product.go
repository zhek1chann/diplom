package converter

import (
	"diploma/modules/product/model"
	repoModel "diploma/modules/product/repository/product/model"
)

func toProductFromRepo(p repoModel.Product) model.Product {
	return model.Product{
		ID:                    p.ID,
		Name:                  p.Name,
		ImageUrl:              p.ImageUrl,
		GTIN:                  p.GTIN,
		CreatedAt:             p.CreatedAt,
		UpdatedAt:             p.UpdatedAt,
		LowestProductSupplier: ToProductSupplierFromRepo(p.LowestSupplier),
	}
}

func ToProductSupplierFromRepo(sp repoModel.ProductSupplier) model.ProductSupplier {
	return model.ProductSupplier{
		Price:      sp.Price,
		SellAmount: sp.SellAmount,
		Supplier:   model.Supplier(sp.Supplier),
	}
}
