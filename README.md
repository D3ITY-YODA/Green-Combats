# Green Compass

**Know Your Place. Move with Change.**

Green Compass is a location-aware climate and environmental information platform that delivers clear, actionable environmental updates for the places that matter to you. It combines a mobile-first user experience with a robust backend infrastructure and an institutional admin console.

---

## Overview

Green Compass helps communities stay informed about local environmental conditions including:

- **Weather & Climate Updates** — Current conditions, forecasts, and alerts
- **Water Quality** — River levels, contamination reports, and advisories
- **Air Quality** — Pollution indices and health recommendations
- **Agricultural Reports** — Crop damage, drought conditions, and pest alerts
- **Flood Warnings** — Real-time flood risk assessments and notifications
- **Community Observations** — Crowdsourced environmental reports from local users

The platform supports **multi-language localization**, **offline capabilities**, and **multiple notification channels** (push, SMS, voice) to ensure critical information reaches users regardless of connectivity.

---

## Project Structure

```
green-compass/
├── mobile/                  # Android mobile application
├── frontend-web/            # Next.js admin console & marketing site
├── green-compass-backend/   # Go backend API server
└── README.md
```

---

## Mobile App (`mobile/`)

An Android application built with Kotlin and Jetpack Compose, designed for community members to receive environmental updates and submit observations.

### Tech Stack

| Technology | Purpose |
|------------|---------|
| Kotlin | Primary language |
| Jetpack Compose | Modern declarative UI |
| Hilt | Dependency injection |
| Room | Local database (offline support) |
| Retrofit + OkHttp | HTTP client for API communication |
| WorkManager | Background sync tasks |
| DataStore | Preferences storage |
| Security Crypto | Encrypted credential storage |
| MapLibre | Map rendering |
| Coil | Image loading |
| Coroutines | Async operations |

### Features

- **Onboarding Flow** — Language selection, interests, place setup, permissions
- **Today Dashboard** — Daily environmental summary with weather, alerts, and conditions
- **Explore** — Discover environmental data by location
- **Updates** — Real-time environmental alerts and notifications
- **Reports** — Submit observations (flood, drought, water quality, air quality, crop damage)
- **Profile & Settings** — Notification preferences, language, units, theme
- **Saved Places** — Manage favorite locations
- **Location Switcher** — Quickly change active location
- **Offline Mode** — View cached data without connectivity

### Screens

| Screen | Description |
|--------|-------------|
| Welcome | App introduction |
| Language Selection | Choose preferred language |
| Account Choice | Sign in or create account |
| Sign In / Up | Authentication |
| Place Setup | Configure primary location |
| Interests | Select environmental topics |
| Permission | Location & notification permissions |
| Today | Daily environmental summary |
| Explore | Location-based data discovery |
| Updates | Real-time alerts feed |
| Report Form | Submit environmental observation |
| Profile | User settings and preferences |
| Settings | App configuration |
| Notifications | Notification channel preferences |

---

## Frontend Web (`frontend-web/`)

A Next.js application serving as both a marketing site and an institutional admin console.

### Tech Stack

| Technology | Purpose |
|------------|---------|
| Next.js 16 | React framework with App Router |
| React 19 | UI library |
| TypeScript | Type-safe development |
| Tailwind CSS 4 | Utility-first styling |
| MapLibre GL | Interactive maps |
| React Query | Server state management |
| React Hook Form | Form handling |
| Zod | Schema validation |
| Lucide React | Icon library |
| Playwright | End-to-end testing |
| Vitest | Unit testing |

### Routes

| Route | Description |
|-------|-------------|
| `/` | Marketing landing page |
| `/sign-in` | User authentication |
| `/sign-up` | Account registration |
| `/console` | Admin console overview |
| `/console/local-conditions` | Local environmental conditions |
| `/console/updates/new` | Create new update |
| `/report` | Submit report |
| `/updates` | View all updates |
| `/profile` | User profile management |
| `/privacy` | Privacy policy |
| `/terms` | Terms of service |

### Admin Console

The console provides institutional users with:

- **Overview Dashboard** — Key metrics (important updates, community reports, delayed sources, pending reviews)
- **Priority Actions** — Urgent items requiring attention
- **Local Conditions** — Detailed environmental data view
- **Update Management** — Create and manage environmental updates

---

## Backend API (`green-compass-backend/`)

A Go-based REST API server handling authentication, data processing, and business logic.

### Tech Stack

| Technology | Purpose |
|------------|---------|
| Go 1.24 | Primary language |
| Gin | HTTP framework |
| PostgreSQL + PostGIS | Database with geospatial support |
| pgx | PostgreSQL driver |
| JWT (golang-jwt) | Authentication tokens |
| Golang Migrate | Database migrations |
| Docker | Containerization |

### Architecture

The backend follows a clean architecture pattern:

