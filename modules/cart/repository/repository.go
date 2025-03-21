package repository

import (
	"context"
	"diploma/modules/cart/model"
	"diploma/modules/cart/repository/converter"
	modelRepo "diploma/modules/cart/repository/model"
	"diploma/pkg/client/db"

	sq "github.com/Masterminds/squirrel"
)

const (
	cartsTable = "carts"

	idColumn         = "id"
	totalColumn      = "total"
	customerIDColumn = "customer_id"
	createdAtColumn  = "created_at"
	updatedAtColumn  = "updated_at"

	//============= cart_items table columns =============
	cartItemTable = "cart_items"

	cartIDColumn     = "cart_id"
	supplierIDColumn = "supplier_id"
	productIDColumn  = "product_id"
	quantityColumn   = "quantity"
	priceColumn      = "price"

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

func (r *cartRepo) CreateCart(ctx context.Context, userID int64) (int64, error) {
	builder := sq.Insert(cartsTable).
		PlaceholderFormat(sq.Dollar).
		Columns(customerIDColumn).
		Values(userID).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()

	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "cart_repository.Create",
		QueryRaw: query,
	}

	var cartID int64

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&cartID)
	if err != nil {
		return 0, err
	}
	return cartID, nil
}

func (r *cartRepo) AddItem(ctx context.Context, input *model.PutCartQuery) error {
	builder := sq.Insert(cartItemTable).
		PlaceholderFormat(sq.Dollar).
		Columns(cartIDColumn, supplierIDColumn, productIDColumn, quantityColumn, priceColumn).
		Values(input.CartID, input.SupplierID, input.ProductID, input.Quantity, input.Price).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()

	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "cart_repository.AddItem",
		QueryRaw: query,
	}

	var cartItemID int64

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&cartItemID)
	if err != nil {
		return err
	}

	return nil
}

func (r *cartRepo) UpdateCartTotal(ctx context.Context, cartID int64, total int) error {
	builder := sq.Update(cartsTable).
		Set(totalColumn, total).
		Where(sq.Eq{idColumn: cartID}).
		Suffix("RETURNING id")

	query, args, err := builder.ToSql()

	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "cart_repository.UpdateTotal",
		QueryRaw: query,
	}

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&cartID)
	if err != nil {
		return err
	}

	return nil
}

func (r *cartRepo) GetCart(ctx context.Context, userID int64) (*model.Cart, error) {
	builder := sq.Select(idColumn, totalColumn).
		PlaceholderFormat(sq.Dollar).
		From(cartsTable).
		Where(sq.Eq{cartIDColumn: userID})

	query, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "cart_repository.Get",
		QueryRaw: query,
	}

	var cart modelRepo.Cart

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&cart.ID, &cart.Total)
	if err != nil {
		return nil, err
	}
	return converter.ToServiceCartFromRepo(&cart), nil
}

func (r *cartRepo) GetCartItems(ctx context.Context, cartID int64) ([]modelRepo.CartItem, error) {
	// Query to select cart items and product details
	builder := sq.Select(
		cartItemTable+"."+productIDColumn,
		cartItemTable+"."+supplierIDColumn,
		cartItemTable+"."+quantityColumn,
		cartItemTable+"."+priceColumn,
		productsTable+"."+nameColumn+" AS product_name",
		productsTable+"."+imageURLColumn,
	).
		PlaceholderFormat(sq.Dollar).
		From(cartItemTable).
		Join(productsTable + " ON " + cartItemTable + "." + productIDColumn + "=" + productsTable + "." + "id").
		Where(sq.Eq{cartIDColumn: cartID})

	query, args, err := builder.ToSql()

	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "cart_repository.GetCartItems",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []modelRepo.CartItem
	for rows.Next() {
		var item modelRepo.CartItem

		if err := rows.Scan(&item.ProductID, &item.SupplierID, &item.Quantity, &item.Price, &item.ProductName, &item.ProductImageURL); err != nil {
			return nil, err
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}
