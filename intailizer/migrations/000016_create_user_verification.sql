-- +goose Up
CREATE TABLE IF NOT EXISTS user_verification (
    verification_id SERIAL PRIMARY KEY,
    auth_id INT NOT NULL UNIQUE REFERENCES cometoseeauth(auth_id) ON DELETE CASCADE,
    citizenship_front TEXT NOT NULL,
    citizenship_back TEXT NOT NULL,
    verification_status TEXT NOT NULL DEFAULT 'pending',
    rejection_reason TEXT,
    verified_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS user_verification;

