-- +goose Up
ALTER TABLE shopping_cart_items
ADD discount NUMERIC(10,2) NOT NULL DEFAULT 0.0;

-- +goose Down
ALTER TABLE shopping_cart_items
DROP COLUMN discount;
