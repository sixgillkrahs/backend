# Coding Conventions

**Analysis Date:** 2026-06-22

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

---

*Convention analysis: 2026-06-22*
*Update when patterns change*
