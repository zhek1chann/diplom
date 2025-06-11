package product

import (
	"context"
	"database/sql"
	"diploma/modules/product/model"
	"diploma/modules/product/repository/product/converter"
	repoModel "diploma/modules/product/repository/product/model"
	"diploma/pkg/client/db"
	"fmt"
	"strings"

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
			"p."+pIdCol,
			"p."+pNameCol,
			"p."+pImageUrlCol,
			"p."+pGTINCol,
			"p."+pCreatedAtCol,
			"p."+pUpdatedAtCol,
			"p.category_id",
			"p.subcategory_id",
			"c.name AS category_name",
			"sc.name AS subcategory_name",
			"ps."+psPriceCol+" AS ps_price",
			"ps."+psSellAmountCol+" AS ps_sell_amount",
			"s."+sNameCol+" AS supplier_name",
			"s."+sOrderAmountCol+" AS supplier_order_amount",
			"dc."+dcFreeDeliveryAmountCol+" AS dc_min_free_delivery_amount",
			"dc."+dcDeliveryFeeCol+" AS dc_delivery_fee",
		).
		From(productsTbl + " AS p").
		LeftJoin("categories AS c ON p.category_id = c.id").
		LeftJoin("subcategories AS sc ON p.subcategory_id = sc.id").
		LeftJoin(productsSupplierTbl + " AS ps ON ps." + psProductIDCol + " = p." + pIdCol +
			" AND ps." + psSupplierIDCol + " = p." + pLowestSupplierIDCol).
		LeftJoin(supplierTbl + " AS s ON s." + sIDCol + " = ps." + psSupplierIDCol).
		LeftJoin(deliveryConditionTbl + " AS dc ON dc." + dcIDCol + " = s." + sDeliveryConditionIDCol).
		Where(sq.Eq{"p." + pIdCol: id}).
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
	var ps repoModel.ProductSupplier
	var s repoModel.Supplier

	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&product.ID,
		&product.Name,
		&product.ImageUrl,
		&product.GTIN,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.CategoryID,
		&product.SubcategoryID,
		&product.CategoryName,
		&product.SubcategoryName,
		&ps.Price,
		&ps.SellAmount,
		&s.Name,
		&s.OrderAmount,
		&s.FreeDeliveryAmount,
		&s.DeliveryFee,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	ps.Supplier = s
	product.LowestSupplier = ps

	return converter.ToProductFromRepo(product), nil
}

// GetProductList retrieves a list of products with their lowest supplier info
func (r *repo) GetProductList(ctx context.Context, queryParam *model.ProductListQuery) ([]model.Product, error) {
	builder := sq.
		Select(
			"p."+pIdCol,
			"p."+pNameCol,
			"p."+pImageUrlCol,
			"p."+pGTINCol,
			"p."+pCreatedAtCol,
			"p."+pUpdatedAtCol,
			"p.category_id",
			"p.subcategory_id",
			"c.name AS category_name",
			"sc.name AS subcategory_name",
			"ps."+psPriceCol+" AS ps_price",
			"ps."+psSellAmountCol+" AS ps_sell_amount",
			"s."+sNameCol+" AS supplier_name",
			"s."+sOrderAmountCol+" AS supplier_order_amount",
			"dc."+dcFreeDeliveryAmountCol+" AS dc_min_free_delivery_amount",
			"dc."+dcDeliveryFeeCol+" AS dc_delivery_fee",
		).
		From(productsTbl + " AS p").
		LeftJoin("categories AS c ON p.category_id = c.id").
		LeftJoin("subcategories AS sc ON p.subcategory_id = sc.id").
		LeftJoin(productsSupplierTbl + " AS ps ON ps." + psProductIDCol + " = p." + pIdCol +
			" AND ps." + psSupplierIDCol + " = p." + pLowestSupplierIDCol).
		LeftJoin(supplierTbl + " AS s ON s." + sIDCol + " = ps." + psSupplierIDCol).
		LeftJoin(deliveryConditionTbl + " AS dc ON dc." + dcIDCol + " = s." + sDeliveryConditionIDCol)

	if queryParam.CategoryID != nil {
		builder = builder.Where(sq.Eq{"p.category_id": *queryParam.CategoryID})
	}

	if queryParam.SubcategoryID != nil {
		builder = builder.Where(sq.Eq{"p.subcategory_id": *queryParam.SubcategoryID})
	}

	builder = builder.
		OrderBy("p." + pIdCol).
		Limit(uint64(queryParam.Limit)).
		Offset(uint64(queryParam.Offset)).
		PlaceholderFormat(sq.Dollar)

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

	var products []repoModel.Product
	for rows.Next() {
		var product repoModel.Product
		var ps repoModel.ProductSupplier
		var s repoModel.Supplier

		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.ImageUrl,
			&product.GTIN,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.CategoryID,
			&product.SubcategoryID,
			&product.CategoryName,
			&product.SubcategoryName,
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
		product.LowestSupplier = ps
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating products: %w", err)
	}

	return converter.ToProductListFromRepo(products), nil
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

// GetProductByGTIN retrieves a product by its GTIN
func (r *repo) GetProductByGTIN(ctx context.Context, gtin int64) (*model.Product, error) {
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
		Where(sq.Eq{pGTINCol: gtin}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get product by GTIN query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetProductByGTIN",
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
		return nil, fmt.Errorf("failed to get product by GTIN: %w", err)
	}

	return converter.ToProductFromRepo(product), nil
}

// CreateProduct creates a new product in the database
func (r *repo) CreateProduct(ctx context.Context, product *model.Product) (int64, error) {
	builder := sq.
		Insert(productsTbl).
		Columns(
			pNameCol,
			pImageUrlCol,
			pGTINCol,
			"category_id",
			"subcategory_id",
			pCreatedAtCol,
			pUpdatedAtCol,
		).
		Values(
			product.Name,
			product.ImageUrl,
			product.GTIN,
			product.CategoryID,
			product.SubcategoryID,
			product.CreatedAt,
			product.UpdatedAt,
		).
		Suffix(fmt.Sprintf("RETURNING %s", pIdCol)).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build create product query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.CreateProduct",
		QueryRaw: query,
	}

	var id int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to create product: %w", err)
	}

	return id, nil
}