```
cmd/           # Application entry points
├── api/       # Main API server
├── migrate/   # Database migration tool
├── worker/    # Background job worker
├── notifier/  # Notification dispatcher
└── scheduler/ # Task scheduler

internal/      # Domain-specific business logic
├── auth/      # Authentication & authorization
├── users/     # User management
├── places/    # Location management
├── updates/   # Environmental updates
├── observations/  # Community observations
├── reports/   # Report generation
├── notifications/ # Push/SMS/voice notifications
├── organizations/ # Organization management
├── projects/  # Project management
├── indicators/ # Environmental indicators
├── assessment/ # Risk assessment
├── preferences/ # User preferences
├── context/   # Request context handling
├── audit/     # Audit logging
├── permissions/ # Permission management
├── health/    # Health check endpoints
├── reporting/ # Analytics & reporting
├── ingestion/ # Data ingestion
├── normalization/ # Data normalization
├── sources/   # Data source management
├── content/   # Content management
└── config/    # Configuration loading

pkg/           # Shared utilities
├── database/  # Database connection pool
├── cache/     # Caching layer
├── storage/   # File storage
├── geo/       # Geospatial calculations
├── logging/   # Structured logging
├── clock/     # Time utilities
├── httpx/     # HTTP utilities
├── security/  # Security helpers
├── idempotency/ # Request idempotency
├── broker/    # Message broker
└── telemetry/ # Observability
```

### API Endpoints

#### Authentication (`/v1/auth`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/register` | Create new account |
| POST | `/login` | Authenticate user |
| POST | `/refresh` | Refresh access token |
| POST | `/logout` | Invalidate session |
| GET | `/me` | Get current user |

#### Places (`/v1/places`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | List places |
| POST | `/` | Create place |
| GET | `/:id` | Get place details |
| PUT | `/:id` | Update place |
| DELETE | `/:id` | Delete place |

#### Observations (`/v1/observations`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/` | Submit observation |
| GET | `/` | List my observations |
| GET | `/:id` | Get observation details |
| GET | `/institutional/pending` | List pending observations (admin) |
| PATCH | `/institutional/:id/verify` | Verify observation (admin) |

#### Updates (`/v1/updates`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | List updates |
| GET | `/:id` | Get update details |

#### Notifications (`/v1/notifications`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | List notifications |
| PATCH | `/:id/read` | Mark as read |

#### Reports (`/v1/reports`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/` | Submit report |
| GET | `/` | List reports |

#### Context (`/v1/context`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Get user context (place, preferences, updates) |

#### Preferences (`/v1/preferences`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Get user preferences |
| PUT | `/` | Update preferences |

#### Health (`/v1/health`)
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Service health check |

### Database Schema

The PostgreSQL database uses PostGIS for geospatial queries with the following core tables:

- **users** — User accounts with contact info and preferences
- **places** — Named locations with geospatial coordinates
- **observations** — Community-submitted environmental reports
- **updates** — Official environmental updates and alerts
- **notifications** — User notification records
- **organizations** — Institutional users and groups
- **projects** — Environmental monitoring projects
- **indicators** — Environmental measurement definitions
- **reports** — Generated analytics reports
- **audit_logs** — System audit trail

---

## Getting Started

### Prerequisites

- **Mobile**: Android Studio, JDK 17, Android SDK 34
- **Frontend**: Node.js 18+, npm
- **Backend**: Go 1.24+, Docker (for PostgreSQL + PostGIS)

### Backend Setup

```bash
cd green-compass-backend

# Copy environment file
cp .env.example .env

# Start PostgreSQL with PostGIS
make dev-up

# Run migrations
make migrate-up

# Start the API server
make run
```

### Frontend Setup

```bash
cd frontend-web

# Copy environment file
cp .env.example .env.local

# Install dependencies
npm install

# Start development server
npm run dev
```

The app will be available at `http://localhost:3000`.

### Mobile Setup

```bash
cd mobile

# Open in Android Studio
# Or build from command line
./gradlew assembleDebug
```

---

## Development Commands

### Backend

| Command | Description |
|---------|-------------|
| `make dev-up` | Start PostgreSQL dev stack |
| `make dev-down` | Stop PostgreSQL dev stack |
| `make run` | Start API server |
| `make test` | Run unit tests |
| `make test-integration` | Run integration tests |
| `make migrate-up` | Apply migrations |
| `make migrate-down` | Rollback migrations |
| `make migrate-status` | Check migration status |
| `make lint` | Run linter |
| `make fmt` | Format code |
| `make build` | Build binary |

### Frontend

| Command | Description |
|---------|-------------|
| `npm run dev` | Start dev server |
| `npm run build` | Production build |
| `npm run start` | Start production server |
| `npm run lint` | Run linter |

### Mobile

| Command | Description |
|---------|-------------|
| `./gradlew assembleDebug` | Build debug APK |
| `./gradlew assembleRelease` | Build release APK |
| `./gradlew test` | Run unit tests |

---

## Docker

### Backend

