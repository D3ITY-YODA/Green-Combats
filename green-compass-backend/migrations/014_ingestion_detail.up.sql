-- Datasets: individual data products from a source
CREATE TABLE datasets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_id       UUID NOT NULL REFERENCES data_sources(id) ON DELETE RESTRICT,
    external_id     TEXT NOT NULL,
    key             TEXT NOT NULL,
    name            TEXT NOT NULL,
    description     TEXT,
    format          TEXT NOT NULL,
    access_method   TEXT NOT NULL,
    license         TEXT,
    refresh_schedule TEXT,
    active          BOOLEAN NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (source_id, external_id)
);

CREATE TRIGGER datasets_set_updated_at
    BEFORE UPDATE ON datasets
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Source runs: track each ingestion execution
CREATE TABLE source_runs (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    dataset_id          UUID NOT NULL REFERENCES datasets(id) ON DELETE RESTRICT,
    place_id            UUID NOT NULL REFERENCES places(id) ON DELETE RESTRICT,
    status              TEXT NOT NULL DEFAULT 'created' CHECK (status IN ('created', 'queued', 'running', 'fetched', 'validated', 'processed', 'completed', 'failed', 'quarantined')),
    started_at          TIMESTAMPTZ,
    finished_at         TIMESTAMPTZ,
    records_fetched     INTEGER NOT NULL DEFAULT 0 CHECK (records_fetched >= 0),
    records_accepted    INTEGER NOT NULL DEFAULT 0 CHECK (records_accepted >= 0),
    records_rejected    INTEGER NOT NULL DEFAULT 0 CHECK (records_rejected >= 0),
    error_code          TEXT,
    error_message       TEXT,
    idempotency_key     TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT source_runs_finished_after_started CHECK (finished_at IS NULL OR finished_at >= started_at),
    CONSTRAINT source_runs_terminal_status CHECK (
        (status IN ('created', 'queued', 'running', 'fetched', 'validated', 'processed') AND finished_at IS NULL) OR
        (status IN ('completed', 'failed', 'quarantined') AND finished_at IS NOT NULL)
    )
);

CREATE INDEX idx_source_runs_dataset_place_started
    ON source_runs (dataset_id, place_id, started_at DESC);

CREATE UNIQUE INDEX idx_source_runs_idempotency_key
    ON source_runs (idempotency_key) WHERE idempotency_key IS NOT NULL;

-- Raw assets: archived raw responses from sources
CREATE TABLE raw_assets (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_run_id   UUID NOT NULL REFERENCES source_runs(id) ON DELETE RESTRICT,
    object_key      TEXT NOT NULL,
    content_type    TEXT NOT NULL,
    size_bytes      BIGINT NOT NULL CHECK (size_bytes >= 0),
    checksum        TEXT NOT NULL CHECK (checksum ~ '^[a-f0-9]{64}$'),
    retrieved_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    license         TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_raw_assets_source_run
    ON raw_assets (source_run_id);

-- Update ingestion_runs to reference datasets instead of data_sources directly
-- (keep existing ingestion_runs for backward compatibility during transition)
-- New ingestion will use source_runs table above