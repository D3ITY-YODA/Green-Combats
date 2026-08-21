-- Notifications: delivery preferences, notification records, and delivery log

-- Per-user notification preferences (which channels, which event types)
CREATE TABLE notification_preferences (
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel             TEXT NOT NULL CHECK (channel IN ('push', 'sms', 'ussd', 'email')),
    event_type          TEXT NOT NULL,
    enabled             BOOLEAN NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id, channel, event_type)
);

CREATE TRIGGER notification_preferences_set_updated_at
    BEFORE UPDATE ON notification_preferences
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- Notification records: one row per notification to be delivered
CREATE TABLE notifications (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    channel             TEXT NOT NULL CHECK (channel IN ('push', 'sms', 'ussd', 'email')),
    event_type          TEXT NOT NULL,
    title               TEXT NOT NULL,
    body                TEXT NOT NULL,
    metadata            JSONB,
    status              TEXT NOT NULL CHECK (status IN ('pending', 'sent', 'failed', 'cancelled')),
    created_at          TIMESTAMPTZ NOT NULL DEFAULT now(),
    sent_at             TIMESTAMPTZ,
    error_message       TEXT
);

CREATE INDEX notifications_status_created_idx
    ON notifications (status, created_at);
CREATE INDEX notifications_user_created_idx
    ON notifications (user_id, created_at DESC);

-- Delivery log: tracks each delivery attempt for observability
CREATE TABLE notification_delivery_log (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_id     UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    attempt             INTEGER NOT NULL CHECK (attempt > 0),
    status              TEXT NOT NULL CHECK (status IN ('sent', 'failed')),
    provider_response   TEXT,
    attempted_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (notification_id, attempt)
);
