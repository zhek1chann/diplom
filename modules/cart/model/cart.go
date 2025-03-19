package model

type PutCartQuery struct {
	Quantity   int
	Price      int
	ProductID  int64
	SupplierID int64
	UserID     int64
}

type Cart struct {
	Total     int
	UserID    int64
	Suppliers []Supplier
}

type Supplier struct {
	MinOrderAmount     int
	TotalAmount        int
	FreeDeliveryAmount int
	DeliveryFee        int
	ID                 int64
	Name               string
	ProductList        []Product
}

type Product struct {
	Price    int
	Quantity int
	ID       int64
	Name     string
	Image    string
}

type DeleteProductQuery struct {
	UserID     int64
	ProductID  int64
	Quantity   int
	SupplierID int64
}
