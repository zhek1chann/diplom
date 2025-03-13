package product

import (
	"context"
	"diploma/modules/product/model"
	"diploma/modules/product/repository/product/converter"
	repoModel "diploma/modules/product/repository/product/model"

	"diploma/pkg/client/db"

	sq "github.com/Masterminds/squirrel"
)

const (
	productsTable       = "products"
	productsSupplierTbl = "products_supplier"

	idColumn        = "id"
	nameColumn      = "name"
	imageUrlColumn  = "image_url"
	GTINColumn      = "gtin"
	minPriceColumn  = "min_price"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"

	// products_supplier table columns
	psProductIDColumn   = "product_id"
	psSupplierIDColumn  = "supplier_id"
	priceColumn         = "price"
	minSellAmountColumn = "min_sell_amount"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) *repo {
	return &repo{db: db}
}

// GetProduct retrieves a product by its id, then converts the repository model into the domain model.
func (r *repo) GetProduct(ctx context.Context, id int64) (*model.Product, error) {
	builder := sq.Select(idColumn, nameColumn, imageUrlColumn, GTINColumn, minPriceColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(productsTable).
		Where(sq.Eq{idColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "product_repository.GetProduct",
		QueryRaw: query,
	}

	var repoProd repoModel.Product
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&repoProd.ID, &repoProd.Name, &repoProd.ImageUrl, &repoProd.GTIN, &repoProd.MinPrice, &repoProd.CreatedAt, &repoProd.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return converter.ToProductFromRepo(&repoProd), nil
}

// GetSupplierInfoListByProduct retrieves supplier information for a given product,
// converts each repository-level supplier info to the domain model.
func (r *repo) GetSupplierInfoListByProduct(ctx context.Context, id int64) ([]model.ProductSupplierInfo, error) {
	builder := sq.Select(psSupplierIDColumn, priceColumn, minSellAmountColumn).
		PlaceholderFormat(sq.Dollar).
		From(productsSupplierTbl).
		Where(sq.Eq{psProductIDColumn: id})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "product_repository.GetSupplierInfoListByProduct",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domainInfos []model.ProductSupplierInfo
	for rows.Next() {
		var repoInfo repoModel.ProductSupplierInfo
		// Assuming repoModel.ProductSupplierInfo has fields: SupplierID, Price, and MinSellAmount.
		if err := rows.Scan(&repoInfo.SupplierID, &repoInfo.Price, &repoInfo.MinSellAmount); err != nil {
			return nil, err
		}
		// Convert each repository supplier info to the domain supplier info.
		domainInfo := converter.ToProductSupplierInfoFromRepo(&repoInfo)
		domainInfos = append(domainInfos, domainInfo)
	}

	return domainInfos, nil
}

// GetProductList retrieves a list of products based on the provided query parameters,
// then converts each repository-level product into the domain model.
func (r *repo) GetProductList(ctx context.Context, queryParam *model.ProductListQuery) ([]model.Product, error) {
	builder := sq.Select(idColumn, nameColumn, imageUrlColumn, GTINColumn, minPriceColumn, createdAtColumn, updatedAtColumn).
		PlaceholderFormat(sq.Dollar).
		From(productsTable)

	if queryParam.Limit > 0 {
		builder = builder.Limit(uint64(queryParam.Limit))
	}
	if queryParam.Offset > 0 {
		builder = builder.Offset(uint64(queryParam.Offset))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "product_repository.GetProductList",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var domainProducts []model.Product
	for rows.Next() {
		var repoProd repoModel.Product
		if err := rows.Scan(
			&repoProd.ID, &repoProd.Name, &repoProd.ImageUrl, &repoProd.GTIN, &repoProd.MinPrice, &repoProd.CreatedAt, &repoProd.UpdatedAt,
		); err != nil {
			return nil, err
		}
		domainProducts = append(domainProducts, *converter.ToProductFromRepo(&repoProd))
	}

	return domainProducts, nil
}

// GetTotalProducts returns the total number of products in the products table.
func (r *repo) GetTotalProducts(ctx context.Context) (int, error) {
	builder := sq.Select("COUNT(*)").
		PlaceholderFormat(sq.Dollar).
		From(productsTable)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := db.Query{
		Name:     "product_repository.GetTotalProducts",
		QueryRaw: query,
	}

	var total int
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
