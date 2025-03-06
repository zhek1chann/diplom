-- +goose Up
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    kztin BIGINT NOT NULL,
    gtin BIGINT, 
    image_url TEXT,
    min_price int,
    supplier_id int,
);

CREATE TABLE products_supplier(
    product_id int
    supplier_id int 
    price int
    min_sell_amount int
);

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    phone_number TEXT NOT NULL,
    hashed_password TEXT NOT NULL,
    role INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- +goose Down
drop table products;
drop table users;