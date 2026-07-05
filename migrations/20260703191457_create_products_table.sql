-- +goose Up
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    category_id INTEGER NOT NULL,
    brand_id INTEGER NOT NULL,
    name VARCHAR(50) NOT NULL UNIQUE,
    description VARCHAR(100),
    image_url VARCHAR(200),
    price NUMERIC(10,2) NOT NULL,
    stock INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_category FOREIGN KEY(category_id) REFERENCES categories(id),
    CONSTRAINT fk_brand FOREIGN KEY(brand_id) REFERENCES brands(id),

    CONSTRAINT chk_products_price CHECK (price >= 0.0),
    CONSTRAINT chk_products_stock CHECK (stock >= 0)
);

CREATE INDEX idx_products_name ON products(name);
CREATE INDEX idx_products_category ON products(category_id);
CREATE INDEX idx_products_brand ON products(brand_id);
CREATE INDEX idx_products_active ON products(is_active);

-- +goose Down
DROP TABLE products;
