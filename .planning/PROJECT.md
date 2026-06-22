# Golang Gin Backend Scaffold

## What This Is

A Golang backend service scaffold using the Gin framework, integrated with PostgreSQL and Redis. It is structured using the Go Standard Layout (`cmd/`, `internal/app/`, `internal/pkg/`) and configured with hot reloading using Air, database migrations using golang-migrate, and local service orchestration via docker-compose.

## Core Value

Provide a clean, production-ready, containerized Go backend structure that makes starting development of new endpoints and services immediate and seamless.

## Current Milestone: v1.1 Authentication & Dynamic RBAC/Policy

**Goal:** Triển khai cơ chế xác thực JWT bảo mật với tính năng khóa IP của client, kết hợp hệ thống phân quyền động dựa trên Resource, Action và Policy.

**Target features:**
- User registration & login with bcrypt password hashing and JWT issuance.
- JWT verification pipeline enforcing Client IP locking.
- Dynamic database schema for Roles, Resources, Actions, and Policies.
- API authorization middleware checking endpoint requests against database policies.
- Management APIs to dynamically create/update resources, policies, and roles.

## Requirements

### Validated

- [x] Go module initialized and core dependency packages installed. (Phase 1)
- [x] Configure `docker-compose.yml` for PostgreSQL, Redis, and hot-reload Air. (Phase 1)
- [x] Setup Air config file (`.air.toml`) for live reloading of the Go application. (Phase 1)
- [x] Configure basic application structure following Go Standard Layout (`cmd/main.go`, `internal/app/...`, `internal/pkg/...`). (Phase 2)
- [x] Configure PostgreSQL database connection and database migration workflow using `golang-migrate`. (Phase 2)
- [x] Configure Redis client connection. (Phase 2)
- [x] Create a sample Ping/Pong or Healthcheck API route to verify connections and hot reloading. (Phase 3)

### Active

- [ ] Implement user registration & login (using bcrypt for hashing and returning JWT).
- [ ] Implement JWT verification middleware enforcing Client IP locking (reject if IP changes).
- [ ] Design and implement PostgreSQL schema for Roles, Resources, Actions, and Policies.
- [ ] Implement dynamic authorization middleware verifying route resource/action against DB policies.
- [ ] Implement dynamic management APIs for Roles, Resources, Policies, and User Role assignments.

### Out of Scope

- Production deployment configurations (Kubernetes, AWS ECS, etc.) — Deferred to operational phases.
- Real business logic controllers (except authentication and authorization) — Out of scope for a basic bootstrap template.

## Context

The developer wants a standard, scalable starting point for a Golang backend. Having completed the base layout and connection setup (v1.0), the next focus is to implement security (authentication with client IP validation) and dynamic policy-based access control (RBAC/ABAC).

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
| JWT Authentication | Token-based stateless authentication suitable for REST | Planned |
| IP Lock Security | Prevent token hijacking by verifying client IP dynamically | Planned |
| Dynamic Database Policies | Allow changing roles, actions, and policies at runtime without redeployment | Planned |

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
*Last updated: 2026-06-22 after milestone v1.0 completion*
