# Acceptance Criteria (MVP)

This file is the pass/fail baseline for the required iteration loop.
Scope is limited to MVP core loop: `create -> join -> exchange -> reputation`.

## 1) Ownership is explicit and traceable

### AC-1.1 Artifact provenance
- **Given** a user publishes an artifact
- **When** the artifact is created
- **Then** system stores creator wallet, content hash/URI, and created timestamp
- **Then** provenance can be read via API and (if on-chain) contract event/log

**Pass check**
- API response includes `creator_wallet`, `content_uri`, `content_hash`, `created_at`
- Contract event includes artifact id + creator + hash/uri pointer

### AC-1.2 Data export path
- **Given** a user requests export
- **When** export endpoint/job runs
- **Then** user can retrieve their artifacts and memberships in machine-readable format

**Pass check**
- Export endpoint or script exists and documented
- Sample export contains user-owned records only

## 2) Exchange requires consent

### AC-2.1 Permissioned use record
- **Given** another user wants to use an artifact
- **When** permission is granted
- **Then** there is an explicit consent record with actor, scope, and timestamp

**Pass check**
- API/DB has a permission object with `granter`, `grantee`, `scope`, `status`, `granted_at`
- No protected usage path bypasses permission check

### AC-2.2 Value attribution
- **Given** value transfer happens for artifact usage (tip/payment)
- **When** payment succeeds
- **Then** transaction references artifact and recipient creator

**Pass check**
- Transfer record links `tx_hash` + `artifact_id` + `recipient_wallet`
- UI/API can query all transfers for one artifact

## 3) Communities form organically and are reversible

### AC-3.1 Join/leave freedom
- **Given** a user joins a Dao circle
- **When** user decides to leave
- **Then** user can leave without hidden lock-in

**Pass check**
- Join and leave endpoints exist
- Membership state changes are auditable (`joined_at`, `left_at`, reason optional)

### AC-3.2 No forced grouping
- **Given** a newly registered user
- **When** onboarding completes
- **Then** no automatic mandatory Dao assignment occurs

**Pass check**
- Onboarding flow has optional Dao join
- No backend default that auto-creates mandatory memberships

## 4) Safety and transparency minimum

### AC-4.1 Visible governance/moderation entry
- Public docs include moderation rules, reporting path, and appeal path (minimal version allowed)

### AC-4.2 Auditability baseline
- API request logs include request id, actor id/wallet (if authenticated), action, timestamp

## 5) Manual verification protocol

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
- issue link (if fail)

## 6) Exit conditions

MVP is accepted only if:
- All AC-1/2/3 mandatory checks pass
- AC-4 has no critical gaps
- Known failures are documented with rollback/mitigation

## 7) Out of scope (MVP)

- Token speculation mechanics
- Complex DAO voting systems
- Full decentralized node federation
- Regulated financial compliance automation
