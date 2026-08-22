#!/bin/bash
set -e

echo "============================================"
echo "  Green Compass — Starting up"
echo "============================================"

# ── Run database migrations ───────────────────────────────────────
if [ -n "$GC_DATABASE_URL" ]; then
  echo "[1/3] Running database migrations..."
  /usr/local/bin/migrate up 2>&1 || echo "[1/3] Migration warning (non-fatal, continuing...)"
else
  echo "[1/3] No GC_DATABASE_URL set, skipping migrations"
fi

# ── Start backend in background ───────────────────────────────────
echo "[2/3] Starting backend API on :8080..."
/usr/local/bin/api &
BACKEND_PID=$!

# Wait for backend to be ready
echo "      Waiting for backend to be ready..."
for i in $(seq 1 30); do
  if curl -sf http://localhost:8080/healthz > /dev/null 2>&1; then
    echo "      Backend is ready!"
    break
  fi
  if [ $i -eq 30 ]; then
    echo "      Backend health check timed out, continuing anyway..."
  fi
  sleep 1
done

# ── Start frontend (foreground) ───────────────────────────────────
echo "[3/3] Starting frontend on :3000..."
exec node server.js
