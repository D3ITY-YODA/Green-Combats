CREATE TABLE observations (
    id               UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reporter_id      UUID NOT NULL REFERENCES users(id),
    place_id         UUID REFERENCES places(id),
    location         GEOGRAPHY(POINT, 4326),
    category         TEXT NOT NULL CHECK (category IN ('flood', 'drought', 'water_quality', 'crop_damage', 'air_quality', 'other')),
    description      TEXT NOT NULL,
    photo_object_key TEXT,
    status           TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'verified', 'rejected', 'flagged')),
    verified_by      UUID REFERENCES users(id),
    verified_at      TIMESTAMPTZ,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT observations_location_required CHECK (location IS NOT NULL OR place_id IS NOT NULL)
);

CREATE INDEX observations_location_gix ON observations USING GIST (location);
CREATE INDEX observations_status_created_idx ON observations (status, created_at DESC);
CREATE INDEX observations_reporter_idx ON observations (reporter_id);

CREATE TRIGGER observations_set_updated_at
    BEFORE UPDATE ON observations
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
