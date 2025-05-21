package address

import (
	"context"
	"diploma/pkg/client/db"
	"time"

	sq "github.com/Masterminds/squirrel"
	model "diploma/modules/user/model"
)

const (
	tableName         = "address"
	idColumn          = "id"
	streetColumn      = "street"
	descriptionColumn = "description"
	userIDColumn      = "user_id"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) *repo {
	return &repo{db: db}
}

func (r *repo) CreateAddress(ctx context.Context, addr model.Address) (int64, error) {
	builder := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(streetColumn, descriptionColumn, userIDColumn).
		Values(addr.Street, addr.Description, addr.UserID).
		Suffix("RETURNING " + idColumn)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "address_repository.Create",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) GetByUserId(ctx context.Context, userID int64) ([]model.Address, error) {
	builder := sq.Select(idColumn, streetColumn, descriptionColumn, userIDColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{userIDColumn: userID})
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "address_repository.GetByUserId",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []model.Address
	for rows.Next() {
		var addr model.Address
		if err = rows.Scan(&addr.ID, &addr.Street, &addr.Description, &addr.UserID); err != nil {
			return nil, err
		}
		addresses = append(addresses, addr)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return addresses, nil
}

func (r *repo) GetById(ctx context.Context, id int64) (*model.Address, error) {
	builder := sq.Select(idColumn, streetColumn, descriptionColumn, userIDColumn).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "address_repository.GetById",
		QueryRaw: query,
	}

	var addr model.Address
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&addr.ID, &addr.Street, &addr.Description, &addr.UserID)
	if err != nil {
		return nil, err
	}

	return &addr, nil
}

func (r *repo) Update(ctx context.Context, id int64, street, description string) error {
	builder := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		Set(streetColumn, street).
		Set(descriptionColumn, description).
		Set("updated_at", time.Now()) // Optional: remove if table doesn't have updated_at column

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "address_repository.Update",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}

func (r *repo) Delete(ctx context.Context, id int64) error {
	builder := sq.Delete(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := db.Query{
		Name:     "address_repository.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	return err
}
