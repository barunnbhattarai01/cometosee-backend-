-- +goose Up
CREATE TABLE IF NOT EXISTS userdetailinfo (
    user_detail_id SERIAL PRIMARY KEY,
    auth_id INT NOT NULL REFERENCES cometoseeauth(auth_id) ON DELETE CASCADE,
    calling_name TEXT NOT NULL,
    sport TEXT NOT NULL,
    skill TEXT NOT NULL,
    avatar TEXT,
    bio TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS userdetailinfo;

