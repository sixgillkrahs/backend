# Codebase Structure

**Analysis Date:** 2026-06-22

## Directory Layout

```
backend/
├── .agents/           # Local GSD Core configuration & skills
├── .git/              # Git version control metadata
└── .planning/         # GSD Project planning artifacts & codebase maps
```

## Directory Purposes

**.agents/**
- Purpose: Stores local GSD Core framework configurations, hooks, scripts, and skills.
- Contains: Configuration JSON files, JavaScript hooks, and markdown skills.
- Key files: [settings.json](file:///d:/works/business-chat/backend/.agents/settings.json) - GSD hook configuration.
- Subdirectories: `skills/` (agent capabilities), `hooks/` (gated actions), `agents/` (agent prompt definitions).

**.git/**
- Purpose: Git source control directory.
- Contains: Commit history, refs, object store, hooks.

**.planning/**
- Purpose: Stores GSD spec-driven development documents.
- Contains: Markdown planning documents.
- Subdirectories: `codebase/` (stores files capturing codebase intelligence).

## Key File Locations

**Configuration:**
- `.agents/settings.json` - GSD hook settings.
- `.planning/codebase/STACK.md` - Tech stack documentation.
- `.planning/codebase/ARCHITECTURE.md` - Architecture overview.
- `.planning/codebase/STRUCTURE.md` - Directory layout and naming conventions.

## Naming Conventions

**Files:**
- kebab-case.md: Markdown documents.
- UPPERCASE.md: General project documentation (like PROJECT.md, STACK.md).

## Where to Add New Code

**Application Source Code:**
- TBD (Will be defined as backend codebase structure is established).

---

*Structure analysis: 2026-06-22*
*Update when directory structure changes*
