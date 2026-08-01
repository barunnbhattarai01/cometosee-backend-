-- +goose Up
CREATE TABLE IF NOT EXISTS slot_participants (
    id SERIAL PRIMARY KEY,
    slot_id INT NOT NULL REFERENCES post_slots(slot_id) ON DELETE CASCADE,
    auth_id INT NOT NULL REFERENCES cometoseeauth(auth_id) ON DELETE CASCADE,
    qr_token TEXT UNIQUE,
    qr_expires_at TIMESTAMPTZ,
    checked_in BOOLEAN DEFAULT FALSE,
    checked_in_at TIMESTAMPTZ,
    checked_in_by INT REFERENCES cometoseeauth(auth_id),
    joined_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (slot_id, auth_id)
);

-- +goose Down
DROP TABLE IF EXISTS slot_participants;

