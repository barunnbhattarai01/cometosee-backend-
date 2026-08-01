-- +goose Up
CREATE EXTENSION IF NOT EXISTS postgis;

-- +goose Down
-- PostGIS is shared database infrastructure and is intentionally not removed.

