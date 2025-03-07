-- +goose Up
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    supplier_id INTEGER,
    min_price INTEGER,
    image_url TEXT,
    name TEXT NOT NULL,
    kzt_price INTEGER,
    kzt_min_sell_amount INTEGER,
    kztin BIGINT NOT NULL,
    gtin BIGINT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE products_supplier (
    product_id INTEGER NOT NULL,
    supplier_id INTEGER NOT NULL,
    price INTEGER NOT NULL,
    min_sell_amount INTEGER NOT NULL,
    PRIMARY KEY (product_id, supplier_id)
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    phone_number TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    role INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);



-- +goose Down
DROP TABLE products_supplier;
DROP TABLE products;
DROP TABLE users;
