-- +goose Up
CREATE TABLE IF NOT EXISTS post_slots (
    slot_id SERIAL PRIMARY KEY,
    post_id INT NOT NULL REFERENCES post(post_id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    max_participants INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT NOW(),
    CHECK (end_time > start_time)
);

-- +goose Down
DROP TABLE IF EXISTS post_slots;

