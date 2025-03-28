package repository

import (
	"context"
	"database/sql"
	"diploma/modules/cart/model"
	"diploma/pkg/client/db"
	"errors"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

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

	if err != nil {
		if errors.As(sql.ErrNoRows, &err) {
			return 0, model.ErrNoRows
		}
		return 0, err
	}
	return quantity, nil
}

func (r *cartRepo) UpdateItemQuantity(ctx context.Context, cartID, productId, supplierId int64, quantity int) error {
	builder := sq.Update(cartItemTable).
		PlaceholderFormat(sq.Dollar).
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
		Values(input.CartID, input.SupplierID, input.ProductID, input.Price, input.Quantity)

	query, args, err := builder.ToSql()

	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "cart_repository.AddItem",
		QueryRaw: query,
	}
	res, err := r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return fmt.Errorf("no rows inserted")
	}

	return nil
}
