# Phase 2: Directory Layout & Connections - Context

**Gathered:** 2026-06-22
**Status:** Ready for planning

<domain>
## Phase Boundary

Scaffold project layout following Go Standard Layout (`cmd/`, `internal/app/`, `internal/pkg/`) and configure PostgreSQL (with golang-migrate) & Redis connections with configuration management.

</domain>

<decisions>
## Implementation Decisions

### Configuration Management
- **D-01:** Implement a configuration package under `internal/pkg/config/config.go` containing a `Config` struct.
- **D-02:** Use `github.com/joho/godotenv` to automatically load `.env` environment variables when running in local development mode.

### Directory Structure & Layout
- **D-03:** Apply Go Standard Layout:
  - `cmd/server/main.go` — Entrypoint.
  - `internal/app/` — Application bootstrapper.
  - `internal/pkg/` — Reusable shared packages.
- **D-04:** Place PostgreSQL connection under `internal/pkg/db/postgres.go` and Redis connection under `internal/pkg/cache/redis.go`.

### Database & Migrations
- **D-05:** Database connection pool is initialized using Gorm with default Postgres driver.
- **D-06:** Database migrations will be structured SQL files located under `./migrations/` directory.

### Logging
- **D-07:** Use Go's built-in structured logger `slog` (introduced in Go 1.21) for logging events.

### the agent's Discretion
- The agent has discretion over the exact implementation of Gorm config, Redis connection options, logging format tweaks, and configuration struct fields.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Project Context & Specifications
- `.planning/PROJECT.md` — Core value, constraints, and constraints.
- `.planning/REQUIREMENTS.md` — Scoped v1 requirements.
- `.planning/ROADMAP.md` — Phase definition and success criteria.
- `.planning/phases/01-environment-tooling/01-01-SUMMARY.md` — Setup information for dependency packages.
- `.planning/phases/01-environment-tooling/01-02-SUMMARY.md` — Docker compose and Air settings.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `go.mod`, `go.sum` — Loaded packages.
- `docker-compose.yml`, `.env.example` — Base parameters for Postgres/Redis services.

### Established Patterns
- Defined standard modular layout in D-03.

### Integration Points
- This phase builds the Go application scaffolding. The next phase will run the Gin router on top of this scaffolding.

</code_context>

<specifics>
- Use `github.com/joho/godotenv` to load `.env` variables for local config.

</specifics>

<deferred>
- None.
</deferred>

---

*Phase: 2-Directory Layout & Connections*
*Context gathered: 2026-06-22*
