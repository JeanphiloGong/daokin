# MVP Iteration 01 Verification Template

Source baseline: `references/acceptance-criteria.md`

## 0) Run Metadata
- Iteration: 01
- Date:
- Verifier:
- Environment (local/staging):
- Commit SHA:
- API base URL:
- Contract/RPC endpoint (if used):

## 1) Evidence Rules (Auditability)
- Store raw evidence under: `references/verification/evidence/iteration-01/`
- One folder per criterion id (example: `ac-1.1/`, `ac-2.2/`).
- Every result must include:
  1. Steps executed
  2. Evidence artifact path(s)
  3. Pass/Fail decision
  4. Issue id/link when failed (`SEC-YYYYMMDD-###`)

## 2) Result Legend
- PASS: all mandatory checks under the criterion succeed.
- FAIL: any mandatory check fails or evidence is missing.
- BLOCKED: cannot execute due to missing feature/environment dependency.

## 3) Manual Flow Run Log (AC section 5)
Use this first-pass sequence once per iteration.

| Step | Endpoint or action | Observed response | Evidence path | Result (PASS/FAIL/BLOCKED) | Issue |
| --- | --- | --- | --- | --- | --- |
| 1 | Create wallet-authenticated user |  |  |  |  |
| 2 | Publish artifact |  |  |  |  |
| 3 | Join a Dao |  |  |  |  |
| 4 | Grant or deny permission |  |  |  |  |
| 5 | Execute tip/payment for authorized use |  |  |  |  |
| 6 | Export user data |  |  |  |  |
| 7 | Leave Dao |  |  |  |  |

## 4) Criterion-by-Criterion Verification Records

### AC-1.1 Artifact provenance
- Steps:
  1. Publish artifact through configured API/action.
  2. Read artifact details via API.
  3. Read corresponding contract event/log (if on-chain is enabled).
- Mandatory checks:
  - API includes `creator_wallet`, `content_uri`, `content_hash`, `created_at`.
  - Event/log includes artifact id + creator + hash/uri pointer.
- Evidence:
  - API response capture:
  - Event/log capture:
- Result: PASS / FAIL / BLOCKED
- Issue: N/A or issue id/link

### AC-1.2 Data export path
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

### AC-2.1 Permissioned use record
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

### AC-2.2 Value attribution
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

### AC-3.1 Join/leave freedom
- Steps:
  1. Join a Dao circle.
  2. Leave the same Dao.
  3. Query membership audit trail.
- Mandatory checks:
  - Join and leave endpoints/actions exist.
  - Membership state includes `joined_at`, `left_at` (reason optional).
- Evidence:
  - Join response:
  - Leave response:
  - Membership audit record:
- Result: PASS / FAIL / BLOCKED
- Issue: N/A or issue id/link

### AC-3.2 No forced grouping
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

### AC-4.1 Visible governance/moderation entry
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

### AC-4.2 Auditability baseline
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

## 5) Release Gate Checklist and Severity Criteria

### Severity classification
- Critical:
  1. Unauthorized ownership/value action is possible.
  2. Permission bypass is reproducible.
  3. Replay can cause duplicate settlement or grant.
  4. Security-relevant actions cannot be traced to actor and request.
- High:
  1. Exploitable control weakness with limited blast radius.
  2. Tampering is possible or detection is incomplete.
  3. Required control exists but has no reproducible evidence.
- Medium:
  1. Hardening gap with compensating controls.
  2. Requires privileged/internal precondition.
  3. No immediate ownership/value impact.

### Gate rules
- NO-GO when:
  1. Any open Critical issue exists.
  2. Any AC-1/AC-2 mandatory check fails.
- CONDITIONAL GO when:
  1. Critical = 0.
  2. High <= 1 with owner, mitigation, rollback, and due date <= 7 days.
  3. Medium issues are logged with target iteration.
- GO when:
  1. Critical = 0 and High = 0.
  2. All AC-1/2/3 mandatory checks pass.
  3. AC-4 has no Critical gaps.

## 6) Final Decision Record
- Decision date:
- Decision: GO / CONDITIONAL GO / NO-GO
- Critical count:
- High count:
- Medium count:
- Approved by:
- Notes and rollback links:
