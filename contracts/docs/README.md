---
id: GUIDE-CONTRACTS-DOCS-001
title: Contracts Docs Map
type: guide
level: domain
domain: contracts
status: active
owner: contracts-owner
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
tags: [contracts, docs]
---

# Contracts Docs Map

This folder is for smart-contract-local documentation owned by `contracts`.

## Put Contract Docs Here

- local contract architecture notes
- deployment and test runbooks
- contract-specific threat notes
- interface usage notes for the Solidity package
- local audit preparation notes

## Do Not Put These Here

- project-wide policy
- shared cross-module decisions
- the authoritative shared decision record for contract behavior that changes backend or frontend contracts at the same time

Those belong in root `docs/`.

## Existing Anchors

- [AGENTS.md](../AGENTS.md)
- [README.md](../README.md)
- [../contracts](../contracts)
- [../test](../test)

## Likely First Docs

- artifact registry operating notes
- tip flow assumptions
- contract test and audit checklist
