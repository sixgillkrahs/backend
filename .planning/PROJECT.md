# Golang Gin Backend Scaffold

## What This Is

A Golang backend service scaffold using the Gin framework, integrated with PostgreSQL and Redis. It is structured using the Go Standard Layout (`cmd/`, `internal/app/`, `internal/pkg/`) and configured with hot reloading using Air, database migrations using golang-migrate, and local service orchestration via docker-compose.

## Core Value

Provide a clean, production-ready, containerized Go backend structure that makes starting development of new endpoints and services immediate and seamless.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] Initialize the Go module and install dependencies (Gin, Gorm/pg/redis client, etc.).
- [ ] Configure `docker-compose.yml` for PostgreSQL, Redis, and hot-reload Air.
- [ ] Setup Air config file (`.air.toml`) for live reloading of the Go application.
- [ ] Configure basic application structure following Go Standard Layout (`cmd/main.go`, `internal/app/...`, `internal/pkg/...`).
- [ ] Configure PostgreSQL database connection and database migration workflow using `golang-migrate`.
- [ ] Configure Redis client connection.
- [ ] Create a sample Ping/Pong or Healthcheck API route to verify connections and hot reloading.

### Out of Scope

- Production deployment configurations (Kubernetes, AWS ECS, etc.) — Deferred to operational phases.
- Real business logic controllers — Out of scope for a basic bootstrap template.

## Context

The developer wants a standard, scalable starting point for a Golang backend. Using Gin for high performance, Gorm or pg client for DB access, Redis for caching/caching utility, and Air for a fast local developer feedback loop.

## Constraints

- **Language**: Go (v1.22+ recommended) — Core language requirement.
- **Framework**: Gin-Gonic — Core HTTP web framework.
- **Data Stores**: PostgreSQL & Redis — Required state and caching systems.
- **Hot Reload**: Air — Tool for hot reloading.
- **Orchestration**: Docker Compose — Local development orchestration.
- **Migrations**: golang-migrate — Core migration utility.

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Go Standard Layout | Scalable, clean separation of concerns for large-scale Go apps | — Pending |
| Docker-compose local env | Standardizes PostgreSQL and Redis environment for all devs | — Pending |
| golang-migrate | Clean SQL-based migrations, independent of application code | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-06-22 after initialization*
