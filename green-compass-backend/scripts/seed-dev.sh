#!/usr/bin/env bash
set -euo pipefail

# Seed development data into the local Postgres database.
# Usage: make seed-dev
# Requires: psql (PostgreSQL client)

DSN="${GC_MIGRATE_DSN:-postgres://greencompass:greencompass@127.0.0.1:5432/greencompass?sslmode=disable}"

echo "Seeding dev data into: $DSN"

psql "$DSN" <<-'SQL'
BEGIN;

-- Users (password: "password123" — bcrypt hash)
INSERT INTO users (phone_number, email, display_name, password_hash, language) VALUES
    ('+256700000001', 'admin@greencompass.org', 'Admin User', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'en'),
    ('+256700000002', 'reviewer@water.go.ug', 'Water Authority Reviewer', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'en'),
    ('+256700000003', 'farmer@local.ug', 'Joseph the Farmer', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'sw'),
    ('+256700000004', 'community@local.ug', 'Community Leader', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'en'),
    ('+256700000005', 'sms@local.ug', 'SMS Only User', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'sw')
ON CONFLICT DO NOTHING;

-- Mark first user as platform admin
UPDATE users SET is_platform_admin = TRUE
WHERE email = 'admin@greencompass.org';

-- Organizations
INSERT INTO organizations (name, status) VALUES
    ('Ministry of Water and Environment', 'verified'),
    ('Kampala Capital City Authority', 'verified'),
    ('Local Community Network', 'pending')
ON CONFLICT DO NOTHING;

-- Places (using PostGIS ST_SetSRID + ST_MakePoint for lat/lng)
INSERT INTO places (name, place_type, location) VALUES
    ('Kampala', 'ward', ST_SetSRID(ST_MakePoint(32.5825, 0.3476), 4326)),
    ('Entebbe', 'ward', ST_SetSRID(ST_MakePoint(32.4637, 0.0562), 4326)),
    ('Mukono', 'district', ST_SetSRID(ST_MakePoint(32.9283, 0.3533), 4326)),
    ('Jinja', 'district', ST_SetSRID(ST_MakePoint(33.2041, 0.4478), 4326)),
    ('Gulu', 'district', ST_SetSRID(ST_MakePoint(32.2979, 2.7747), 4326))
ON CONFLICT DO NOTHING;

-- Save primary places for users
INSERT INTO user_saved_places (user_id, place_id, is_primary)
SELECT u.id, p.id, TRUE
FROM users u, places p
WHERE u.email = 'farmer@local.ug' AND p.name = 'Kampala'
ON CONFLICT DO NOTHING;

INSERT INTO user_saved_places (user_id, place_id, is_primary)
SELECT u.id, p.id, TRUE
FROM users u, places p
WHERE u.email = 'community@local.ug' AND p.name = 'Jinja'
ON CONFLICT DO NOTHING;

-- Notification preferences
INSERT INTO notification_preferences (user_id, channel, event_type, enabled)
SELECT u.id, 'push', 'content.generated', TRUE
FROM users u WHERE u.email = 'farmer@local.ug'
ON CONFLICT DO NOTHING;

INSERT INTO notification_preferences (user_id, channel, event_type, enabled)
SELECT u.id, 'sms', 'content.generated', TRUE
FROM users u WHERE u.email = 'sms@local.ug'
ON CONFLICT DO NOTHING;

COMMIT;

SQL

echo "✓ Dev data seeded successfully"
echo ""
echo "Test accounts:"
echo "  admin@greencompass.org / password123 (platform admin)"
echo "  reviewer@water.go.ug / password123  (institutional reviewer)"
echo "  farmer@local.ug / password123       (community farmer)"
echo "  community@local.ug / password123    (community leader)"
echo "  sms@local.ug / password123          (SMS-only user)"
