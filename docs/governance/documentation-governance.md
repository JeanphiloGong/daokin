---
id: POLICY-DOCS-001
title: Documentation Governance
type: policy
level: system
domain: shared
status: active
owner: project-owner
created_at: 2026-04-07
updated_at: 2026-04-07
last_verified_at: 2026-04-07
review_by: 2026-10-07
version: v1
source_of_truth: true
related_issues: [1, 2, 3, 4, 5]
related_docs:
  - docs/README.md
  - docs/decisions/README.md
  - docs/templates/formal-doc-template.md
supersedes: []
superseded_by: []
tags: [policy, docs, governance]
---

# Documentation Governance

## Goal

This policy defines how DaoKin stores durable knowledge so architecture, roadmap, security, and operating documents do not drift across issues, docs, and code.

## Governance Scope

This policy applies to:

- formal docs under `docs/`
- future module-local docs under `daokin-api/docs/`, `daokin-web/docs/`, and `contracts/docs/`
- high-value documents that are edited after this policy is adopted

This policy does not force an immediate rewrite of all legacy markdown files.

## Document Types

DaoKin uses these formal document types:

- `policy`: repository rules and governance requirements
- `rfc`: proposed cross-module or domain changes
- `adr`: durable decision records
- `architecture`: current real design
- `spec`: stable contracts such as API or event behavior
- `guide`: human-facing usage or contribution guidance
- `runbook`: repeatable operational procedures
- `postmortem`: durable incident learning
- `roadmap`: milestone intent and planning state

Each formal doc should have one primary job. If a document is trying to propose, decide, and describe current reality at the same time, split it.

## Scope Levels

DaoKin uses four scope levels:

- `system`: the whole repository or product
- `domain`: a major area such as backend, frontend, contracts, security, or shared
- `module`: one service, flow, or major feature slice
- `component`: a narrower sub-area inside one module

Do not use a deeper level model unless the repo proves it needs one.

## Metadata Standard

New formal docs should use YAML front matter.

Required fields:

- `id`
- `title`
- `type`
- `level`
- `domain`
- `status`
- `owner`
- `created_at`
- `updated_at`

Strongly recommended fields:

- `last_verified_at`
- `review_by`
- `version`
- `source_of_truth`
- `related_issues`
- `related_docs`
- `supersedes`
- `superseded_by`

## Lifecycle States

DaoKin uses these lifecycle states:

- `draft`
- `review`
- `accepted`
- `implemented`
- `active`
- `deprecated`
- `superseded`
- `archived`

Expected flows:

- `rfc`: `draft -> review -> accepted -> implemented`
- `adr`: `accepted -> active`
- `architecture`, `spec`, `guide`, `runbook`, `policy`, `roadmap`: normally `active`

When a document is replaced:

1. do not delete it by default
2. mark it `superseded`
3. link the replacement in `superseded_by`
4. link the old document in the new document's `supersedes`

## Source of Truth Rules

- Issues are authoritative for work tracking, delivery scope, and acceptance status.
- Repository docs are authoritative for durable decisions, contracts, architecture, governance, and operations.
- Code comments are authoritative only for local implementation details.

Only one document in the same scope should be marked `source_of_truth: true` for the same contract or design surface.

Current controlled baseline:

- `docs/roadmap/mvp.md` is the MVP roadmap, acceptance baseline, and verification template.

## Storage Rules

### Root `docs/`

Use root `docs/` for:

- project-level and cross-module policies
- key decisions
- shared architecture
- shared specs and contracts
- project-level security and roadmap material

### Module-Local `*/docs/`

Use module-local docs for:

- implementation notes
- local architecture details
- local runbooks
- developer guides scoped to one module

Module-local docs must not become a parallel source of truth for project-wide decisions or shared contracts.

### Verification Evidence

Keep durable acceptance criteria and verification templates in their owning docs. For MVP, that owner is `docs/roadmap/mvp.md`.

Record raw evidence in the issue, PR, release note, or verification report that owns the iteration.

### Promotion Rule

If a module-local document becomes shared across multiple modules or starts acting as the authoritative project decision or contract, promote it into root `docs/`.

## Naming Rules

- Use lowercase kebab case for filenames.
- Use topic-based names, not generic names such as `notes.md` or `misc.md`.
- Keep stable, long-lived docs on stable filenames.
- Prefer ASCII names for new docs.

## Ownership and Review

Default owners:

- project-wide docs: `project-owner`
- backend docs: `backend-owner`
- frontend docs: `frontend-owner`
- contracts docs: `contracts-owner`

Recommended review cadence:

- fast-moving API and MVP docs: every 1 to 3 months
- normal architecture and spec docs: every 3 to 6 months
- slow-moving policy and governance docs: every 6 to 12 months

Update a formal doc when:

- behavior changes
- a contract changes
- a key decision is made or reversed
- verification proves the document is stale

## Adoption and Enforcement

This policy is required for:

- new formal docs
- touched high-value docs in `docs/`
- new module-local docs created after adoption

This policy is advisory for untouched legacy docs until they are edited or reviewed.

## Minimal Automation Baseline

Automation is not required in the first pass, but the target baseline is:

- validate required front matter fields
- validate `type`, `level`, `domain`, and `status`
- check broken internal links
- flag expired `review_by`
- detect duplicate `source_of_truth: true` in the same scope
