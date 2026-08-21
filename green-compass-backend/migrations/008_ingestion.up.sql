CREATE TABLE data_sources (
    id                    UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code                  TEXT NOT NULL UNIQUE,
    display_name          TEXT NOT NULL,
    poll_interval_seconds INTEGER NOT NULL CHECK (poll_interval_seconds > 0),
    enabled               BOOLEAN NOT NULL DEFAULT FALSE,
    config                JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at            TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at            TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER data_sources_set_updated_at
    BEFORE UPDATE ON data_sources
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

INSERT INTO data_sources (code, display_name, poll_interval_seconds, enabled) VALUES
    ('open_meteo', 'Open-Meteo', 21600, TRUE),
    ('nasa_power', 'NASA POWER', 86400, TRUE),
    ('freshwater', 'Freshwater provider', 21600, FALSE),
    ('satellite', 'Satellite provider', 86400, FALSE),
    ('community', 'Community reports', 3600, FALSE);

CREATE TABLE ingestion_runs (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id           UUID NOT NULL REFERENCES data_sources(id) ON DELETE RESTRICT,
    place_id            UUID NOT NULL REFERENCES places(id) ON DELETE RESTRICT,
    status              TEXT NOT NULL CHECK (status IN ('running', 'succeeded', 'failed')),
    requested_from      TIMESTAMPTZ,
    requested_to        TIMESTAMPTZ,
    started_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    finished_at         TIMESTAMPTZ,
    landed_record_count INTEGER NOT NULL DEFAULT 0 CHECK (landed_record_count >= 0),
    error_message       TEXT,
    CONSTRAINT ingestion_runs_finished_after_started CHECK (finished_at IS NULL OR finished_at >= started_at),
    CONSTRAINT ingestion_runs_finished_status CHECK (
        (status = 'running' AND finished_at IS NULL) OR
        (status IN ('succeeded', 'failed') AND finished_at IS NOT NULL)
    )
);

CREATE INDEX ingestion_runs_source_place_started_idx
    ON ingestion_runs (source_id, place_id, started_at DESC);

CREATE TABLE raw_records (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    ingestion_run_id    UUID NOT NULL REFERENCES ingestion_runs(id) ON DELETE RESTRICT,
    source_id           UUID NOT NULL REFERENCES data_sources(id) ON DELETE RESTRICT,
    place_id            UUID NOT NULL REFERENCES places(id) ON DELETE RESTRICT,
    source_observed_at  TIMESTAMPTZ NOT NULL,
    fetched_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    source_url          TEXT NOT NULL,
    content_type        TEXT NOT NULL,
    payload             JSONB NOT NULL,
    payload_checksum    TEXT NOT NULL CHECK (payload_checksum ~ '^[a-f0-9]{64}$'),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (source_id, place_id, source_observed_at, payload_checksum)
);

CREATE INDEX raw_records_source_place_observed_idx
    ON raw_records (source_id, place_id, source_observed_at DESC);
