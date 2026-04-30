---
id: ROADMAP-MVP-001
title: MVP Roadmap And Acceptance Criteria
type: roadmap
level: system
domain: shared
status: active
owner: project-owner
created_at: 2026-04-07
updated_at: 2026-04-30
last_verified_at: 2026-04-30
review_by: 2026-07-30
version: v2
source_of_truth: true
related_issues: [1, 2, 3, 4, 5]
related_docs:
  - docs/system-overview.md
  - docs/architecture/api-contract-v1.md
  - docs/architecture/domain-model.md
  - docs/security/threat-model-mvp.md
supersedes: []
superseded_by: []
tags: [mvp, roadmap, acceptance]
---

# MVP Roadmap And Acceptance Criteria

This is the canonical MVP plan, acceptance baseline, and iteration verification template for DaoKin. The delivery plan and pass/fail gate live in one place.

## Milestone M1: Narrative + skeleton baseline (Week 1)

### Objective
Establish a shared execution baseline and remove ambiguity.

### Deliverables
- Mission + README aligned (already present)
- Acceptance criteria baseline in this document
- Sub-agent role boundaries: [../agents/roles.md](../agents/roles.md)
- API skeleton health checks and router structure (already present)

### Exit criteria
- Team can point to this roadmap as the MVP plan and acceptance baseline.
- Domain terms are consistent across docs and code.

### Rollback plan
- If naming/model is unstable, freeze on minimal terms: User, Dao, Artifact, Permission.

## Milestone M2: Core loop implementation (Weeks 2-4)

### Objective
Run the MVP loop end to end:

```text
create artifact -> join dao -> permissioned exchange -> attribution/reputation
```

### Workstreams

#### WS-1 Auth and identity
- Wallet challenge + signature verify endpoint
- Session/token issuance and auth middleware
- Minimal audit fields in logs

#### WS-2 Artifact ownership
- Artifact create/read endpoints with provenance fields
- Storage pointer + hash persistence model
- Optional contract event linkage

#### WS-3 Dao membership
- Join/leave endpoints
- Membership history (`joined_at` / `left_at`)
- No forced assignment path

#### WS-4 Consent and exchange
- Permission grant/revoke records
- Tip/payment linkage to artifact and recipient
- Query endpoint for attribution trail

#### WS-5 Frontend flow
- Wallet auth -> create artifact -> join dao -> tip/pay -> view attribution
- Error states for denied permission and payment failures

### Exit criteria
- All mandatory checks in AC-1, AC-2, and AC-3 pass.
- AC-4 has no critical unresolved risk.

### Rollback plan
- If on-chain integration blocks progress, keep contract calls optional and preserve provenance via API records.
- If payment integration blocks progress, keep transfer records in pending/mock mode with clear flags.

## Milestone M3: Hardening pass (Week 5)

### Objective
Close top risks before broader testing.

### Deliverables
- Threat checklist for auth/permission bypass/replay/tampering
- Manual verification report against the acceptance criteria in this document
- Prioritized bug list with owner + due date

### Exit criteria
- No critical ownership/consent/attribution defect open
- Known medium risks have mitigation and owner

### Rollback plan
- Disable non-essential features; keep only core loop endpoints and UI paths.

## Acceptance Criteria

Scope is limited to the MVP core loop:

```text
create -> join -> exchange -> reputation
```

### AC-1 Ownership is explicit and traceable

#### AC-1.1 Artifact provenance
- Given a user publishes an artifact
- When the artifact is created
- Then the system stores creator wallet, content hash/URI, and created timestamp
- Then provenance can be read via API and, if on-chain, contract event/log

Pass check:
- API response includes `creator_wallet`, `content_uri`, `content_hash`, `created_at`
- Contract event includes artifact id + creator + hash/URI pointer

#### AC-1.2 Data export path
- Given a user requests export
- When export endpoint/job runs
- Then user can retrieve their artifacts and memberships in machine-readable format

Pass check:
- Export endpoint or script exists and is documented
- Sample export contains user-owned records only

### AC-2 Exchange requires consent

#### AC-2.1 Permissioned use record
- Given another user wants to use an artifact
- When permission is granted
- Then there is an explicit consent record with actor, scope, and timestamp

