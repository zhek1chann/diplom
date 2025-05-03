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

func (r *cartRepo) Cart(ctx context.Context, userID int64) (*model.Cart, error) {
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

func (r *cartRepo) UpdateCartTotal(ctx context.Context, cartID int64, total int) error {
	builder := sq.Update(cartsTable).
		PlaceholderFormat(sq.Dollar).
		Set(cTotalColumn, total).
		Where(sq.Eq{cIDColumn: cartID})

	query, args, err := builder.ToSql()

	if err != nil {
		return err
	}

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

func (r *cartRepo) GetCartItems(ctx context.Context, cartID int64) ([]model.Supplier, error) {
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
		if errors.As(sql.ErrNoRows, &err) {
			return nil, model.ErrNoRows
		}
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

	return converter.ToServiceSupplierFromRepo(items), nil
}

func (r *cartRepo) DeleteCart(ctx context.Context, cartID int64) error {
	builder := sq.Delete(cartsTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{cIDColumn: cartID})

	query, args, err := builder.ToSql()

	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "cart_repository.Delete",
		QueryRaw: query,
	}

	result, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("error executing query: %v", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("no rows deleted")
	}

	return nil
}