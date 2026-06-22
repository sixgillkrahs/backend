<!-- GSD:project-start source:PROJECT.md -->

## Project

**Golang Gin Backend Scaffold**

A Golang backend service scaffold using the Gin framework, integrated with PostgreSQL and Redis. It is structured using the Go Standard Layout (`cmd/`, `internal/app/`, `internal/pkg/`) and configured with hot reloading using Air, database migrations using golang-migrate, and local service orchestration via docker-compose.

**Core Value:** Provide a clean, production-ready, containerized Go backend structure that makes starting development of new endpoints and services immediate and seamless.

### Constraints

- **Language**: Go (v1.22+ recommended) — Core language requirement.
- **Framework**: Gin-Gonic — Core HTTP web framework.
- **Data Stores**: PostgreSQL & Redis — Required state and caching systems.
- **Hot Reload**: Air — Tool for hot reloading.
- **Orchestration**: Docker Compose — Local development orchestration.
- **Migrations**: golang-migrate — Core migration utility.

<!-- GSD:project-end -->

<!-- GSD:stack-start source:codebase/STACK.md -->

## Technology Stack

## Languages

- None (Greenfield project - no application code yet)
- JavaScript (Node.js) - Scripting and GSD configuration hooks

## Runtime

- Node.js v22.22.1
- npm (configured via global/local npx scripts)
- Lockfile: None (greenfield project)

## Frameworks

- None (Greenfield project)
- None (Greenfield project)
- GSD Core v1.5.0 - Meta-prompting and context engineering framework

## Key Dependencies

- @opengsd/gsd-core v1.5.0 - Local agent workflows and GSD core tools
- Git - Version control

## Configuration

- Local GSD environment in `.agents` directory
- None (Greenfield project)

## Platform Requirements

- Windows (Current OS environment)
- Node.js v22.22.1+
- TBD (To be defined as requirements emerge)

<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->

## Conventions

## Naming Patterns

- **Files:** TBD (Recommendation: kebab-case for modules/scripts, PascalCase for classes/components)
- **Functions:** TBD (Recommendation: camelCase)
- **Variables:** TBD (Recommendation: camelCase, UPPER_SNAKE_CASE for constants)
- **Types:** TBD (Recommendation: PascalCase for types and interfaces, no I-prefix for interfaces)

## Code Style

- **Formatting:** TBD (Prettier / ESLint settings will be configured as development begins)
- **Indentation:** TBD (Recommendation: 2 spaces)

## Import Organization

- **Order:** TBD (Recommendation: External libraries first, then internal alias modules, then relative paths)

## Error Handling

- **Patterns:** TBD (Recommendation: Use structured error throwing/catching, use standard Node.js custom error classes)

## Logging

- **Framework:** TBD (Recommendation: Structured logger like Pino or Winston, fallback to console.log/console.error)

## Comments

- **When to Comment:** TBD (Recommendation: Explain "why", not "what")
- **TODO Comments:** TBD (Recommendation: Use `// TODO: description` or `// TODO(issue-id): description`)

## Function Design

- **Size:** TBD (Recommendation: Keep functions small, <50 lines, single responsibility)
- **Parameters:** TBD (Recommendation: Prefer options objects for >3 parameters)

## Module Design

- **Exports:** TBD (Recommendation: Named exports preferred for modules, index.ts for folder exports)

<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->

## Architecture

## Pattern Overview

- Greenfield state
- GSD-driven workflow configured locally
- No application servers, microservices, or client-side assets present yet

## Layers

- **TBD (Application Layers):** The application layers will be defined during the initial design and requirements definition phase.
- **GSD Meta-Planning Layer:** Local GSD settings and workflows are managed inside the `.agents/` folder.

## Data Flow

- **TBD (Application Data Flow):** Standard requests, API flows, and event processing pipelines will be documented as soon as backend architecture is defined.

## Key Abstractions

- **TBD (Application Abstractions):** Domain services, models, and controllers are yet to be created.

## Entry Points

- **TBD (Application Entry Points):** The main entry points (e.g., HTTP server listener, worker entry points) will be established as part of the initial implementation.

## Error Handling

- **Strategy:** TBD

## Cross-Cutting Concerns

- **Logging:** TBD
- **Validation:** TBD
- **Authentication:** TBD

<!-- GSD:architecture-end -->

<!-- GSD:skills-start source:skills/ -->

## Project Skills

No project skills found. Add skills to any of: `.agents/skills/`, `.agents/skills/`, `.cursor/skills/`, `.github/skills/`, or `.codex/skills/` with a `SKILL.md` index file.
<!-- GSD:skills-end -->

<!-- GSD:workflow-start source:GSD defaults -->

## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:

- `/gsd-quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd-debug` for investigation and bug fixing
- `/gsd-execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->

<!-- GSD:profile-start -->

## Developer Profile

> Profile not yet configured. Run `/gsd-profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->
