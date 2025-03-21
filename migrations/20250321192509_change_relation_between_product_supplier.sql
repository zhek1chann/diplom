-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

ALTER TABLE products
ADD CONSTRAINT lowest_supplier_id
FOREIGN KEY (lowest_supplier_id)
REFERENCES suppliers (user_id);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
