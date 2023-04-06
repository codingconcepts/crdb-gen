CREATE DATABASE curious_cupcake;
USE curious_cupcake;

CREATE TABLE products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name STRING NOT NULL,
    price DECIMAL NOT NULL,
    available_from TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE members (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email STRING NOT NULL
);

CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    member_id UUID NOT NULL REFERENCES members(id),
    product_ids UUID[],
    total DECIMAL NOT NULL DEFAULT 0,
    checkout_at TIMESTAMPTZ
);