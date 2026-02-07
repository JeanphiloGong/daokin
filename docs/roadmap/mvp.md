# MVP Roadmap

This roadmap is execution-first and reversible. It follows the current project standards (M1, M2) and the acceptance criteria gate.

## Milestone M1: Narrative + skeleton baseline (Week 1)

### Objective
Establish a shared execution baseline and remove ambiguity.

### Deliverables
- Mission + README aligned (already present)
- Acceptance criteria baseline: `references/acceptance-criteria.md`
- Sub-agent role boundaries: `docs/agents/roles.md`
- API skeleton health checks and router structure (already present)

### Exit criteria
- Team can point to one acceptance baseline file
- Domain terms are consistent across docs and code

### Rollback plan
- If naming/model is unstable, freeze on minimal terms: User, Dao, Artifact, Permission.

## Milestone M2: Core loop implementation (Weeks 2-4)

### Objective
Run the MVP loop end to end:
`create artifact -> join dao -> permissioned exchange -> attribution/reputation`

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
- Membership history (joined_at/left_at)
- No forced assignment path

#### WS-4 Consent and exchange
- Permission grant/revoke records
- Tip/payment linkage to artifact and recipient
- Query endpoint for attribution trail

#### WS-5 Frontend flow
- Wallet auth -> create artifact -> join dao -> tip/pay -> view attribution
- Error states for denied permission and payment failures

### Exit criteria
- All mandatory checks in `references/acceptance-criteria.md` AC-1/2/3 pass
- AC-4 has no critical unresolved risk

### Rollback plan
- If on-chain integration blocks progress, keep contract calls optional and preserve provenance via API records
- If payment integration blocks progress, keep transfer records in pending/mock mode with clear flags

## Milestone M3: Hardening pass (Week 5)

### Objective
Close top risks before broader testing.

### Deliverables
- Threat checklist for auth/permission bypass/replay/tampering
- Manual verification report against acceptance criteria
- Prioritized bug list with owner + due date

### Exit criteria
- No critical ownership/consent/attribution defect open
- Known medium risks have mitigation and owner

### Rollback plan
- Disable non-essential features; keep only core loop endpoints and UI paths.

## Tracking Format (for each PR)
- Intent
- Scope files
- Acceptance criteria touched (AC ids)
- Test evidence (automated/manual)
- Rollback note

## Current Risks and Mitigations
- **Risk:** central dependency lock-in
  - **Mitigation:** keep interfaces storage-agnostic; avoid provider-specific types in domain layer.
- **Risk:** speculation-first drift
  - **Mitigation:** no token incentive features in MVP.
- **Risk:** moderation gap
  - **Mitigation:** publish minimal policy and appeal path before public testing.
