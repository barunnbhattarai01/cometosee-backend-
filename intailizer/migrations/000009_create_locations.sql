-- +goose Up
CREATE TABLE IF NOT EXISTS location (
    id SERIAL PRIMARY KEY,
    user_detail_id INT NOT NULL REFERENCES userdetailinfo(user_detail_id) ON DELETE CASCADE,
    country TEXT NOT NULL,
    city TEXT NOT NULL,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL
);
ALTER TABLE location ADD COLUMN IF NOT EXISTS geom geography(Point, 4326);
CREATE INDEX IF NOT EXISTS idx_location_geom ON location USING GIST (geom);

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_geom_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.geom = ST_SetSRID(ST_MakePoint(NEW.longitude, NEW.latitude), 4326)::geography;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

DROP TRIGGER IF EXISTS set_geom ON location;
CREATE TRIGGER set_geom BEFORE INSERT OR UPDATE ON location
FOR EACH ROW EXECUTE FUNCTION update_geom_column();

-- +goose Down
DROP TRIGGER IF EXISTS set_geom ON location;
DROP FUNCTION IF EXISTS update_geom_column();
DROP TABLE IF EXISTS location;
