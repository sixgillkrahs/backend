# Roadmap: Golang Gin Backend Scaffold

## Overview

Initialize a Golang backend skeleton following the Go Standard Layout with Gin, PostgreSQL (managed via golang-migrate), Redis, and Air for hot reloading. Subsequent milestones implement authentication with IP locking, dynamic RBAC, and interactive Swagger API documentation.

## Phases

- [x] **Phase 1: Environment & Tooling** - Set up Go module, dependencies, Docker Compose, and Air. (completed 2026-06-22)
- [x] **Phase 2: Directory Layout & Connections** - Scaffold project layout and configure Postgres/Redis connections with migrations. (completed 2026-06-22)
- [x] **Phase 3: Router & Healthcheck APIs** - Configure Gin router and verify status with `/ping` and `/healthz` API endpoints. (completed 2026-06-22)
- [x] **Phase 4: User Authentication & IP-Locked JWTs** - Set up signup, login, password hashing, and token verification middleware with client IP verification. (completed 2026-06-22)
- [x] **Phase 5: Dynamic RBAC Database Schema & Models** - Create SQL migrations for dynamic roles, resources, actions, policies, and user roles, and define GORM models. (completed 2026-06-22)
- [x] **Phase 6: Dynamic Authorization Middleware & RBAC Management APIs** - Implement authorization middleware and admin endpoints to manage roles, resources, policies, and role assignment. (completed 2026-06-22)
- [x] **Phase 7: Swagger Dependencies & Handler Annotations** - Install swaggo packages and add annotations to all existing API endpoints. (completed 2026-06-22)
- [x] **Phase 8: Swagger Routing & Interactive Verification** - Generate Swagger documentation files, mount Swagger UI route, and perform integration testing. (completed 2026-06-22)
- [x] **Phase 9: Auth & RBAC Testing (Unit & Integration)** - Write unit and integration testing suites covering tokens, middleware, and handlers. (completed 2026-06-22)

---

## Phase Details

### Phase 1: Environment & Tooling (Complete)
- **Goal**: Set up Go module, dependencies, Docker Compose, and Air.
- **Depends on**: Nothing
- **Requirements**: ENV-01, ENV-02, ENV-03, ENV-04, ENV-05, ENV-06
- **Plans**:
  - [x] 01-01: Go module initialization & dependency setup
  - [x] 01-02: Docker Compose & Air configuration

### Phase 2: Directory Layout & Connections (Complete)
- **Goal**: Scaffold project layout and configure Postgres/Redis connections with migrations.
- **Depends on**: Phase 1
- **Requirements**: DIR-01, DIR-02, DIR-03, DB-01, DB-02, DB-03, CACHE-01
- **Plans**:
  - [x] 02-01: Go Standard Layout scaffolding
  - [x] 02-02: PostgreSQL connection & golang-migrate setup
  - [x] 02-03: Redis connection setup

### Phase 3: Router & Healthcheck APIs (Complete)
- **Goal**: Configure Gin router and verify status with `/ping` and `/healthz` API endpoints.
- **Depends on**: Phase 2
- **Requirements**: API-01, API-02
- **Plans**:
  - [x] 03-01: Gin router & ping route setup
  - [x] 03-02: Healthcheck route with DB & Redis ping verification

### Phase 4: User Authentication & IP-Locked JWTs (Complete)
- **Goal**: Implement signup, login, password hashing, and token verification middleware with client IP verification.
- **Depends on**: Phase 3
- **Requirements**: AUTH-01, AUTH-02, AUTH-03, AUTH-04, AUTH-05
- **Plans**:
  - [x] 04-01: Auth endpoints for Signup and Login
  - [x] 04-02: JWT token generation with client IP claims
  - [x] 04-03: JWT Authentication middleware with IP locking

### Phase 5: Dynamic RBAC Database Schema & Models (Complete)
- **Goal**: Create SQL migrations for dynamic roles, resources, actions, policies, and user roles, and define GORM models.
- **Depends on**: Phase 4
- **Requirements**: AUTH-06, AUTH-07
- **Plans**:
  - [x] 05-01: Database migrations for RBAC schema and seed data
  - [x] 05-02: Define GORM structs/models for RBAC entities

