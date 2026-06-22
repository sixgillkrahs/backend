# Milestones

## v1.3 v1.3 (Shipped: 2026-06-22)

**Phases completed:** 1 phase, 5 plans, 10 tasks

**Key accomplishments:**

- Installed testify, go-sqlmock, and miniredis testing dependencies.
- Implemented comprehensive unit tests for JWT helper functions.
- Implemented unit tests for AuthMiddleware (missing header, invalid token, IP mismatch, success).
- Implemented unit tests for RBACMiddleware with GORM database mocking.
- Implemented unit/integration tests for /healthz diagnostics endpoint with Postgres and Redis ping checks.
- Implemented unit/integration tests for Signup, Login, and Profile handlers.
- Implemented unit/integration tests for roles and user-role mapping handlers.

---

## v1.2 v1.2 (Shipped: 2026-06-22)

**Phases completed:** 4 phases, 0 plans, 0 tasks

**Key accomplishments:**

- Installed swaggo packages and tools for Go Gin routing.
- Documented and annotated general, auth, and RBAC endpoints.
- Configure Swagger UI to support Bearer token authentication testing.
- Mount /swagger/*any route and host Swagger UI locally.

---

## v1.1 v1.1 (Shipped: 2026-06-22)

**Phases completed:** 3 phases, 0 plans, 0 tasks

- Implemented user registration & login (using bcrypt for hashing and returning JWT).
- Implemented JWT verification middleware enforcing Client IP locking (reject if IP changes).
- Designed and implemented PostgreSQL schema migrations and GORM models for Roles, Resources, Actions, Policies, and User-Role mappings.
- Implemented dynamic authorization middleware verifying route resource/action against DB policies.
- Implemented dynamic management APIs for Roles, Resources, Actions, Policies, and User Role assignments.
- Performed full end-to-end integration and verification testing of dynamic RBAC policies.

---

## v1.0 v1.0 (Shipped: 2026-06-22)

**Phases completed:** 3 phases, 7 plans, 22 tasks

**Key accomplishments:**

- Go module initialized and core dependency packages (Gin, Gorm, Postgres driver, Redis client) downloaded and tidied.
- Docker Compose services configured for PostgreSQL and Redis; Air configured for Go hot reloading.
- Scaffolded standard Go layout directories, configured godotenv for configuration management, and booted cmd/server/main.go with structured logging.
- Implemented PostgreSQL GORM connection pool, added golang-migrate auto-migration support on boot, and created the initial SQL migration files.
- Configured Redis client connection pool using go-redis/v9 and verified the connection on boot.
- Initialized Gin HTTP engine, registered structured request logging middleware using slog, and exposed the `/ping` route.
- Implemented the dynamic `/healthz` diagnostics endpoint testing PostgreSQL and Redis active connections.

---
