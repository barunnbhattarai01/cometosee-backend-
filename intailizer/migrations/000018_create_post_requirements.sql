-- +goose Up
CREATE TABLE IF NOT EXISTS post_requirements (
    requirement_id SERIAL PRIMARY KEY,
    post_id INT NOT NULL REFERENCES post(post_id) ON DELETE CASCADE,
    min_age INT,
    max_age INT,
    gender TEXT,
    skill_level TEXT,
    verification_required BOOLEAN DEFAULT FALSE,
    player_document_required BOOLEAN DEFAULT FALSE,
    description TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    CHECK (min_age IS NULL OR max_age IS NULL OR min_age <= max_age),
    UNIQUE (post_id)
);

-- +goose Down
DROP TABLE IF EXISTS post_requirements;