### Phase 6: Dynamic Authorization Middleware & RBAC Management APIs (Complete)
- **Goal**: Implement authorization middleware and admin endpoints to manage roles, resources, policies, and role assignment.
- **Depends on**: Phase 5
- **Requirements**: AUTH-08, AUTH-09, AUTH-10, AUTH-11, AUTH-12
- **Plans**:
  - [x] 06-01: Dynamic authorization middleware logic
  - [x] 06-02: RBAC management controllers and router endpoints
  - [x] 06-03: User-role mapping management endpoints

### Phase 7: Swagger Dependencies & Handler Annotations (Complete)
- **Goal**: Install swaggo packages and add annotations to all existing API endpoints.
- **Depends on**: Phase 6
- **Requirements**: SWAG-01, SWAG-02, SWAG-03, SWAG-04, SWAG-05
- **Success Criteria**:
  1. Go package dependencies `github.com/swaggo/gin-swagger` and `github.com/swaggo/files` added to module.
  2. Swagger general config annotations added to `cmd/server/main.go`.
  3. API annotation comments added to all controller handlers (Health, Ping, Signup, Login, Profile, Roles, Resources, Actions, Policies, User-Roles).
- **Plans**:
  - [x] 07-01: Install swaggo dependencies and configure CLI tool
  - [x] 07-02: Add Swagger annotations to general and authentication endpoints
  - [x] 07-03: Add Swagger annotations to RBAC management endpoints

### Phase 8: Swagger Routing & Interactive Verification (Complete)
- **Goal**: Generate Swagger documentation files, mount Swagger UI route, and perform integration testing.
- **Depends on**: Phase 7
- **Requirements**: SWAG-06, SWAG-07
- **Success Criteria**:
  1. Swagger spec files generated under `docs/` using `swag init`.
  2. Swagger routing endpoint `/swagger/*any` mounted and accessible.
  3. Interactive UI successfully loaded, permitting request execution and Bearer authentication token setup.
- **Plans**:
  - [x] 08-01: Mount Swagger UI route and integrate generated docs
  - [x] 08-02: Perform interactive endpoint validation tests via Swagger UI

### Phase 9: Auth & RBAC Testing (Unit & Integration) (Complete)
- **Goal**: Implement unit and integration testing suites using mock connections.
- **Depends on**: Phase 8
- **Requirements**: TEST-01, TEST-02, TEST-03
- **Success Criteria**:
  1. Unit tests written for Auth tokens, middleware, and handlers.
  2. Integration tests written for RBAC middleware and handlers using sqlmock/miniredis.
  3. Full test suite executes successfully with >80% coverage on core packages.
- **Plans**:
  - [x] 09-01: Setup testify, go-sqlmock, and miniredis testing environments
  - [x] 09-02: Implement token helper and middleware unit tests
  - [x] 09-03: Implement health, authentication, and RBAC handler integration tests

---

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3 → 4 → 5 → 6 → 7 → 8 → 9

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Environment & Tooling | 2/2 | Complete | 2026-06-22 |
| 2. Directory Layout & Connections | 3/3 | Complete | 2026-06-22 |
| 3. Router & Healthcheck APIs | 2/2 | Complete | 2026-06-22 |
| 4. User Authentication & IP-Locked JWTs | 3/3 | Complete | 2026-06-22 |
| 5. Dynamic RBAC Database Schema & Models | 2/2 | Complete | 2026-06-22 |
| 6. Dynamic Authorization Middleware & RBAC Management APIs | 3/3 | Complete | 2026-06-22 |
| 7. Swagger Dependencies & Handler Annotations | 3/3 | Complete | 2026-06-22 |
| 8. Swagger Routing & Interactive Verification | 2/2 | Complete | 2026-06-22 |
| 9. Auth & RBAC Testing (Unit & Integration) | 3/3 | Complete | 2026-06-22 |
