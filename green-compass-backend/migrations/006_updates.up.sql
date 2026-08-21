CREATE TABLE updates (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    place_id       UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    indicator      TEXT NOT NULL,
    severity       TEXT NOT NULL CHECK (severity IN ('info', 'watch', 'warning')),
    title          TEXT NOT NULL,
    body           TEXT NOT NULL,
    effective_from TIMESTAMPTZ NOT NULL,
    expires_at     TIMESTAMPTZ,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT updates_expiry_after_effective CHECK (expires_at IS NULL OR expires_at > effective_from)
);

CREATE INDEX updates_place_effective_idx ON updates (place_id, effective_from DESC);
