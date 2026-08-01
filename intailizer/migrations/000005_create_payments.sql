-- +goose Up
CREATE TABLE IF NOT EXISTS paymenttable (
    id SERIAL PRIMARY KEY,
    auth_id INT NOT NULL REFERENCES cometoseeauth(auth_id) ON DELETE CASCADE,
    transaction_uuid TEXT NOT NULL UNIQUE,
    amount DECIMAL(10,2) NOT NULL,
    plan TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending',
    payment_date TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS paymenttable;

