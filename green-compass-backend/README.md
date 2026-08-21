# GreenCompass Backend

Go backend for location-aware climate and environmental information.

## Development

| Command | Purpose |
| --- | --- |
| `make dev-up` / `make dev-down` | PostGIS dev stack (host port `GC_DEV_PG_PORT`, default 5432 — see `.env.example`) |
| `make test` | Unit tests with `-race` |
| `make test-integration` | DB-backed tests (`Integration`-named; self-skip without `GC_TEST_DATABASE_URL`) |
| `make ci-test` | Full suite with the integration database forced on — what CI must run |
| `make migrate-up` / `migrate-down` / `migrate-status` | Schema migrations via `cmd/migrate` (`MIGRATE_STEPS=2` for partial) |
| `make run` / `make build` / `make lint` / `make fmt` / `make vet` | Everyday targets |

## Architecture notes

- **Migration runner**: hand-rolled (`internal/migrate`) instead of golang-migrate/goose. Keeps versioned SQL embedded in the module with zero extra dependencies, enforces strict `NNN_name.up/down` pairing at parse time, and serializes concurrent runners through a Postgres advisory lock with a `schema_migrations` ledger table.
- **Distance/proximity policy**: `pkg/geo` is for non-DB-bound calculations only (coordinate validation, in-memory checks). All proximity search and radius queries go through PostGIS (`ST_DWithin`/`ST_Distance` on `GEOGRAPHY` columns) inside repository SQL so the GiST indexes are actually used. Repositories must not filter or sort by `pkg/geo.MetersBetween`.
