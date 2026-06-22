# Roadmap: Golang Gin Backend Scaffold

## Overview

Initialize a Golang backend skeleton following the Go Standard Layout with Gin, PostgreSQL (managed via golang-migrate), Redis, and Air for hot reloading.

## Phases

- [x] **Phase 1: Environment & Tooling** - Set up Go module, dependencies, Docker Compose, and Air. (completed 2026-06-22)
- [x] **Phase 2: Directory Layout & Connections** - Scaffold project layout and configure Postgres/Redis connections with migrations. (completed 2026-06-22)
- [x] **Phase 3: Router & Healthcheck APIs** - Configure Gin router and verify status with `/ping` and `/healthz` API endpoints. (completed 2026-06-22)

## Phase Details

### Phase 1: Environment & Tooling

**Goal**: Set up Go module, dependencies, Docker Compose, and Air.
**Depends on**: Nothing
**Requirements**: ENV-01, ENV-02, ENV-03, ENV-04, ENV-05, ENV-06
**Success Criteria** (what must be TRUE):

  1. Go module initialized and core dependency packages installed.
  2. docker-compose.yml runs PostgreSQL and Redis without error.
  3. Air config is set up and ready to watch changes.

**Plans**: 2 plans

Plans:

- [x] 01-01: Go module initialization & dependency setup
- [x] 01-02: Docker Compose & Air configuration

### Phase 2: Directory Layout & Connections

**Goal**: Scaffold project layout and configure Postgres/Redis connections with migrations.
**Depends on**: Phase 1
**Requirements**: DIR-01, DIR-02, DIR-03, DB-01, DB-02, DB-03, CACHE-01
**Success Criteria** (what must be TRUE):

  1. Directories cmd/, internal/app/, and internal/pkg/ structured.
  2. Database migration logic and initial SQL migration files created.
  3. Postgres and Redis connection pools established upon app boot.

**Plans**: 3 plans

Plans:
**Wave 1**

- [x] 02-01: Go Standard Layout scaffolding

**Wave 2** *(blocked on Wave 1 completion)*

- [x] 02-02: PostgreSQL connection & golang-migrate setup
- [x] 02-03: Redis connection setup

### Phase 3: Router & Healthcheck APIs

**Goal**: Configure Gin router and verify status with `/ping` and `/healthz` API endpoints.
**Depends on**: Phase 2
**Requirements**: API-01, API-02
**Success Criteria** (what must be TRUE):

  1. Gin router parses requests and maps endpoints.
  2. `/ping` endpoint returns valid JSON.
  3. `/healthz` dynamically pings database and Redis connection and returns healthy status.

**Plans**: 2 plans

Plans:
**Wave 1**

- [x] 03-01: Gin router & ping route setup

**Wave 2** *(blocked on Wave 1 completion)*

- [x] 03-02: Healthcheck route with DB & Redis ping verification

## Progress

**Execution Order:**
Phases execute in numeric order: 1 → 2 → 3

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Environment & Tooling | 2/2 | Complete    | 2026-06-22 |
| 2. Directory Layout & Connections | 3/3 | Complete    | 2026-06-22 |
| 3. Router & Healthcheck APIs | 2/2 | Complete    | 2026-06-22 |
