-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
    ALTER TABLE suppliers
    ADD COLUMN  min_order_amount INTEGER NOT NULL;
    
    ALTER TABLE products_supplier
    RENAME COLUMN min_sell_amount TO sell_amount;

    
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
