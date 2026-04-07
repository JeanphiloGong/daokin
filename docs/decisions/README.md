---
id: GUIDE-DECISIONS-001
title: Key Decisions Folder
type: guide
level: system
domain: shared
status: active
owner: project-owner
created_at: 2026-04-07
updated_at: 2026-04-07
last_verified_at: 2026-04-07
review_by: 2026-08-07
version: v1
source_of_truth: true
related_issues: [1]
related_docs:
  - docs/README.md
  - docs/governance/documentation-governance.md
supersedes: []
superseded_by: []
tags: [decisions, adr, rfc]
---

# Key Decisions Folder

This folder is the main home for durable project decisions that affect more than one module.

## Put A Document Here When

- the decision changes project direction
- the decision affects multiple modules
- the decision changes a shared contract or operating rule
- the decision should still be understandable months later without reopening issue threads

## Do Not Put It Here When

- the note is only for one module's local implementation
- the content is only a task checklist
- the content is just progress tracking for an issue

## File Conventions

Use one of these prefixes:

- `adr-YYYYMMDD-topic.md` for accepted decisions
- `rfc-YYYYMMDD-topic.md` for proposals under review

Examples:

- `adr-20260407-doc-governance.md`
- `rfc-20260407-api-auth-shape.md`

## Relationship To Module Docs

If a note starts in a module-local docs folder but becomes a cross-module decision, promote it here and leave a pointer in the original module docs if needed.
