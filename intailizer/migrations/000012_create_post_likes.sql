-- +goose Up
CREATE TABLE IF NOT EXISTS post_likes (
    like_id SERIAL PRIMARY KEY,
    created_at TIMESTAMP DEFAULT NOW(),
    auth_id INT NOT NULL REFERENCES cometoseeauth(auth_id) ON DELETE CASCADE,
    post_id INT NOT NULL REFERENCES post(post_id) ON DELETE CASCADE,
    UNIQUE (post_id, auth_id)
);

-- +goose Down
DROP TABLE IF EXISTS post_likes;

