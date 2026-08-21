ALTER TABLE users
    ADD COLUMN phone_verified_at TIMESTAMPTZ,
    ADD COLUMN email_verified_at TIMESTAMPTZ,
    ADD COLUMN is_platform_admin BOOLEAN NOT NULL DEFAULT FALSE;

ALTER TABLE organizations
    ADD COLUMN status TEXT NOT NULL DEFAULT 'pending'
    CONSTRAINT organizations_status_check CHECK (status IN ('pending', 'verified', 'rejected'));

CREATE TABLE refresh_tokens (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash  TEXT NOT NULL UNIQUE,
    issued_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at  TIMESTAMPTZ NOT NULL,
    revoked_at  TIMESTAMPTZ,
    replaced_by UUID REFERENCES refresh_tokens(id)
);

CREATE INDEX refresh_tokens_user_id_idx ON refresh_tokens (user_id);