// CreateProductSupplier creates a new product-supplier record in the database.
func (r *repo) CreateProductSupplier(ctx context.Context, supplierID, productID int64, price int) error {
	builder := sq.
		Insert(productsSupplierTbl).
		Columns(
			psProductIDCol,
			psSupplierIDCol,
			psPriceCol,
			psSellAmountCol,
		).
		Values(
			productID,
			supplierID,
			price,
			0, // default sell_amount value
		).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build create product supplier query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.CreateProductSupplier",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to create product supplier: %w", err)
	}
	return nil
}

// GetProductListBySupplier retrieves a list of products for a specific supplier
func (r *repo) GetProductListBySupplier(ctx context.Context, supplierID int64, limit, offset int) ([]model.Product, error) {
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
		Join(productsSupplierTbl + " AS ps ON ps." + psProductIDCol + " = p." + pIdCol).
		Join(supplierTbl + " AS s ON s." + sIDCol + " = ps." + psSupplierIDCol).
		LeftJoin(deliveryConditionTbl + " AS dc ON dc." + dcIDCol + " = s." + sDeliveryConditionIDCol).
		Where(sq.Eq{"ps." + psSupplierIDCol: supplierID}).
		PlaceholderFormat(sq.Dollar)

	if limit > 0 {
		builder = builder.Limit(uint64(limit))
	} else {
		builder = builder.Limit(30)
	}
	if offset > 0 {
		builder = builder.Offset(uint64(offset))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get product list by supplier query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetProductListBySupplier",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get product list by supplier: %w", err)
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

// GetTotalProductsBySupplier returns the total number of products for a specific supplier
func (r *repo) GetTotalProductsBySupplier(ctx context.Context, supplierID int64) (int, error) {
	builder := sq.
		Select("COUNT(DISTINCT p." + pIdCol + ")").
		From(productsTbl + " AS p").
		Join(productsSupplierTbl + " AS ps ON ps." + psProductIDCol + " = p." + pIdCol).
		Where(sq.Eq{"ps." + psSupplierIDCol: supplierID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build get total products by supplier query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetTotalProductsBySupplier",
		QueryRaw: query,
	}

	var total int
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total products by supplier: %w", err)
	}

	return total, nil
}

// FindOrCreateCategory finds a category by name or creates it if it doesn't exist
func (r *repo) FindOrCreateCategory(ctx context.Context, name string) (int, error) {
	// First try to find existing category
	findBuilder := sq.
		Select("id").
		From("categories").
		Where(sq.Eq{"name": name}).
		PlaceholderFormat(sq.Dollar)

	findQuery, findArgs, err := findBuilder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build find category query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.FindCategory",
		QueryRaw: findQuery,
	}

	var categoryID int
	err = r.db.DB().QueryRowContext(ctx, q, findArgs...).Scan(&categoryID)
	if err == nil {
		return categoryID, nil // Found existing category
	}

	// Check for "no rows" errors more broadly - could be sql.ErrNoRows or wrapped/custom error
	if err == sql.ErrNoRows || strings.Contains(err.Error(), "no rows") {
		// Category doesn't exist, create it
		createBuilder := sq.
			Insert("categories").
			Columns("name").
			Values(name).
			Suffix("RETURNING id").
			PlaceholderFormat(sq.Dollar)

		createQuery, createArgs, err := createBuilder.ToSql()
		if err != nil {
			return 0, fmt.Errorf("failed to build create category query: %w", err)
		}

		createQ := db.Query{
			Name:     "product_repository.CreateCategory",
			QueryRaw: createQuery,
		}

		err = r.db.DB().QueryRowContext(ctx, createQ, createArgs...).Scan(&categoryID)
		if err != nil {
			return 0, fmt.Errorf("failed to create category: %w", err)
		}

		return categoryID, nil
	}

	// Some other error occurred during find
	return 0, fmt.Errorf("failed to find category: %w", err)
}

// FindOrCreateSubcategory finds a subcategory by name and category_id or creates it if it doesn't exist
func (r *repo) FindOrCreateSubcategory(ctx context.Context, name string, categoryID int) (int, error) {
	// First try to find existing subcategory
	findBuilder := sq.
		Select("id").
		From("subcategories").
		Where(sq.And{
			sq.Eq{"name": name},
			sq.Eq{"category_id": categoryID},
		}).
		PlaceholderFormat(sq.Dollar)

	findQuery, findArgs, err := findBuilder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build find subcategory query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.FindSubcategory",
		QueryRaw: findQuery,
	}

	var subcategoryID int
	err = r.db.DB().QueryRowContext(ctx, q, findArgs...).Scan(&subcategoryID)
	if err == nil {
		return subcategoryID, nil // Found existing subcategory
	}

	// Check for "no rows" errors more broadly - could be sql.ErrNoRows or wrapped/custom error
	if err == sql.ErrNoRows || strings.Contains(err.Error(), "no rows") {
		// Subcategory doesn't exist, create it
		createBuilder := sq.
			Insert("subcategories").
			Columns("name", "category_id").
			Values(name, categoryID).
			Suffix("RETURNING id").
			PlaceholderFormat(sq.Dollar)

		createQuery, createArgs, err := createBuilder.ToSql()
		if err != nil {
			return 0, fmt.Errorf("failed to build create subcategory query: %w", err)
		}

		createQ := db.Query{
			Name:     "product_repository.CreateSubcategory",
			QueryRaw: createQuery,
		}

		err = r.db.DB().QueryRowContext(ctx, createQ, createArgs...).Scan(&subcategoryID)
		if err != nil {
			return 0, fmt.Errorf("failed to create subcategory: %w", err)
		}

		return subcategoryID, nil
	}

	// Some other error occurred during find
	return 0, fmt.Errorf("failed to find subcategory: %w", err)
}

// GetCategories retrieves all categories
func (r *repo) GetCategories(ctx context.Context) ([]model.Category, error) {
	builder := sq.
		Select("id", "name", "description", "created_at", "updated_at").
		From("categories").
		OrderBy("name").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get categories query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetCategories",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var category model.Category
		err := rows.Scan(
			&category.ID,
			&category.Name,
			&category.Description,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categories, nil
}

// GetCategory retrieves a single category by ID
func (r *repo) GetCategory(ctx context.Context, id int) (*model.Category, error) {
	builder := sq.
		Select("id", "name", "description", "created_at", "updated_at").
		From("categories").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get category query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetCategory",
		QueryRaw: query,
	}

	var category model.Category
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return &category, nil
}

// GetSubcategories retrieves all subcategories for a specific category
func (r *repo) GetSubcategories(ctx context.Context, categoryID int) ([]model.Subcategory, error) {
	builder := sq.
		Select("id", "category_id", "name", "description", "created_at", "updated_at").
		From("subcategories").
		Where(sq.Eq{"category_id": categoryID}).
		OrderBy("name").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get subcategories query: %w", err)
	}

	q := db.Query{
		Name:     "product_repository.GetSubcategories",
		QueryRaw: query,
	}

	rows, err := r.db.DB().QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get subcategories: %w", err)
	}
	defer rows.Close()

	var subcategories []model.Subcategory
	for rows.Next() {
		var subcategory model.Subcategory
		err := rows.Scan(
			&subcategory.ID,
			&subcategory.CategoryID,
			&subcategory.Name,
			&subcategory.Description,
			&subcategory.CreatedAt,
			&subcategory.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subcategory: %w", err)
		}
		subcategories = append(subcategories, subcategory)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating subcategories: %w", err)
	}

	return subcategories, nil
}
