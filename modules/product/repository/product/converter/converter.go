package converter

import (
	"diploma/modules/product/model"
	repoModel "diploma/modules/product/repository/product/model"
)

func ToProductFromRepo(product repoModel.Product) *model.Product {
	var categoryName, subcategoryName string
	if product.CategoryName != nil {
		categoryName = *product.CategoryName
	}
	if product.SubcategoryName != nil {
		subcategoryName = *product.SubcategoryName
	}

	var price, sellAmount int
	if product.LowestSupplier.Price != nil {
		price = *product.LowestSupplier.Price
	}
	if product.LowestSupplier.SellAmount != nil {
		sellAmount = *product.LowestSupplier.SellAmount
	}

	return &model.Product{
		ID:              product.ID,
		GTIN:            product.GTIN,
		Name:            product.Name,
		ImageUrl:        product.ImageUrl,
		CreatedAt:       product.CreatedAt,
		UpdatedAt:       product.UpdatedAt,
		CategoryID:      product.CategoryID,
		SubcategoryID:   product.SubcategoryID,
		CategoryName:    categoryName,
		SubcategoryName: subcategoryName,
		LowestProductSupplier: model.ProductSupplier{
			Price:      &price,
			SellAmount: &sellAmount,
			Supplier:   ToSupplierFromRepo(product.LowestSupplier.Supplier),
		},
	}
}

func ToSupplierFromRepo(supplier repoModel.Supplier) model.Supplier {
	var name string
	var orderAmount, freeDeliveryAmount, deliveryFee int

	if supplier.Name != nil {
		name = *supplier.Name
	}
	if supplier.OrderAmount != nil {
		orderAmount = *supplier.OrderAmount
	}
	if supplier.FreeDeliveryAmount != nil {
		freeDeliveryAmount = *supplier.FreeDeliveryAmount
	}
	if supplier.DeliveryFee != nil {
		deliveryFee = *supplier.DeliveryFee
	}

	return model.Supplier{
		ID:                 supplier.ID,
		Name:               &name,
		OrderAmount:        &orderAmount,
		FreeDeliveryAmount: &freeDeliveryAmount,
		DeliveryFee:        &deliveryFee,
	}
}

func ToProductListFromRepo(products []repoModel.Product) []model.Product {
	result := make([]model.Product, len(products))
	for i, product := range products {
		result[i] = *ToProductFromRepo(product)
	}
	return result
}
