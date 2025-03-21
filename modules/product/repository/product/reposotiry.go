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
	// ======== products table ========
	productsTbl = "products"

	pIdCol               = "id"
	pNameCol             = "name"
	pImageUrlCol         = "image_url"
	pGTINCol             = "gtin"
	pLowestSupplierIDCol = "lowest_supplier_id"
	pCreatedAtCol        = "created_at"
	pUpdatedAtCol        = "updated_at"

	// ======== product-supplier table ========
	productsSupplierTbl = "products_supplier"
	psProductIDCol      = "product_id"
	psSupplierIDCol     = "supplier_id"
	psPriceCol          = "price"
	psSellAmountCol     = "sell_amount"

	// ======== supplier table ========
	supplierTbl             = "suppliers"
	sIDCol                  = "user_id"
	sNameCol                = "name"
	sOrderAmountCol         = "order_amount"
	sDeliveryConditionIDCol = "condition_id"

	// ======== delivery conditions ========
	deliveryConditionTbl    = "delivery_conditions"
	dcIDCol                 = "condition_id"
	dcFreeDeliveryAmountCol = "minimum_free_delivery_amount"
	dcDeliveryFeeCol        = "delivery_fee"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) *repo {
	return &repo{db: db}
}

func (r *repo) GetProduct(ctx context.Context, id int64) (*model.Product, error) {
	builder := sq.
		Select(
			pIdCol,
			pNameCol,
			pImageUrlCol,
			pGTINCol,
			pCreatedAtCol,
			pUpdatedAtCol,
		).
		From(productsTbl).
		Where(sq.Eq{pIdCol: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "product_repository.GetProduct",
		QueryRaw: query,
	}

	var product repoModel.Product
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&product.ID,
		&product.Name,
		&product.ImageUrl,
		&product.GTIN,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return converter.ToProductFromRepo(product), nil
}

func (r *repo) GetSupplierProductListByProduct(ctx context.Context, id int64) ([]model.ProductSupplier, error) {
	// Build the query to fetch suppliers for a specific product
	builder := sq.
		Select(
			"ps."+psSupplierIDCol+" AS supplier_id",
			"ps."+psPriceCol+" AS price",
			"ps."+psSellAmountCol+" AS sell_amount",

			"s."+sNameCol+" AS supplier_name",
			"s."+sOrderAmountCol+" AS order_amount",

			"dc."+dcFreeDeliveryAmountCol+" AS minimum_free_delivery_amount",
			"dc."+dcDeliveryFeeCol+" AS delivery_fee",
		).
		From(productsSupplierTbl + " AS ps").
		// Inner join to get only valid suppliers
		Join(supplierTbl + " AS s ON s." + sIDCol + " = ps." + psSupplierIDCol).
		// Optional: Left join delivery conditions
		LeftJoin(deliveryConditionTbl + " AS dc ON dc." + dcIDCol + " = s." + sDeliveryConditionIDCol).
		Where(sq.Eq{psProductIDCol: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := db.Query{
		Name:     "product_repository.GetSupplierProductListByProduct",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []model.ProductSupplier
	for rows.Next() {
		var ps model.ProductSupplier
		var s model.Supplier
		if err := rows.Scan(
			&s.ID,
			&ps.Price,
			&ps.SellAmount,
			&s.Name,
			&s.OrderAmount,
			&s.FreeDeliveryAmount,
			&s.DeliveryFee,
		); err != nil {
			return nil, err
		}
		ps.Supplier = s
		results = append(results, ps)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return results, nil
}

// -----------------------------------------------------------------------------
// GetProductList
// -----------------------------------------------------------------------------
//
// Retrieves a list of products with their "lowest supplier" info.
func (r *repo) GetProductList(ctx context.Context, queryParam *model.ProductListQuery) ([]model.Product, error) {

	builder := sq.Select(
		// Products (p)
		"p."+pIdCol+" AS product_id",
		"p."+pNameCol+" AS product_name",
		"p."+pImageUrlCol+" AS product_image_url",
		"p."+pGTINCol+" AS product_gtin",
		// We'll omit `lowest_supplier` if you don't need to scan/store it:
		// "p."+pLowestSupplierIDCol+" AS product_lowest_supplier",

		// products_supplier (ps)
		"ps."+psPriceCol+" AS ps_price",
		"ps."+psSellAmountCol+" AS ps_sell_amount",

		// supplier (s)
		"s."+sNameCol+" AS supplier_name",
		"s."+sOrderAmountCol+" AS supplier_order_amount",

		// delivery_conditions (dc)
		"dc."+dcFreeDeliveryAmountCol+" AS dc_min_free_delivery_amount",
		"dc."+dcDeliveryFeeCol+" AS dc_delivery_fee",
	).
		From(productsTbl + " AS p").
		// Join products_supplier using the known `lowest_supplier`
		LeftJoin(productsSupplierTbl + " AS ps ON ps." + psProductIDCol + " = p." + pIdCol +
			" AND ps." + psSupplierIDCol + " = p." + pLowestSupplierIDCol).
		// Join supplier on ps.supplier_id
		LeftJoin(supplierTbl + " AS s ON s." + sIDCol + " = ps." + psSupplierIDCol).
		// Join delivery_conditions on s.condition_id
		LeftJoin(deliveryConditionTbl + " AS dc ON dc." + dcIDCol + " = s." + sDeliveryConditionIDCol).
		PlaceholderFormat(sq.Dollar)

	// Optional limit & offset
	if queryParam.Limit > 0 {
		builder = builder.Limit(uint64(queryParam.Limit))
	} else {
		builder = builder.Limit(30)
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

	// We'll scan into these repoModel structs:
	var productList []repoModel.Product

	for rows.Next() {
		var (
			p  repoModel.Product
			ps repoModel.ProductSupplier
			s  repoModel.Supplier
		)

		// The order and count of Scan fields must match SELECT columns exactly.
		err := rows.Scan(
			&p.ID,       // product_id
			&p.Name,     // product_name
			&p.ImageUrl, // product_image_url
			&p.GTIN,     // product_gtin

			&ps.Price,      // ps_price
			&ps.SellAmount, // ps_sell_amount

			&s.Name,        // supplier_name
			&s.OrderAmount, // supplier_order_amount

			&s.FreeDeliveryAmount, // dc_min_free_delivery_amount
			&s.DeliveryFee,        // dc_delivery_fee
		)
		if err != nil {
			return nil, err
		}

		ps.Supplier = s
		p.LowestSupplier = ps

		productList = append(productList, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return converter.ToProductListFromRepo(productList), nil
}

func (r *repo) GetTotalProducts(ctx context.Context) (int, error) {
	builder := sq.
		Select("COUNT(*)").
		From(productsTbl).
		PlaceholderFormat(sq.Dollar)

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
