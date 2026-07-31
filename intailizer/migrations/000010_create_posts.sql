-- +goose Up
CREATE TABLE IF NOT EXISTS post (
    post_id SERIAL PRIMARY KEY,
    auth_id INT NOT NULL REFERENCES cometoseeauth(auth_id) ON DELETE CASCADE,
    caption TEXT NOT NULL,
    images_url TEXT,
    venue TEXT NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    sport TEXT NOT NULL,
    room_id TEXT UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS post;
