-- +goose Up
CREATE TABLE IF NOT EXISTS cometoseeauth (
    auth_id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    is_verified BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS cometoseeauth;

