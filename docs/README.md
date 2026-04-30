---
id: GUIDE-DOCS-001
title: Documentation Map
type: guide
level: system
domain: shared
status: active
owner: project-owner
created_at: 2026-04-07
updated_at: 2026-04-30
last_verified_at: 2026-04-30
review_by: 2026-07-07
version: v1
source_of_truth: true
related_issues: [1, 2, 3, 4, 5]
related_docs:
  - docs/governance/documentation-governance.md
  - docs/system-overview.md
supersedes: []
superseded_by: []
tags: [docs, governance]
---

# Documentation Map

This is the documentation landing page for DaoKin. Start here when you need the reading path, the source of truth for a topic, or the boundary between project docs and module-local docs.

## Start Here

- Repository summary and quick start: [../README.md](../README.md)
- Long-term mission and manifesto: [../MISSION.md](../MISSION.md)
- System overview and module boundaries: [system-overview.md](system-overview.md)
- Documentation governance: [governance/documentation-governance.md](governance/documentation-governance.md)
- MVP acceptance baseline: [../references/acceptance-criteria.md](../references/acceptance-criteria.md)

## Reading Paths

- New project reader: [../README.md](../README.md) -> [../MISSION.md](../MISSION.md) -> [system-overview.md](system-overview.md)
- MVP implementer: [system-overview.md](system-overview.md) -> [roadmap/mvp.md](roadmap/mvp.md) -> [architecture/api-contract-v1.md](architecture/api-contract-v1.md)
- Backend contributor: [../daokin-api/README.md](../daokin-api/README.md) -> [architecture/api-contract-v1.md](architecture/api-contract-v1.md) -> [architecture/domain-model.md](architecture/domain-model.md)
- Frontend contributor: [../daokin-web/README.md](../daokin-web/README.md) -> [architecture/api-contract-v1.md](architecture/api-contract-v1.md)
- Contract contributor: [../contracts/README.md](../contracts/README.md) -> [architecture/contract-event-map.md](architecture/contract-event-map.md)
- Security or acceptance reviewer: [security/threat-model-mvp.md](security/threat-model-mvp.md) -> [../references/verification/mvp-iteration-01.md](../references/verification/mvp-iteration-01.md)

## Authority Map

- `governance/`: project rules for documentation and other cross-cutting policy.
- `decisions/`: key cross-module decisions and future ADR/RFC records.
- `architecture/`: current system design that should match the repository as it exists today.
- `security/`: threat models, security baselines, and auditability expectations.
- `roadmap/`: milestone-level plans and delivery intent.
- `agents/`: role boundaries and coordination guidance for project contributors.
- `templates/`: copyable templates for new formal documents.

## Root Docs vs Module Entrypoints

Use root `docs/` when a document:

- affects multiple modules
- records a project-level decision
- defines a shared contract or architecture
- sets governance, policy, or review expectations

Use module root README files when a reader needs the purpose, boundary, commands, or local orientation for one owned module:

- [../daokin-api/README.md](../daokin-api/README.md)
- [../daokin-web/README.md](../daokin-web/README.md)
- [../contracts/README.md](../contracts/README.md)

Create module-local docs only when a module has durable local knowledge that is too detailed for its README, such as a local runbook, local current-state note, or implementation guide.

## References vs Docs

Use `references/` for bounded working artifacts such as:

- acceptance baselines
- verification templates
- evidence captures
- iteration run logs

Do not use `references/` as the default home for long-lived architecture, policy, or contract knowledge.

## Placement Rules

- Put system purpose, shared architecture, shared contracts, governance, and roadmap material in root `docs/`.
- Put module purpose, commands, boundaries, and first-hop navigation in each module `README.md`.
- Put module-local details under `<module>/docs/` only after that module has real durable docs to index.
- Keep parent docs as summaries and routing pages; keep detailed implementation and verification notes at the lowest owning node.
- Keep test coverage, fixtures, harness notes, and verification evidence near the owning test or verification artifact when the test asset is the subject.

## Issues vs Docs vs Code Comments

- Issues track why work exists, scope, progress, and acceptance.
- Repository docs hold durable knowledge, contracts, decisions, and operating guidance.
- Code comments explain local implementation details only.

## Adoption Rule

- New important docs should follow the documentation governance policy.
- Touched high-value docs should be upgraded when edited.
- Untouched legacy docs can remain in place until they are reviewed or updated.
