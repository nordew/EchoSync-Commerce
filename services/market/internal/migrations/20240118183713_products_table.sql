-- +goose Up
-- +goose StatementBegin
CREATE TABLE products (
                          product_id UUID PRIMARY KEY,
                          store_id UUID REFERENCES stores(store_id),
                          product_name VARCHAR(100) NOT NULL,
                          price DECIMAL(10, 2) NOT NULL,
                          quantity INT DEFAULT 0,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE products IF EXISTS;
-- +goose StatementEnd
