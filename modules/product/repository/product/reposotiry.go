package product

import (
	"context"
	"database/sql"
	"diploma/modules/product/model"
	"diploma/modules/product/repository/product/converter"
	repoModel "diploma/modules/product/repository/product/model"
	"diploma/pkg/client/db"
	"fmt"

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

// GetProduct retrieves a product by its ID
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
		return nil, fmt.Errorf("failed to build get product query: %w", err)
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
		if err == sql.ErrNoRows {
			return nil, model.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return converter.ToProductFromRepo(product), nil
}

// GetProductListByIDList retrieves a list of products by their IDs
func (r *repo) GetProductListByIDList(ctx context.Context, idList []int64) ([]*model.Product, error) {
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
		Where(sq.Eq{pIdCol: idList}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get product list query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetProductListByIDList",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get product list: %w", err)
	}
	defer rows.Close()

	var products []repoModel.Product
	for rows.Next() {
		var product repoModel.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.ImageUrl,
			&product.GTIN,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	var result []*model.Product
	for _, product := range products {
		result = append(result, converter.ToProductFromRepo(product))
	}

	return result, nil
}

// GetSupplierProductListByProduct retrieves a list of suppliers for a specific product
func (r *repo) GetSupplierProductListByProduct(ctx context.Context, id int64) ([]model.ProductSupplier, error) {
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
		Join(supplierTbl + " AS s ON s." + sIDCol + " = ps." + psSupplierIDCol).
		LeftJoin(deliveryConditionTbl + " AS dc ON dc." + dcIDCol + " = s." + sDeliveryConditionIDCol).
		Where(sq.Eq{psProductIDCol: id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get supplier product list query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetSupplierProductListByProduct",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get supplier product list: %w", err)
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
			return nil, fmt.Errorf("failed to scan supplier product: %w", err)
		}
		ps.Supplier = s
		results = append(results, ps)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating supplier products: %w", err)
	}

	return results, nil
}

// GetProductList retrieves a list of products with their lowest supplier info
func (r *repo) GetProductList(ctx context.Context, queryParam *model.ProductListQuery) ([]model.Product, error) {
	builder := sq.Select(
		"p."+pIdCol+" AS product_id",
		"p."+pNameCol+" AS product_name",
		"p."+pImageUrlCol+" AS product_image_url",
		"p."+pGTINCol+" AS product_gtin",
		"ps."+psPriceCol+" AS ps_price",
		"ps."+psSellAmountCol+" AS ps_sell_amount",
		"s."+sNameCol+" AS supplier_name",
		"s."+sOrderAmountCol+" AS supplier_order_amount",
		"dc."+dcFreeDeliveryAmountCol+" AS dc_min_free_delivery_amount",
		"dc."+dcDeliveryFeeCol+" AS dc_delivery_fee",
	).
		From(productsTbl + " AS p").
		LeftJoin(productsSupplierTbl + " AS ps ON ps." + psProductIDCol + " = p." + pIdCol +
			" AND ps." + psSupplierIDCol + " = p." + pLowestSupplierIDCol).
		LeftJoin(supplierTbl + " AS s ON s." + sIDCol + " = ps." + psSupplierIDCol).
		LeftJoin(deliveryConditionTbl + " AS dc ON dc." + dcIDCol + " = s." + sDeliveryConditionIDCol).
		PlaceholderFormat(sq.Dollar)

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
		return nil, fmt.Errorf("failed to build get product list query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetProductList",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get product list: %w", err)
	}
	defer rows.Close()

	var productList []repoModel.Product
	for rows.Next() {
		var (
			p  repoModel.Product
			ps repoModel.ProductSupplier
			s  repoModel.Supplier
		)

		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.ImageUrl,
			&p.GTIN,
			&ps.Price,
			&ps.SellAmount,
			&s.Name,
			&s.OrderAmount,
			&s.FreeDeliveryAmount,
			&s.DeliveryFee,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}

		ps.Supplier = s
		p.LowestSupplier = ps
		productList = append(productList, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return converter.ToProductListFromRepo(productList), nil
}

// GetTotalProducts returns the total number of products
func (r *repo) GetTotalProducts(ctx context.Context) (int, error) {
	builder := sq.
		Select("COUNT(*)").
		From(productsTbl).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build get total products query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetTotalProducts",
		QueryRaw: query,
	}

	var total int
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total products: %w", err)
	}

	return total, nil
}

// GetProductPriceBySupplier retrieves the price of a product from a specific supplier
func (r *repo) GetProductPriceBySupplier(ctx context.Context, productID, supplierID int64) (int, error) {
	builder := sq.
		Select(psPriceCol).
		From(productsSupplierTbl).
		Where(sq.And{
			sq.Eq{psProductIDCol: productID},
			sq.Eq{psSupplierIDCol: supplierID},
		}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build get product price query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetProductPriceBySupplier",
		QueryRaw: query,
	}

	var price int
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&price)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, model.ErrNoRows
		}
		return 0, fmt.Errorf("failed to get product price: %w", err)
	}
	return price, nil
}
