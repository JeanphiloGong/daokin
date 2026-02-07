# MVP Threat Model (Iteration 01)

## 0) Document Control
- Owner role: Security & Acceptance
- Last updated: 2026-02-07
- Scope: MVP loop `create -> join -> exchange -> reputation`
- In-scope families: auth, permission bypass, replay, tampering, logging/audit
- Out of scope: token speculation, complex DAO voting, full federation

## 1) Security Objectives
1. Authentication cannot be forged, replayed, or bypassed.
2. Protected usage requires explicit and valid permission.
3. Permission and payment actions are replay resistant.
4. Ownership, provenance, and transfer records are tamper evident.
5. Sensitive actions are auditable with actor and request traceability.

## 2) Trust Boundaries and Assets
### Boundaries
1. Client wallet <-> API
2. API <-> persistence (DB/object storage)
3. API <-> contract/RPC/event source
4. API <-> logging/audit sink

### Security-Critical Assets
- Auth state: wallet address, nonce, signature, session token
- Permission record: `granter`, `grantee`, `scope`, `status`, `granted_at`
- Provenance record: `creator_wallet`, `content_uri`, `content_hash`, `created_at`
- Value attribution: `tx_hash`, `artifact_id`, `recipient_wallet`
- Audit tuple: `request_id`, actor (`id` or wallet), action, timestamp, outcome

## 3) Severity Rubric (Release Gate Input)
- Critical: unauthorized ownership/value impact OR no security auditability for sensitive actions.
- High: exploitable control weakness with bounded blast radius OR controls cannot be evidenced.
- Medium: defense-in-depth weakness requiring extra preconditions or privileged context.

## 4) Threat Catalog and Required Controls

### TM-AUTH-01: Wallet Auth Replay or Forgery
- Threat: attacker reuses a signed payload or forges login context.
- Impact: unauthorized account/session access.
- Required controls:
  1. Nonce is single-use and expires in <= 5 minutes.
  2. Signature binds `domain`, `chain_id`, `nonce`, and wallet address.
  3. Session/token TTL is finite and revocable.
  4. Auth failures are fail-closed (no partial session).
- Verification steps:
  1. Request nonce and sign once; login succeeds.
  2. Replay same signature/nonce; API must reject.
  3. Tamper signed fields (`domain` or `chain_id`); API must reject.
- Evidence to collect:
  - Request/response captures for valid and replay attempts.
  - Audit logs containing request id, actor wallet, action, and deny reason.
- Severity if present: Critical.

### TM-PERM-01: Permission Bypass on Protected Usage
- Threat: artifact usage endpoint can be called without valid permission.
- Impact: unauthorized use and value leakage.
- Required controls:
  1. Server-side authorization check on every protected endpoint.
  2. Deny-by-default when permission is missing, expired, revoked, or scope mismatch.
  3. Decision references exact tuple: actor + artifact_id + scope + status.
  4. Revocation takes effect immediately.
- Verification steps:
  1. Call protected usage path without permission; must deny.
  2. Grant permission and retry; must allow.
  3. Revoke permission and retry; must deny.
- Evidence to collect:
  - API responses for deny/allow/deny sequence.
  - Permission record snapshots before and after revocation.
  - Audit logs for authorization decisions.
- Severity if present: Critical.

### TM-REPLAY-01: Duplicate Permission/Payment Execution
- Threat: same grant/payment intent is processed multiple times.
- Impact: duplicate grants, double transfer, inconsistent state.
- Required controls:
  1. Idempotency key for grant/payment write paths.
  2. Replay cache or unique constraint on signed intent hash / tx hash.
  3. Expiration window on signed intents.
  4. Deterministic response for duplicates (no second state mutation).
- Verification steps:
  1. Submit same signed grant twice; second call must be no-op/rejected.
  2. Submit duplicate payment callback/confirmation; no double settlement.
- Evidence to collect:
  - API traces with same idempotency key.
  - DB/event records proving one effective mutation.
  - Transfer records showing unique `tx_hash` linkage.
- Severity if present: Critical (if value moved), otherwise High.

### TM-TAMP-01: Provenance or Transfer Tampering
- Threat: metadata or transfer linkage is modified without trace.
- Impact: ownership dispute, broken attribution, unverifiable history.
- Required controls:
  1. `content_hash` immutable after publish (new version requires new record/event).
  2. Provenance and transfer changes are append-only or fully versioned.
  3. Contract event/log pointer and API record remain consistent.
  4. Export includes checksums for integrity verification.
- Verification steps:
  1. Attempt post-create mutation of immutable fields; must fail.
  2. Compare API provenance with contract/event data for one artifact.
  3. Verify export checksums after download.
- Evidence to collect:
  - Rejected mutation response.
  - API and event comparison output.
  - Export checksum verification output.
- Severity if present: High (Critical if unauthorized ownership transfer is possible).

### TM-AUD-01: Incomplete or Non-Correlatable Audit Trail
- Threat: logs miss identity/context and cannot reconstruct security decisions.
- Impact: undetected abuse and no forensic accountability.
- Required controls:
  1. Mandatory fields in security-relevant logs:
     - `request_id`, actor id/wallet (if authenticated), action, timestamp, outcome.
  2. Both allow and deny paths are logged.
  3. Log retention and immutability policy is documented.
  4. Clock source is consistent (UTC) for event ordering.
- Verification steps:
  1. Execute auth, permission deny, permission allow, payment attempt.
  2. For each request, find one matching audit record by `request_id`.
  3. Confirm actor/action/outcome fields exist and are readable.
- Evidence to collect:
  - Log excerpts for each critical flow.
  - Mapping sheet from request id -> log record -> decision.
- Severity if present: Critical when sensitive actions are not traceable; otherwise High.

## 5) Security Release Gate Checklist (for this threat model)
1. No open Critical security issues.
2. Every in-scope threat has at least one executed verification record and evidence artifact.
3. Any open High issue has owner, mitigation, rollback path, and due date.
4. Medium issues are tracked with target iteration and monitoring signal.
5. Gate decision (`GO`/`NO-GO`) is recorded with reviewer and date.

## 6) Known Current Risks (Repository Snapshot)
1. API currently exposes only health/ping routes, so MVP auth/permission/replay/tampering controls are not yet observable.
2. Request logger currently records method/path/status/bytes/duration/remote but does not explicitly log actor or action semantics.
3. No explicit audit retention/immutability policy is documented yet.
