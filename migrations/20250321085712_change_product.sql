-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

ALTER TABLE products
    RENAME COLUMN supplier_id TO lowest_supplier_id;

ALTER table products
    Rename COLUMN min_price to lowest_price;

ALTER TABLE suppliers
    Add COLUMN name VARCHAR(40);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
