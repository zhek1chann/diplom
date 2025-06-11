-- +goose Up
-- SQL в секции Up выполняется при применении миграции
------------------------------------------------------------------------

-- Create price_history table for tracking product price changes
CREATE TABLE price_history (
    id SERIAL PRIMARY KEY,
    product_id INTEGER NOT NULL,
    supplier_id INTEGER NOT NULL,
    price INTEGER NOT NULL,
    date TIMESTAMP NOT NULL,
    change_reason VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE,
    FOREIGN KEY (supplier_id) REFERENCES suppliers(user_id) ON DELETE CASCADE
);

-- Create indexes for better performance
CREATE INDEX idx_price_history_product_id ON price_history(product_id);
CREATE INDEX idx_price_history_supplier_id ON price_history(supplier_id);
CREATE INDEX idx_price_history_date ON price_history(date);
CREATE INDEX idx_price_history_product_supplier ON price_history(product_id, supplier_id);

------------------------------------------------------------------------
-- +goose Down
-- SQL в секции Down выполняется при откате миграции
------------------------------------------------------------------------

-- Drop indexes first
DROP INDEX IF EXISTS idx_price_history_product_supplier;
DROP INDEX IF EXISTS idx_price_history_date;
DROP INDEX IF EXISTS idx_price_history_supplier_id;
DROP INDEX IF EXISTS idx_price_history_product_id;

-- Drop table
DROP TABLE IF EXISTS price_history; 