# Golang Gin Backend Scaffold

## What This Is

A Golang backend service scaffold using the Gin framework, integrated with PostgreSQL and Redis. It is structured using the Go Standard Layout (`cmd/`, `internal/app/`, `internal/pkg/`) and configured with hot reloading using Air, database migrations using golang-migrate, and local service orchestration via docker-compose.

## Core Value

Provide a clean, production-ready, containerized Go backend structure that makes starting development of new endpoints and services immediate and seamless.

## Current Milestone: v1.2 Swagger API Documentation

**Goal:** Tích hợp Swagger UI vào dự án để hỗ trợ tài liệu hóa API từ code annotations và cung cấp giao diện kiểm thử trực quan trên trình duyệt.

**Target features:**
- Install swaggo packages and tools for Go Gin routing.
- Document and annotate general, auth, and RBAC endpoints.
- Configure Swagger UI to support Bearer token authentication testing.
- Mount `/swagger/*any` route and host Swagger UI locally.

## Requirements

### Validated

- [x] Go module initialized and core dependency packages installed. (Phase 1)
- [x] Configure `docker-compose.yml` for PostgreSQL, Redis, and hot-reload Air. (Phase 1)
- [x] Setup Air config file (`.air.toml`) for live reloading of the Go application. (Phase 1)
- [x] Configure basic application structure following Go Standard Layout (`cmd/main.go`, `internal/app/...`, `internal/pkg/...`). (Phase 2)
- [x] Configure PostgreSQL database connection and database migration workflow using `golang-migrate`. (Phase 2)
- [x] Configure Redis client connection. (Phase 2)
- [x] Create a sample Ping/Pong or Healthcheck API route to verify connections and hot reloading. (Phase 3)
- [x] User registration & login with bcrypt password hashing and JWT issuance. (Phase 4)
- [x] JWT verification pipeline enforcing Client IP locking. (Phase 4)
- [x] Dynamic database schema for Roles, Resources, Actions, Policies, and User-Roles. (Phase 5)
- [x] API authorization middleware checking endpoint requests against database policies. (Phase 6)
- [x] Management APIs to dynamically create/update resources, policies, and roles. (Phase 6)

### Active

- [ ] Install Swagger Go dependencies and setup the `swag` generating CLI.
- [ ] Add Swagger general annotations in `cmd/server/main.go`.
- [ ] Annotate Ping and Health check handlers.
- [ ] Annotate Signup, Login, and Profile handlers (with JWT Security Definition).
- [ ] Annotate all dynamic RBAC management handlers (Roles, Resources, Actions, Policies, User-Roles).
- [ ] Generate Swagger docs and mount `/swagger/*any` route in router.

### Out of Scope

- Production deployment configurations (Kubernetes, AWS ECS, etc.) — Deferred to operational phases.
- Real business logic controllers (except authentication and authorization) — Out of scope for a basic bootstrap template.

## Context

The developer wants a standard, scalable starting point for a Golang backend. Having completed the base layout, connections, auth, and dynamic RBAC, the focus is now on making the API discovery and integration tests easier by introducing interactive Swagger API docs.

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
| Go Standard Layout | Scalable, clean separation of concerns for Go apps | Implemented |
| Docker-compose local env | Standardizes PostgreSQL and Redis environment | Implemented |
| golang-migrate | Clean SQL-based migrations, independent of application code | Implemented |
| slog Logging | Standard library structured logger, avoiding external deps | Implemented |
| godotenv Config | Simple local env loader, standard for Twelve-Factor apps | Implemented |
| JWT Authentication | Token-based stateless authentication suitable for REST | Implemented |
| IP Lock Security | Prevent token hijacking by verifying client IP dynamically | Implemented |
| Dynamic Database Policies | Allow changing roles, actions, and policies at runtime without redeployment | Implemented |
| Swagger API Docs | Auto-generate interactive OpenAPI documentation from source annotations | Planned |

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
*Last updated: 2026-06-22 after milestone v1.1 completion*
