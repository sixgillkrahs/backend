# Phase 3: Router & Healthcheck APIs - Context

**Gathered:** 2026-06-22
**Status:** Ready for planning

<domain>
## Phase Boundary

Configure Gin router and verify status with `/ping` and `/healthz` API endpoints. The `/healthz` endpoint must dynamically check both the PostgreSQL database and Redis connections and return a detailed status report.

</domain>

<decisions>
## Implementation Decisions

### Router Organization
- **D-08:** Setup Gin router and route registration inside `internal/app/router.go` or `internal/app/app.go` as part of the application bootstrapper.
- **D-09:** Set the Gin running mode dynamically depending on environment configuration (if `ENV=production`, set Release mode; otherwise, Debug mode).

### API Endpoints
- **D-10:** `/healthz` returns a detailed JSON response showing the status of each dependent connection pool (PostgreSQL and Redis) with appropriate HTTP response codes (200 OK if all up, 500 Internal Server Error if any down).
- **D-11:** `/ping` returns a simple `{"message": "pong"}` JSON payload.

### Logging & Instrumentation
- **D-12:** Use Go's built-in `slog` structured logger to log incoming API requests and startup details.

### the agent's Discretion
- The agent has discretion over the exact routing middleware (e.g., recovery, logger), CORS setup (if any), database/cache ping timeout thresholds, and port binding logic.

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Project Context & Specifications
- `.planning/PROJECT.md` — Core value, constraints, and constraints.
- `.planning/REQUIREMENTS.md` — Scoped v1 requirements.
- `.planning/ROADMAP.md` — Phase definition and success criteria.
- `.planning/phases/02-directory-layout-connections/02-01-SUMMARY.md` — Layout scaffolding details.
- `.planning/phases/02-directory-layout-connections/02-02-SUMMARY.md` — Database setup details.
- `.planning/phases/02-directory-layout-connections/02-03-SUMMARY.md` — Cache setup details.

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/pkg/config` — Provides config variables (port, environment, etc.).
- `internal/pkg/db` — Provides Postgres pool initialization (`InitPostgres`).
- `internal/pkg/cache` — Provides Redis pool initialization (`InitRedis`).
- `cmd/server/main.go` — Entrypoint where the router and server will be initialized.

### Established Patterns
- Modular initialization and structured logging with `slog` at startup.

### Integration Points
- This phase wires the Gin HTTP engine on top of the connection pools initialized in Phase 2.

</code_context>

<specifics>
## Specific Ideas
- Dynamic check in `/healthz` will perform a `Ping` query on the GORM/Postgres DB object and a `Ping` check on the Redis client.

</specifics>

<deferred>
## Deferred Ideas
- None.
</deferred>

---

*Phase: 3-Router & Healthcheck APIs*
*Context gathered: 2026-06-22*