Pass check:
- API/DB has a permission object with `granter`, `grantee`, `scope`, `status`, `granted_at`
- No protected usage path bypasses permission check

#### AC-2.2 Value attribution
- Given value transfer happens for artifact usage (tip/payment)
- When payment succeeds
- Then transaction references artifact and recipient creator

Pass check:
- Transfer record links `tx_hash` + `artifact_id` + `recipient_wallet`
- UI/API can query all transfers for one artifact

### AC-3 Communities form organically and are reversible

#### AC-3.1 Join/leave freedom
- Given a user joins a Dao circle
- When user decides to leave
- Then user can leave without hidden lock-in

Pass check:
- Join and leave endpoints exist
- Membership state changes are auditable (`joined_at`, `left_at`, reason optional)

#### AC-3.2 No forced grouping
- Given a newly registered user
- When onboarding completes
- Then no automatic mandatory Dao assignment occurs

Pass check:
- Dao join is optional in onboarding flow
- No backend default auto-creates mandatory memberships

### AC-4 Safety and transparency minimum

#### AC-4.1 Visible governance/moderation entry
Public docs include moderation rules, reporting path, and appeal path. A minimal version is allowed for MVP.

#### AC-4.2 Auditability baseline
API request logs include request id, actor id/wallet if authenticated, action, and timestamp.

## Manual Verification Protocol

Run this sequence per iteration:

1. Create wallet-authenticated user
2. Publish artifact
3. Join a Dao
4. Grant or deny permission for use
5. Execute tip/payment for one authorized use
6. Export user data
7. Leave Dao

Record for each step:

- endpoint/action
- observed response
- pass/fail
- issue link if fail

## Iteration Verification Template

### Run Metadata
- Iteration:
- Date:
- Verifier:
- Environment:
- Commit SHA:
- API base URL:
- Contract/RPC endpoint, if used:

### Evidence Rules
- Store raw evidence in the issue, PR, release note, or manual verification report that owns the iteration.
- Every result must include steps executed, evidence path or link, pass/fail decision, and issue id/link when failed.

### Result Legend
- `PASS`: all mandatory checks under the criterion succeed.
- `FAIL`: any mandatory check fails or evidence is missing.
- `BLOCKED`: cannot execute due to missing feature or environment dependency.

### Manual Flow Run Log

| Step | Endpoint or action | Observed response | Evidence path/link | Result | Issue |
| --- | --- | --- | --- | --- | --- |
| 1 | Create wallet-authenticated user |  |  |  |  |
| 2 | Publish artifact |  |  |  |  |
| 3 | Join a Dao |  |  |  |  |
| 4 | Grant or deny permission |  |  |  |  |
| 5 | Execute tip/payment for authorized use |  |  |  |  |
| 6 | Export user data |  |  |  |  |
| 7 | Leave Dao |  |  |  |  |

### Criterion Verification Records

#### AC-1.1 Artifact provenance
- Steps:
  1. Publish artifact through configured API/action.
  2. Read artifact details via API.
  3. Read corresponding contract event/log if on-chain is enabled.
- Mandatory checks:
  - API includes `creator_wallet`, `content_uri`, `content_hash`, `created_at`.
  - Event/log includes artifact id + creator + hash/URI pointer.
- Evidence:
  - API response capture:
  - Event/log capture:
- Result: PASS / FAIL / BLOCKED
- Issue: N/A or issue id/link

#### AC-1.2 Data export path
- Steps:
  1. Trigger export endpoint/job for one test user.
  2. Download export artifact.
  3. Inspect records for ownership boundaries.
- Mandatory checks:
  - Export endpoint/script exists and is documented.
  - Sample export contains only requesting user records.
- Evidence:
  - Export trigger output:
  - Export file sample path:
- Result: PASS / FAIL / BLOCKED
- Issue: N/A or issue id/link

#### AC-2.1 Permissioned use record
- Steps:
  1. Attempt protected usage without permission.
  2. Grant permission and attempt usage again.
  3. Verify permission object fields in API/DB.
- Mandatory checks:
  - Permission object includes `granter`, `grantee`, `scope`, `status`, `granted_at`.
  - No protected usage path bypasses permission checks.
- Evidence:
  - Deny response:
  - Allow response after grant:
  - Permission record snapshot:
