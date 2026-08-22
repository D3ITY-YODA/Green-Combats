# ── Root Dockerfile — Green Compass full stack ─────────────────────
# Builds both backend and frontend in a single image.
# Backend runs on :8080 (internal), frontend on :3000 (exposed).
# Render routes external traffic to :3000.

# ══════════════════════════════════════════════════════════════════
# Stage 1: Build Go backend
# ══════════════════════════════════════════════════════════════════
FROM golang:1.26-alpine AS backend-builder
ARG VERSION=dev
WORKDIR /src

COPY green-compass-backend/go.mod green-compass-backend/go.sum ./
RUN go mod download

COPY green-compass-backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags "-s -w -X main.version=${VERSION}" \
    -o /out/api ./cmd/api

RUN CGO_ENABLED=0 GOOS=linux go build -trimpath \
    -ldflags "-s -w" \
    -o /out/migrate ./cmd/migrate

# ══════════════════════════════════════════════════════════════════
# Stage 2: Build Next.js frontend
# ══════════════════════════════════════════════════════════════════
FROM node:20-alpine AS frontend-deps
WORKDIR /app
COPY frontend-web/package.json frontend-web/package-lock.json ./
RUN npm ci --ignore-scripts

FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY --from=frontend-deps /app/node_modules ./node_modules
COPY frontend-web/ ./
ENV NEXT_TELEMETRY_DISABLED=1
ENV NEXT_PUBLIC_API_URL=http://localhost:8080
ENV NEXT_PUBLIC_APP_NAME="Green Compass"
ENV NEXT_PUBLIC_DEFAULT_LOCALE=en
RUN npm run build

# ══════════════════════════════════════════════════════════════════
# Stage 3: Final production image
# ══════════════════════════════════════════════════════════════════
FROM node:20-alpine AS runner
RUN apk add --no-cache bash curl

# Create non-root user
RUN addgroup --system --gid 1001 nodejs \
 && adduser  --system --uid 1001 nextjs

# ── Copy backend binaries ─────────────────────────────────────────
COPY --from=backend-builder /out/api      /usr/local/bin/api
COPY --from=backend-builder /out/migrate  /usr/local/bin/migrate
COPY green-compass-backend/configs /app/configs

# ── Copy frontend (standalone) ────────────────────────────────────
WORKDIR /app
COPY --from=frontend-builder --chown=nextjs:nodejs /app/public          ./public
COPY --from=frontend-builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=frontend-builder --chown=nextjs:nodejs /app/.next/static    ./.next/static

# ── Startup script ────────────────────────────────────────────────
COPY startup.sh /app/startup.sh
RUN chmod +x /app/startup.sh

USER nextjs
EXPOSE 3000

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

ENTRYPOINT ["/app/startup.sh"]
