package repository

import (
	"context"
	"database/sql"
	"diploma/modules/cart/model"
	"diploma/modules/cart/repository/converter"
	modelRepo "diploma/modules/cart/repository/model"
	"diploma/pkg/client/db"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

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

func (r *cartRepo) CreateCart(ctx context.Context, userID int64) (int64, error) {
	builder := sq.Insert(cartsTable).
		PlaceholderFormat(sq.Dollar).
		Columns(cCustomerIDColumn).
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

func (r *cartRepo) ItemQuantity(ctx context.Context, cartID, productId, supplierId int64) (int, error) {
	builder := sq.Select(ciQuantityColumn).
		PlaceholderFormat(sq.Dollar).
		From(cartItemTable).
		Where(sq.And{
			sq.Eq{ciCartIDColumn: cartID},
			sq.Eq{ciSupplierIDColumn: supplierId},
			sq.Eq{ciProductIDColumn: productId},
		})

	query, args, err := builder.ToSql()

	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "cart_repository.ItemQuantity",
		QueryRaw: query,
	}

	var quantity int

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&quantity)
	fmt.Println(err)
	if err == sql.ErrNoRows {
		return 0, nil
	} else if err != nil {
		return 0, err
	}
	return quantity, nil
}

func (r *cartRepo) UpdateItemQuantity(ctx context.Context, cartID, productId, supplierId int64, quantity int) error {
	builder := sq.Update(cartItemTable).
		Set(ciQuantityColumn, quantity).
		Where(sq.And{
			sq.Eq{ciCartIDColumn: cartID},
			sq.Eq{ciSupplierIDColumn: supplierId},
			sq.Eq{ciProductIDColumn: productId},
		})

	query, args, err := builder.ToSql()

	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "cart_repository.UpdateItemQuantity",
		QueryRaw: query,
	}

	res, err := r.db.DB().ExecContext(ctx, q, args...)
	fmt.Println(err)
	if err != nil {
		return err
	}

	if res.RowsAffected() == 0 {
		return fmt.Errorf("not updated item")
	}
	return nil
}

func (r *cartRepo) AddItem(ctx context.Context, input *model.PutCartQuery) error {
	builder := sq.Insert(cartItemTable).
		PlaceholderFormat(sq.Dollar).
		Columns(ciCartIDColumn, ciSupplierIDColumn, ciProductIDColumn, ciPriceColumn, ciQuantityColumn).
		Values(input.CartID, input.SupplierID, input.ProductID, input.Price, input.Quantity).
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
		Set(cTotalColumn, total).
		Where(sq.Eq{cIDColumn: cartID})

	query, args, err := builder.ToSql()

	if err != nil {
		return err
	}
	fmt.Printf("Generated SQL Query: %s\n", query)
	fmt.Printf("SQL Arguments: %v\n", args)

	q := db.Query{
		Name:     "cart_repository.UpdateTotal",
		QueryRaw: query,
	}

	result, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no rows updated")
	}

	return nil
}

func (r *cartRepo) GetCart(ctx context.Context, userID int64) (*model.Cart, error) {
	builder := sq.Select(cIDColumn, cTotalColumn).
		PlaceholderFormat(sq.Dollar).
		From(cartsTable).
		Where(sq.Eq{cCustomerIDColumn: userID})

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
		if errors.As(sql.ErrNoRows, &err) {
			return nil, model.ErrNoRows
		}
		return nil, err
	}
	return converter.ToServiceCartFromRepo(&cart), nil
}

func (r *cartRepo) GetCartItems(ctx context.Context, cartID int64) ([]modelRepo.CartItem, error) {
	// Query to select cart items and product details
	builder := sq.Select(
		cartItemTable+"."+ciProductIDColumn,
		cartItemTable+"."+ciSupplierIDColumn,
		cartItemTable+"."+ciQuantityColumn,
		cartItemTable+"."+ciPriceColumn,
		productsTable+"."+nameColumn+" AS product_name",
		productsTable+"."+imageURLColumn,
	).
		PlaceholderFormat(sq.Dollar).
		From(cartItemTable).
		Join(productsTable + " ON " + cartItemTable + "." + ciProductIDColumn + "=" + productsTable + "." + "id").
		Where(sq.Eq{ciCartIDColumn: cartID})

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
