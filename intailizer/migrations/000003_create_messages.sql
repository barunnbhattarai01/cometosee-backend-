-- +goose Up
CREATE TABLE IF NOT EXISTS messagetable (
    id SERIAL PRIMARY KEY,
    sender TEXT NOT NULL,
    room TEXT NOT NULL,
    message TEXT NOT NULL,
    receiver TEXT NOT NULL,
    sent_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS messagetable;

