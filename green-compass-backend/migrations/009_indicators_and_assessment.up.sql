-- Computed indicators and assessment signals for plain-language content generation

CREATE TABLE indicator_definitions (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code                 TEXT NOT NULL UNIQUE,
    display_name         TEXT NOT NULL,
    category             TEXT NOT NULL CHECK (category IN ('weather', 'water', 'agriculture', 'air_quality')),
    source_variables    TEXT[] NOT NULL,
    description          TEXT NOT NULL,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TRIGGER indicator_definitions_set_updated_at
    BEFORE UPDATE ON indicator_definitions
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Computed indicator values for a place over a time period
CREATE TABLE place_indicators (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    place_id             UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    indicator_id         UUID NOT NULL REFERENCES indicator_definitions(id) ON DELETE RESTRICT,
    computed_at          TIMESTAMPTZ NOT NULL,
    period_start         TIMESTAMPTZ NOT NULL,
    period_end           TIMESTAMPTZ NOT NULL,
    value                NUMERIC(10, 3) NOT NULL,
    unit                 TEXT NOT NULL,
    trend                TEXT CHECK (trend IN ('increasing', 'stable', 'decreasing')),
    trend_confidence     NUMERIC(3, 2) CHECK (trend_confidence IS NULL OR (trend_confidence >= 0 AND trend_confidence <= 1)),
    data_points_count    INTEGER NOT NULL CHECK (data_points_count > 0),
    metadata             JSONB,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (place_id, indicator_id, period_start, period_end)
);

CREATE INDEX place_indicators_place_computed_idx
    ON place_indicators (place_id, computed_at DESC);
CREATE INDEX place_indicators_period_idx
    ON place_indicators (place_id, period_start, period_end);

-- Applicability rules: which indicators matter for which place types
CREATE TABLE indicator_applicability_rules (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    indicator_id         UUID NOT NULL REFERENCES indicator_definitions(id) ON DELETE CASCADE,
    place_type           TEXT NOT NULL,
    applicable           BOOLEAN NOT NULL DEFAULT TRUE,
    min_threshold        NUMERIC(10, 3),
    max_threshold        NUMERIC(10, 3),
    relevance_score      NUMERIC(3, 2) NOT NULL DEFAULT 1.0 CHECK (relevance_score >= 0 AND relevance_score <= 1.0),
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (indicator_id, place_type)
);

CREATE TRIGGER indicator_applicability_rules_set_updated_at
    BEFORE UPDATE ON indicator_applicability_rules
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Assessment signals: backend-only scoring (never exposed in API)
CREATE TABLE place_assessments (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    place_id             UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    assessed_at          TIMESTAMPTZ NOT NULL,
    period_start         TIMESTAMPTZ NOT NULL,
    period_end           TIMESTAMPTZ NOT NULL,
    urgency_score        INTEGER NOT NULL CHECK (urgency_score >= 0 AND urgency_score <= 100),
    confidence_score     INTEGER NOT NULL CHECK (confidence_score >= 0 AND confidence_score <= 100),
    applicable_indicators UUID[] NOT NULL DEFAULT ARRAY[]::UUID[],
    affected_groups      TEXT[] NOT NULL DEFAULT ARRAY[]::TEXT[],
    assessment_summary   TEXT,
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (place_id, period_start, period_end)
);

CREATE INDEX place_assessments_place_assessed_idx
    ON place_assessments (place_id, assessed_at DESC);

-- Generated content: plain-language update text
CREATE TABLE place_content (
    id                   UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    place_id             UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    generated_at         TIMESTAMPTZ NOT NULL,
    period_start         TIMESTAMPTZ NOT NULL,
    period_end           TIMESTAMPTZ NOT NULL,
    content_type         TEXT NOT NULL CHECK (content_type IN ('today', 'forecast', 'alert')),
    language             TEXT NOT NULL DEFAULT 'en',
    headline             TEXT NOT NULL,
    body_text            TEXT NOT NULL,
    call_to_action       TEXT,
    source_indicators    UUID[] NOT NULL DEFAULT ARRAY[]::UUID[],
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (place_id, content_type, period_start, period_end, language)
);

CREATE INDEX place_content_place_generated_idx
    ON place_content (place_id, generated_at DESC);
CREATE INDEX place_content_type_idx
    ON place_content (content_type, generated_at DESC);

-- Seed common indicator definitions
INSERT INTO indicator_definitions (code, display_name, category, source_variables, description)
VALUES
    ('rain_intensity_trend', 'Rain Intensity Trend', 'weather', ARRAY['precipitation_mm'], 'Trend in precipitation over the observation period'),
    ('rain_frequency', 'Rain Frequency', 'weather', ARRAY['precipitation_mm'], 'Number of rainy days in the observation period'),
    ('temperature_trend', 'Temperature Trend', 'weather', ARRAY['temperature_celsius'], 'Mean temperature trend across the period'),
    ('wind_speed_trend', 'Wind Speed Trend', 'weather', ARRAY['wind_speed_ms'], 'Wind speed trend across the period'),
    ('soil_moisture_level', 'Soil Moisture Level', 'agriculture', ARRAY['soil_moisture_m3m3'], 'Current soil moisture as percentage of capacity'),
    ('drought_stress', 'Drought Stress Index', 'agriculture', ARRAY['precipitation_mm', 'temperature_celsius', 'soil_moisture_m3m3'], 'Composite drought risk indicator'),
    ('water_quality_concern', 'Water Quality Concern', 'water', ARRAY[]::TEXT[], 'Aggregate water quality assessment (stubbed: requires water-specific data)'),
    ('air_quality_index', 'Air Quality Index (AQI)', 'air_quality', ARRAY[]::TEXT[], 'Composite air quality from PM2.5 and Ozone (stubbed: requires air data)');

-- Default applicability rules (can be overridden per-place)
INSERT INTO indicator_applicability_rules (indicator_id, place_type, applicable, relevance_score)
SELECT id, 'agricultural', TRUE, 1.0 FROM indicator_definitions WHERE code = 'rain_intensity_trend'
UNION ALL
SELECT id, 'agricultural', TRUE, 1.0 FROM indicator_definitions WHERE code = 'soil_moisture_level'
UNION ALL
SELECT id, 'agricultural', TRUE, 1.0 FROM indicator_definitions WHERE code = 'drought_stress'
UNION ALL
SELECT id, 'urban', TRUE, 0.8 FROM indicator_definitions WHERE code = 'rain_intensity_trend'
UNION ALL
SELECT id, 'water_community', TRUE, 1.0 FROM indicator_definitions WHERE code = 'water_quality_concern'
UNION ALL
SELECT id, 'water_community', TRUE, 0.9 FROM indicator_definitions WHERE code = 'rain_intensity_trend';
