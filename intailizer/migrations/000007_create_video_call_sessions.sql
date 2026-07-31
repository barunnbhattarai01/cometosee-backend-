-- +goose Up
CREATE TABLE IF NOT EXISTS video_call_sessions (
    id BIGSERIAL PRIMARY KEY,
    connection_id BIGINT NOT NULL,
    initiated_by_user_id BIGINT NOT NULL,
    agora_channel_name TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'initiated',
    started_at TIMESTAMP NULL,
    ended_at TIMESTAMP NULL,
    duration_seconds INT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS video_call_sessions;

