CREATE TABLE places (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name          TEXT NOT NULL,
    place_type    TEXT NOT NULL CHECK (place_type IN ('community', 'ward', 'district', 'custom')),
    location      GEOGRAPHY(POINT, 4326) NOT NULL,
    external_code TEXT,
    created_by    UUID REFERENCES users(id),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX places_location_gix ON places USING GIST (location);
CREATE INDEX places_name_idx ON places (name);

CREATE TRIGGER places_set_updated_at
    BEFORE UPDATE ON places
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

CREATE TABLE user_saved_places (
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    place_id   UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    label      TEXT,
    is_primary BOOLEAN NOT NULL DEFAULT FALSE,
    saved_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, place_id)
);

CREATE UNIQUE INDEX user_saved_places_single_primary ON user_saved_places (user_id) WHERE is_primary;
