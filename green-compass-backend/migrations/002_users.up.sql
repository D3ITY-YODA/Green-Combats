CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone_number  TEXT,
    email         CITEXT,
    password_hash TEXT,
    display_name  TEXT NOT NULL,
    language      TEXT NOT NULL DEFAULT 'en',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT users_contact_required CHECK (email IS NOT NULL OR phone_number IS NOT NULL)
);

CREATE UNIQUE INDEX users_phone_number_key ON users (phone_number) WHERE phone_number IS NOT NULL;
CREATE UNIQUE INDEX users_email_key ON users (email) WHERE email IS NOT NULL;

CREATE TRIGGER users_set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();
