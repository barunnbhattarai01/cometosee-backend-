-- +goose Up
CREATE TABLE IF NOT EXISTS connectionstable (
    id BIGSERIAL PRIMARY KEY,
    user_id_1 TEXT NOT NULL,
    user_id_2 TEXT NOT NULL,
    status TEXT NOT NULL,
    requested_by TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CHECK (user_id_1 < user_id_2)
);
CREATE UNIQUE INDEX IF NOT EXISTS unique_pair ON connectionstable(user_id_1, user_id_2);

-- +goose Down
DROP TABLE IF EXISTS connectionstable;

