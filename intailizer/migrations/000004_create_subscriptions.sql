-- +goose Up
CREATE TABLE IF NOT EXISTS subscriptiontable (
    id SERIAL PRIMARY KEY,
    auth_id INT NOT NULL REFERENCES cometoseeauth(auth_id) ON DELETE CASCADE,
    plan TEXT NOT NULL DEFAULT 'standard',
    status TEXT NOT NULL DEFAULT 'active',
    start_date TIMESTAMP NOT NULL DEFAULT NOW(),
    end_date TIMESTAMP NOT NULL,
    UNIQUE (auth_id)
);

-- +goose Down
DROP TABLE IF EXISTS subscriptiontable;

