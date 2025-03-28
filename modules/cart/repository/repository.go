package repository

import "diploma/pkg/client/db"

const (
	cartsTable = "carts"

	cIDColumn         = "id"
	cTotalColumn      = "total"
	cCustomerIDColumn = "customer_id"
	cCreatedAtColumn  = "created_at"
	cUpdatedAtColumn  = "updated_at"

	//============= cart_items table columns =============
	cartItemTable      = "cart_items"
	ciID               = "id"
	ciCartIDColumn     = "cart_id"
	ciSupplierIDColumn = "supplier_id"
	ciProductIDColumn  = "product_id"
	ciQuantityColumn   = "quantity"
	ciPriceColumn      = "price"

	//============= cart_items_suppliers table columns =============
	productsTable  = "products"
	imageURLColumn = "image_url"
	nameColumn     = "name"
)

type cartRepo struct {
	db db.Client
}

func NewRepository(db db.Client) *cartRepo {
	return &cartRepo{db: db}
}
