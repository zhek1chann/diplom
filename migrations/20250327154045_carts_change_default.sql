-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
    ALTER TABLE carts ALTER COLUMN total SET DEFAULT 0;

    ALTER TABLE cart_items DROP COLUMN id;

    -- Add composite primary key using cart_id, supplier_id, and product_id
    ALTER TABLE cart_items ADD PRIMARY KEY (cart_id, supplier_id, product_id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
