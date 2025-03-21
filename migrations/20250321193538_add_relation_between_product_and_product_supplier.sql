-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose 

ALTER TABLE products
ADD CONSTRAINT lowest_supplier_id
FOREIGN KEY (lowest_supplier_id)
REFERENCES products_supplier (supplier_id);


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
