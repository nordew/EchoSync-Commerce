-- +goose Up

-- +goose StatementBegin
CREATE TABLE stores (
                        store_id UUID PRIMARY KEY,
                        store_name VARCHAR(100) NOT NULL,
                        owner_user_id UUID REFERENCES users(user_id),
                        products_count INT DEFAULT 0,
                        is_active BOOLEAN DEFAULT false,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stores;
-- +goose StatementEnd
