-- +goose Up
CREATE TABLE IF NOT EXISTS post_shares (
    share_id SERIAL PRIMARY KEY,
    post_id INT NOT NULL REFERENCES post(post_id) ON DELETE CASCADE,
    auth_id INT NOT NULL REFERENCES cometoseeauth(auth_id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (post_id, auth_id)
);

-- +goose Down
DROP TABLE IF EXISTS post_shares;

