#!/bin/sh
set -e

echo "greencompass: running database migrations..."
/usr/local/bin/migrate up 2>&1 || echo "greencompass: migration warning (non-fatal)"

echo "greencompass: starting API server..."
exec /usr/local/bin/api
