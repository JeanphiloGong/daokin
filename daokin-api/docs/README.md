---
id: GUIDE-BACKEND-DOCS-001
title: Backend Docs Map
type: guide
level: domain
domain: backend
status: active
owner: backend-owner
created_at: 2026-04-07
updated_at: 2026-04-07
last_verified_at: 2026-04-07
review_by: 2026-07-07
version: v1
source_of_truth: true
related_issues: [2, 4]
related_docs:
  - docs/README.md
  - docs/governance/documentation-governance.md
supersedes: []
superseded_by: []
tags: [backend, docs]
---

# Backend Docs Map

This folder is for backend-local documentation owned by `daokin-api`.

## Put Backend Docs Here

- local architecture notes for the Go service
- backend implementation guides
- local runbooks
- module-specific data flow notes
- local troubleshooting notes that do not define project-wide policy

## Do Not Put These Here

- cross-module decisions
- project-wide governance or policy
- the authoritative shared API or contract source of truth if it affects multiple modules

Those belong in root `docs/`.

## Existing Anchors

- [AGENTS.md](../AGENTS.md)
- [README.md](../README.md)
- [../internal](../internal)

## Likely First Docs

- auth flow notes
- permission and transfer flow notes
- persistence migration notes
- API runbooks for local verification
