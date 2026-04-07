---
id: GUIDE-TEMPLATE-001
title: Formal Document Template
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
source_of_truth: false
related_issues: []
related_docs:
  - docs/governance/documentation-governance.md
supersedes: []
superseded_by: []
tags: [template, docs]
---

# Formal Document Template

Use this template for new formal docs in `docs/` and `*/docs/`.

```yaml
---
id: SPEC-YYYY-NNN
title: Example Title
type: spec
level: module
domain: backend
status: draft
owner: backend-owner
created_at: 2026-04-07
updated_at: 2026-04-07
last_verified_at: 2026-04-07
review_by: 2026-07-07
version: v1
source_of_truth: false
related_issues: []
related_docs: []
supersedes: []
superseded_by: []
tags: [example]
---
```

## Suggested Sections

Choose the sections that match the document type.

- Background / Problem
- Goal / Scope
- Current Design or Proposed Change
- External Contract
- Risks / Dependencies
- Verification
- Replacement / Migration Notes

## Notes

- Keep one document type per file.
- If the file is an RFC, use proposal language and lifecycle states.
- If the file is an architecture or spec doc, describe current reality rather than intent.
