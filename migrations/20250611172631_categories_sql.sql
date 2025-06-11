-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

-- Create categories table
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create subcategories table
CREATE TABLE subcategories (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

-- Add foreign key constraints to existing products table
ALTER TABLE products
    ADD CONSTRAINT fk_product_category    FOREIGN KEY (category_id)    REFERENCES categories(id),
    ADD CONSTRAINT fk_product_subcategory FOREIGN KEY (subcategory_id) REFERENCES subcategories(id);


CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_subcategory ON products(subcategory_id);

-- +goose Down
-- SQL in section 'Down' is executed when this migration is rolled back
DROP INDEX IF EXISTS idx_subcategories_category;
DROP INDEX IF EXISTS idx_products_subcategory;
DROP INDEX IF EXISTS idx_products_category;

ALTER TABLE products
DROP CONSTRAINT IF EXISTS fk_product_subcategory,
DROP CONSTRAINT IF EXISTS fk_product_category;


DROP TABLE IF EXISTS subcategories;
DROP TABLE IF EXISTS categories; 