- Result: PASS / FAIL / BLOCKED
- Issue: N/A or issue id/link

#### AC-2.2 Value attribution
- Steps:
  1. Execute one authorized usage payment/tip.
  2. Query transfer record by artifact id.
  3. Cross-check tx reference with recipient wallet.
- Mandatory checks:
  - Transfer links `tx_hash` + `artifact_id` + `recipient_wallet`.
  - API/UI query returns transfers for one artifact.
- Evidence:
  - Payment result:
  - Transfer query output:
- Result: PASS / FAIL / BLOCKED
- Issue: N/A or issue id/link

#### AC-3.1 Join/leave freedom
- Steps:
  1. Join a Dao circle.
  2. Leave the same Dao.
  3. Query membership audit trail.
- Mandatory checks:
  - Join and leave endpoints/actions exist.
  - Membership state includes `joined_at`, `left_at`, reason optional.
- Evidence:
  - Join response:
  - Leave response:
  - Membership audit record:
- Result: PASS / FAIL / BLOCKED
- Issue: N/A or issue id/link

#### AC-3.2 No forced grouping
- Steps:
  1. Register/onboard a fresh user.
  2. Inspect membership state immediately after onboarding.
  3. Verify onboarding UI/API does not require mandatory Dao assignment.
- Mandatory checks:
  - Dao join is optional in onboarding flow.
  - No backend default mandatory membership.
- Evidence:
  - Onboarding flow capture:
  - Post-onboarding membership query:
- Result: PASS / FAIL / BLOCKED
- Issue: N/A or issue id/link

#### AC-4.1 Visible governance/moderation entry
- Steps:
  1. Locate public moderation/governance docs.
  2. Verify docs include reporting and appeal paths.
- Mandatory checks:
  - Public docs contain moderation rules.
  - Reporting path is visible.
  - Appeal path is visible.
- Evidence:
  - Doc path(s):
  - Screenshot or text excerpt path:
- Result: PASS / FAIL / BLOCKED
- Issue: N/A or issue id/link

#### AC-4.2 Auditability baseline
- Steps:
  1. Execute authenticated and unauthenticated API actions.
  2. Collect request logs for each action.
  3. Verify minimum fields in each log entry.
- Mandatory checks:
  - Logs include request id.
  - Logs include actor id/wallet when authenticated.
  - Logs include action and timestamp.
- Evidence:
  - Log sample path:
  - Request-to-log correlation sheet path:
- Result: PASS / FAIL / BLOCKED
- Issue: N/A or issue id/link

## Release Gate

### Severity classification
- Critical: unauthorized ownership/value action is possible, permission bypass is reproducible, replay can cause duplicate settlement or grant, or security-relevant actions cannot be traced to actor and request.
- High: exploitable control weakness with limited blast radius, tampering is possible or detection is incomplete, or a required control exists but has no reproducible evidence.
- Medium: hardening gap with compensating controls, privileged/internal precondition, or no immediate ownership/value impact.

### Gate rules
- `NO-GO`: any open Critical issue exists, or any AC-1/AC-2 mandatory check fails.
- `CONDITIONAL GO`: Critical = 0, High <= 1 with owner/mitigation/rollback/due date <= 7 days, and Medium issues are logged with target iteration.
- `GO`: Critical = 0, High = 0, all AC-1/2/3 mandatory checks pass, and AC-4 has no Critical gaps.

### Final Decision Record
- Decision date:
- Decision: GO / CONDITIONAL GO / NO-GO
- Critical count:
- High count:
- Medium count:
- Approved by:
- Notes and rollback links:

## Tracking Format

Use this format for each PR or milestone record:

- Intent
- Scope files
- Acceptance criteria touched (AC ids)
- Test evidence
- Rollback note

## Current Risks and Mitigations

- **Risk:** central dependency lock-in
  - **Mitigation:** keep interfaces storage-agnostic; avoid provider-specific types in domain layer.
- **Risk:** speculation-first drift
  - **Mitigation:** no token incentive features in MVP.
- **Risk:** moderation gap
  - **Mitigation:** publish minimal policy and appeal path before public testing.

## Out of Scope

- Token speculation mechanics
- Complex DAO voting systems
- Full decentralized node federation
- Regulated financial compliance automation
