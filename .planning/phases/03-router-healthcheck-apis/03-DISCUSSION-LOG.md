# Phase 3: Router & Healthcheck APIs - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-06-22
**Phase:** 3-Router & Healthcheck APIs
**Areas discussed:** Router Scaffolding Location, Gin Engine Mode, Healthcheck Payload Format

---

## Router Scaffolding Location

| Option | Description | Selected |
|--------|-------------|----------|
| Under `internal/app/` | Scaffolds router registration under application bootstrapper package | ✓ |
| Under `internal/pkg/router/` | Scaffolds router inside shared packages | |

**User's choice:** Under `internal/app/`
**Notes:** Keeps `internal/pkg` strictly for reusable utility modules, separating routing logic.

---

## Gin Engine Mode (Debug vs Release)

| Option | Description | Selected |
|--------|-------------|----------|
| Dynamic | Selects mode on the fly based on ENV config | ✓ |
| Hardcoded | Debug local, release remote, statically compiled | |

**User's choice:** Dynamic
**Notes:** Allows configuration changes without re-compilation.

---

## Healthcheck (`/healthz`) Payload Format

| Option | Description | Selected |
|--------|-------------|----------|
| Detailed JSON | Returns status object of Postgres and Redis check | ✓ |
| Simple Text | Returns plain text 'OK' | |

**User's choice:** Detailed JSON
**Notes:** Allows monitoring tools to quickly identify which specific dependency is down.

---

## the agent's Discretion
- Logging format details, CORS middlewares, endpoint path details, database/cache check timeout limits.

## Deferred Ideas
None.
