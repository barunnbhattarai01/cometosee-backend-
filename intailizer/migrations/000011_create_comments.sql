-- +goose Up
CREATE TABLE IF NOT EXISTS comments (
    comment_id SERIAL PRIMARY KEY,
    auth_id INT NOT NULL REFERENCES cometoseeauth(auth_id) ON DELETE CASCADE,
    post_id INT NOT NULL REFERENCES post(post_id) ON DELETE CASCADE,
    comment TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS comments;

