-- +goose Up
CREATE TABLE IF NOT EXISTS player_documents (
    document_id SERIAL PRIMARY KEY,
    auth_id INT NOT NULL REFERENCES cometoseeauth(auth_id) ON DELETE CASCADE,
    document_name TEXT NOT NULL,
    document_type TEXT NOT NULL,
    document_url TEXT NOT NULL,
    issued_by TEXT,
    issue_date DATE,
    rejection_reason TEXT,
    verification_status TEXT DEFAULT 'pending',
    verified_at TIMESTAMP DEFAULT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS player_documents;

