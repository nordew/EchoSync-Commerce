-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
                       user_id UUID PRIMARY KEY,
                       username VARCHAR(50) NOT NULL,
                       email VARCHAR(100) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       refresh_token VARCHAR(255),
                       stores_active INT DEFAULT 0,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE stores (
                        store_id UUID DEFAULT PRIMARY KEY,
                        store_name VARCHAR(100) NOT NULL,
                        owner_id UUID REFERENCES users(user_id) ON DELETE CASCADE,
                        description TEXT,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS stores;
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
