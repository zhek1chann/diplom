package model

type ProductInput struct {
	ID int64 `json:"id"`
}

type ProductResponse struct {
	DetailedProduct *DetailedProduct `json:"product"`
}

type Product struct {
	ID                    int64           `json:"id"`
	Name                  string          `json:"name"`
	ImageUrl              string          `json:"image`
	LowestProductSupplier ProductSupplier `json:"lowest_product_supplier"`
}

type ProductSupplier struct {
	Price      int      `json:"price"`
	SellAmount int      `json:"sell_amount"`
	Supplier   Supplier `json:"supplier"`
}

type Supplier struct {
	ID                 int64  `json:"id"`
	Name               string `json:"name"`
	OrderAmount        int    `json:"order_amount"`
	FreeDeliveryAmount int    `json:"free_delivery_amount"`
	DeliveryFee        int    `json:"delivery_fee"`
}

type DetailedProduct struct {
	*Product            `json:"product"`
	ProductSupplierList []ProductSupplier `json:"suppliers"`
}
