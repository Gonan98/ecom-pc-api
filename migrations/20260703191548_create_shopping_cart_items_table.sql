-- +goose Up
CREATE TABLE shopping_cart_items(
    cart_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (cart_id, product_id),
    FOREIGN KEY (cart_id) REFERENCES shopping_carts(id),
    FOREIGN KEY (product_id) REFERENCES products(id),

    CONSTRAINT chk_items_quantity CHECK (quantity > 0)
);

-- +goose Down
DROP TABLE shopping_cart_items;