```bash
cd green-compass-backend

# Build image
docker build -t green-compass-api .

# Run container
docker run -p 8080:8080 green-compass-api
```

### Docker Compose (Development)

The `make dev-up` command starts a PostgreSQL + PostGIS container for local development.

---

## Authentication

Green Compass uses JWT-based authentication:

- **Access Token**: Short-lived (15 minutes), used for API requests
- **Refresh Token**: Long-lived (30 days), used to obtain new access tokens

Tokens are included in the `Authorization` header as `Bearer <token>`.

---

## Localization

The platform supports multiple languages:

- **Backend**: Language field stored per user, localized content delivery
- **Mobile**: Language selection during onboarding, supports English and Spanish
- **Frontend**: Configurable default locale via `NEXT_PUBLIC_DEFAULT_LOCALE`

---

## Geospatial Features

- **PostGIS**: All proximity searches and radius queries use PostGIS (`ST_DWithin`, `ST_Distance`)
- **GEOGRAPHY columns**: Store coordinates in WGS84 (SRID 4326)
- **GiST indexes**: Optimized spatial queries for location-based searches
- **Coordinate validation**: `pkg/geo` handles non-DB-bound calculations

---

## Environment Variables

### Backend (`.env`)

| Variable | Description | Default |
|----------|-------------|---------|
| `GC_APP_ENV` | Environment (local/dev/staging/production) | `local` |
| `GC_SERVER_PORT` | HTTP server port | `8080` |
| `GC_DATABASE_URL` | PostgreSQL connection string | — |
| `GC_AUTH_SECRET` | JWT signing secret | — |
| `GC_AUTH_ACCESS_TTL` | Access token TTL | `15m` |
| `GC_AUTH_REFRESH_TTL` | Refresh token TTL | `720h` |
| `GC_LOG_LEVEL` | Log level (debug/info/warn/error) | `debug` |
| `GC_LOG_FORMAT` | Log format (text/json) | `text` |

### Frontend (`.env.local`)

| Variable | Description | Default |
|----------|-------------|---------|
| `NEXT_PUBLIC_API_URL` | Backend API URL | `http://localhost:8080` |
| `NEXT_PUBLIC_MAP_STYLE_URL` | MapLibre style URL | — |
| `NEXT_PUBLIC_DEFAULT_LOCALE` | Default language | `en` |
| `NEXT_PUBLIC_APP_NAME` | Application name | `Green Compass` |

---

## Testing

### Backend

```bash
# Unit tests
make test

# Integration tests (requires running database)
make test-integration

# Full CI test suite
make ci-test
```

### Frontend

```bash
# Unit tests (Vitest)
npm test

# E2E tests (Playwright)
npx playwright test
```

---

## Deployment

### Backend

- **Docker**: Build and deploy the container image
- **Binary**: Cross-compile with `make build` for target architecture
- **Database**: Run migrations before deployment

### Frontend

- **Vercel**: Recommended for Next.js deployment
- **Docker**: Can be containerized for self-hosted deployment
- **Static Export**: Available for marketing pages

### Render (recommended all-in-one)

A `render.yaml` Blueprint is provided at the repository root that deploys the
full stack on [Render](https://render.com):

1. **`green-compass-db`** — a PostGIS private service (`postgis/postgis:16-3.4`).
   Render's managed Postgres does **not** include the PostGIS extension this
   backend requires, so PostGIS is run as its own internal service.
2. **`green-compass-api`** — the Go API, built from
   `green-compass-backend/Dockerfile`. It reads `PORT` (injected by Render),
   applies database migrations automatically on startup (`GC_RUN_MIGRATIONS=true`),
   and enables CORS for the frontend.
3. **`green-compass-web`** — the Next.js frontend, configured with
   `NEXT_PUBLIC_API_URL` pointing at the API service so the browser can call it.

To deploy:

```bash
# Connect the repo to Render and create the Blueprint from render.yaml,
# or run the Render CLI from the repo root:
render blueprint launch
```

Key environment variables (all wired in `render.yaml`):

| Service | Variable | Purpose |
|---------|----------|---------|
| api | `GC_DATABASE_URL` | Built from the db service's generated password/host |
| api | `GC_RUN_MIGRATIONS` | Applies migrations on startup (`true`) |
| api | `GC_CORS_ALLOWED_ORIGINS` | Browser origins allowed to call the API |
| api | `GC_AUTH_SECRET` | JWT signing secret (auto-generated) |
| web | `NEXT_PUBLIC_API_URL` | Public URL of the API service |

For production, restrict `GC_CORS_ALLOWED_ORIGINS` to your frontend's exact
origin and add a disk to `green-compass-db` if you need data to survive rebuilds.

### Mobile

- **Google Play Store**: Build signed release APK
- **Internal Distribution**: Share debug APK for testing

---

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## License

This project is proprietary software. All rights reserved.

---

## Support

For support and inquiries, please contact the development team.

---

**Built with care for people and the environment** 
