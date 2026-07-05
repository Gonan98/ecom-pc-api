-- +goose Up
CREATE TABLE order_details (
    order_id INTEGER NOT NULL,
    product_id INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    unit_price NUMERIC(10,2) NOT NULL,
    discount FLOAT NOT NULL DEFAULT 0.0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (order_id, product_id),
    FOREIGN KEY (order_id) REFERENCES orders(id),
    FOREIGN KEY (product_id) REFERENCES products(id),

    CONSTRAINT chk_details_quantity CHECK (quantity > 0),
    CONSTRAINT chk_details_unit_price CHECK (unit_price > 0),
    CONSTRAINT chk_details_discount CHECK (discount >= 0 AND discount <= 1)
);

-- +goose Down
DROP TABLE order_details;
