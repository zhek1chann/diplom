-- +goose Up
CREATE TABLE delivery_conditions (
    condition_id SERIAL PRIMARY KEY,
    minimum_free_delivery_amount NUMERIC,
    delivery_fee NUMERIC
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    phone_number VARCHAR(12) NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    role INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE suppliers (
    user_id INTEGER NOT NULL,
    
    condition_id INT,
    PRIMARY KEY (user_id),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (condition_id) REFERENCES delivery_conditions (condition_id)
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    image_url TEXT,
    gtin BIGINT,
    min_price INTEGER,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    supplier_id INTEGER,
    FOREIGN KEY (supplier_id) REFERENCES users (id)
);

CREATE TABLE products_supplier (
    product_id INTEGER NOT NULL,
    supplier_id INTEGER NOT NULL,
    price INTEGER NOT NULL,
    min_sell_amount INTEGER NOT NULL,
    PRIMARY KEY (product_id, supplier_id),
    FOREIGN KEY (product_id) REFERENCES products (id),
    FOREIGN KEY (supplier_id) REFERENCES users (id)
);

-- +goose Down
DROP TABLE products_supplier;
DROP TABLE products;
DROP TABLE suppliers;
DROP TABLE delivery_conditions;
DROP TABLE users;
