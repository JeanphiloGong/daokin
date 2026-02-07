# Sub-agent Roles (MVP)

This document defines execution roles for parallel work while keeping changes auditable and reversible.

## Operating Rules
- One role = one clear responsibility boundary.
- Each role must produce verifiable outputs (files, endpoints, tests, docs).
- No role changes philosophy/mission text without owner approval.
- Cross-role handoff must reference concrete artifacts (path, endpoint, ABI, schema).

## Role 1: Architecture Agent

### Scope
- Domain model and interface contracts.
- API/contract boundary, naming, and versioning strategy.

### Owns
- `docs/architecture/*`
- API surface conventions
- Event naming conventions for contracts

### Outputs
- Domain model: User, Dao, Artifact, Permission, Transfer, Membership
- API contract draft (request/response fields)
- Contract integration contract (ABI usage map)

### Done definition
- Backend and frontend can implement without field ambiguity.

## Role 2: Backend Agent (Go + chi)

### Scope
- HTTP API, business logic orchestration, persistence adapters.

### Owns
- `daokin-api/*`

### Outputs
- Auth/session endpoints (wallet challenge + verify)
- Artifact CRUD (MVP subset)
- Dao membership join/leave
- Permission grant/revoke checks
- Transfer attribution query endpoints

### Done definition
- Endpoints run locally and pass manual acceptance checks in `references/acceptance-criteria.md`.

## Role 3: Contract Agent (Solidity)

### Scope
- Minimal on-chain primitives for provenance and tipping.

### Owns
- `contracts/contracts/*`
- `contracts/scripts/*` (if added)
- `contracts/test/*` (if added)

### Outputs
- Artifact registry event for provenance pointer/hash
- Tip/payment contract integration points
- ABI artifacts for frontend/backend integration

### Done definition
- Contracts compile and expose required events/functions for MVP flow.

## Role 4: Frontend Agent (Svelte)

### Scope
- User flow and API/chain integration UX.

### Owns
- `daokin-web/*`

### Outputs
- Wallet connect/auth flow
- Publish artifact screen
- Dao join/leave interaction
- Tip/payment action UI
- Activity/provenance display

### Done definition
- A new user can complete `create -> join -> exchange -> verify attribution` in one session.

## Role 5: Security & Acceptance Agent

### Scope
- Threat review, quality gates, release readiness.

### Owns
- `references/acceptance-criteria.md`
- `docs/security/*` (if added)
- verification notes in PR/release docs

### Outputs
- Threat checklist for auth, permission bypass, replay, tampering
- Verification report mapped to acceptance criteria
- Open risks with severity and mitigation owner

### Done definition
- No critical unresolved issue in ownership, consent, or attribution path.

## Handoff Matrix
- Architecture -> Backend/Frontend/Contract: schema + naming + version contracts
- Contract -> Backend/Frontend: ABI + event semantics + network config
- Backend -> Frontend: endpoint docs + error codes + pagination/filter semantics
- Security -> All: fail gates + required fixes before merge/release

## Cadence (recommended)
- Daily sync: blockers and contract changes only
- Merge strategy: small PRs, one intent per change
- Verification gate: run manual protocol in `references/acceptance-criteria.md` before milestone tag
