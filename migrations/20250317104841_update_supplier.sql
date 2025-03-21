-- +goose Up

-- Add the 'order_amount' column to 'suppliers' table
ALTER TABLE suppliers
    ADD COLUMN order_amount INTEGER NOT NULL;

-- Rename 'min_sell_amount' to 'sell_amount' in 'products_supplier' table
ALTER TABLE products_supplier
    RENAME COLUMN min_sell_amount TO sell_amount;

-- Create 'carts' table
CREATE TABLE carts (
    id SERIAL PRIMARY KEY,
    customer_id INTEGER,
    total INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    FOREIGN KEY (customer_id) REFERENCES users(id) -- Removed the trailing comma here
);

-- Create 'cart_items' table
CREATE TABLE cart_items (
    id SERIAL PRIMARY KEY,
    cart_id INTEGER REFERENCES carts(id) ON DELETE CASCADE,
    product_id INTEGER,
    supplier_id INTEGER, 
    quantity INTEGER NOT NULL,
    price INTEGER NOT NULL,
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (supplier_id) REFERENCES suppliers(user_id)
);

-- Create 'orders' table
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    cart_id INTEGER REFERENCES carts(id) ON DELETE CASCADE,
    order_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(255) NOT NULL DEFAULT 'Pending',
    customer_id INTEGER,
    supplier_id INTEGER,
    FOREIGN KEY (customer_id) REFERENCES users(id),
    FOREIGN KEY (supplier_id) REFERENCES suppliers(user_id)
);

-- Create index for 'orders' table on 'supplier_id'
CREATE INDEX idx_orders_supplier_id ON orders(supplier_id);

-- +goose Down

-- Revert the changes (for 'down' migration)

-- Drop 'orders' table
DROP TABLE orders;

-- Drop 'cart_items' table
DROP TABLE cart_items;

-- Drop 'carts' table
DROP TABLE carts;

-- Remove 'order_amount' column from 'suppliers' table
ALTER TABLE suppliers
    DROP COLUMN order_amount;

-- Rename 'sell_amount' back to 'min_sell_amount' in 'products_supplier' table
ALTER TABLE products_supplier
    RENAME COLUMN sell_amount TO min_sell_amount;
