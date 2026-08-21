-- Normalized observations: canonical internal format after normalization
CREATE TABLE normalized_observations (
    id                  UUID NOT NULL DEFAULT gen_random_uuid(),
    source_key          TEXT NOT NULL,
    dataset_key         TEXT NOT NULL,
    topic_key           TEXT NOT NULL,
    variable            TEXT NOT NULL,
    geometry            GEOGRAPHY(GEOMETRY, 4326),
    place_id            UUID REFERENCES places(id),
    value               DOUBLE PRECISION,
    text_value          TEXT,
    unit                TEXT,
    observed_at         TIMESTAMPTZ NOT NULL,
    valid_from          TIMESTAMPTZ,
    valid_until         TIMESTAMPTZ,
    retrieved_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    is_forecast         BOOLEAN NOT NULL DEFAULT FALSE,
    quality_status      TEXT NOT NULL DEFAULT 'accepted' CHECK (quality_status IN ('accepted', 'rejected', 'quarantined', 'unverified', 'stale', 'superseded')),
    source_version      TEXT,
    raw_asset_ref       TEXT NOT NULL,
    license             TEXT,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (id, observed_at)
) PARTITION BY RANGE (observed_at);

CREATE INDEX idx_normalized_observations_geometry
    ON normalized_observations USING GIST (geometry);

CREATE INDEX idx_normalized_observations_topic_time
    ON normalized_observations (topic_key, observed_at DESC);

CREATE INDEX idx_normalized_observations_place_time
    ON normalized_observations (place_id, observed_at DESC);

CREATE INDEX idx_normalized_observations_source_dataset
    ON normalized_observations (source_key, dataset_key, observed_at DESC);

-- Monthly partitions for the next 2 years
DO $$
DECLARE
    start_date DATE := DATE_TRUNC('month', NOW()) - INTERVAL '1 month';
    end_date DATE := start_date + INTERVAL '24 months';
    partition_start DATE;
    partition_end DATE;
    partition_name TEXT;
BEGIN
    WHILE start_date < end_date LOOP
        partition_start := start_date;
        partition_end := start_date + INTERVAL '1 month';
        partition_name := 'normalized_observations_' || TO_CHAR(partition_start, 'YYYY_MM');
        
        EXECUTE FORMAT('
            CREATE TABLE %I PARTITION OF normalized_observations
            FOR VALUES FROM (%L) TO (%L)',
            partition_name, partition_start, partition_end);
        
        start_date := partition_end;
    END LOOP;
END $$;

-- Function to create future partitions automatically
CREATE OR REPLACE FUNCTION create_monthly_partition(table_name TEXT, partition_date DATE)
RETURNS VOID AS $$
DECLARE
    partition_start DATE := DATE_TRUNC('month', partition_date);
    partition_end DATE := partition_start + INTERVAL '1 month';
    partition_name TEXT := table_name || '_' || TO_CHAR(partition_start, 'YYYY_MM');
BEGIN
    EXECUTE FORMAT('
        CREATE TABLE IF NOT EXISTS %I PARTITION OF %I
        FOR VALUES FROM (%L) TO (%L)',
        partition_name, table_name, partition_start, partition_end);
END $$ LANGUAGE plpgsql;
