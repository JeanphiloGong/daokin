---
id: GUIDE-DOCS-001
title: Documentation Map
type: guide
level: system
domain: shared
status: active
owner: project-owner
created_at: 2026-04-07
updated_at: 2026-04-07
last_verified_at: 2026-04-07
review_by: 2026-07-07
version: v1
source_of_truth: true
related_issues: [1, 2, 3, 4, 5]
related_docs:
  - docs/governance/documentation-governance.md
supersedes: []
superseded_by: []
tags: [docs, governance]
---

# Documentation Map

This folder is the home for durable project-level knowledge for DaoKin.

## What Lives Here

- `governance/`: project rules for documentation and other cross-cutting policy.
- `decisions/`: key cross-module decisions and future ADR/RFC records.
- `architecture/`: current system design that should match the repository as it exists today.
- `security/`: threat models, security baselines, and auditability expectations.
- `roadmap/`: milestone-level plans and delivery intent.
- `agents/`: role boundaries and coordination guidance for project contributors.
- `templates/`: copyable templates for new formal documents.

## Root Docs vs Module Docs

Use root `docs/` when a document:

- affects multiple modules
- records a project-level decision
- defines a shared contract or architecture
- sets governance, policy, or review expectations

Use module-local docs when a document is mostly owned and consumed inside one area:

- [daokin-api/docs/README.md](../daokin-api/docs/README.md)
- [daokin-web/docs/README.md](../daokin-web/docs/README.md)
- [contracts/docs/README.md](../contracts/docs/README.md)

## References vs Docs

Use `references/` for bounded working artifacts such as:

- acceptance baselines
- verification templates
- evidence captures
- iteration run logs

Do not use `references/` as the default home for long-lived architecture, policy, or contract knowledge.

## Issues vs Docs vs Code Comments

- Issues track why work exists, scope, progress, and acceptance.
- Repository docs hold durable knowledge, contracts, decisions, and operating guidance.
- Code comments explain local implementation details only.

## Adoption Rule

- New important docs should follow the documentation governance policy.
- Touched high-value docs should be upgraded when edited.
- Untouched legacy docs can remain in place until they are reviewed or updated.
