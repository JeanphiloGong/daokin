# AGENTS.md (Project Rules)
## Overview
DaoKin is a decentralized community network where artifacts belong to individuals, value flows to contributors, and circles are self-governed by those who walk the same path.
This project aims to cultivate a soil for civilization-grade communities, not a centralized social platform or speculative vehicle.

## Core Principles
- Data belongs to individuals; default is user ownership, not platform custody.
- Value flows to contributors; usage requires permission and fair exchange.
- Communities arise naturally; no forced top-down group formation.
- No algorithmic manipulation that amplifies conflict or extraction.
- Decentralization is a path; avoid hard dependencies that block future self-hosting.
- Governance must be transparent, forkable, and community-first.

## Domain Philosophies (Master-Level)
### Product
- **Goal**: help people find their place and build lasting communities.
- **Constraints**: no growth tactics that distort purpose or community health.
- **Evidence**: clear user flows; documented community feedback.
- **Failure Cost**: the platform becomes another attention trap.
- **Tradeoffs**: reduce feature breadth to protect meaning.
- **Non-negotiables**: no feature that breaks user ownership or consent.

### Design/UX
- **Goal**: reduce cognitive load; make “belonging and contribution” obvious.
- **Constraints**: avoid addictive patterns and dark UX.
- **Evidence**: usability checks for core flows.
- **Failure Cost**: users feel manipulated or lost.
- **Tradeoffs**: clarity over decorative novelty.
- **Non-negotiables**: critical actions must be explicit and reversible.

### Privacy
- **Goal**: data minimization and explicit consent.
- **Constraints**: no default exposure of identity or relationships.
- **Evidence**: documented data flows and access paths.
- **Failure Cost**: loss of trust and safety.
- **Tradeoffs**: slower growth in exchange for privacy.
- **Non-negotiables**: no use of data without clear permission.

### Security
- **Goal**: protect integrity of ownership and value flows.
- **Constraints**: least privilege; safe defaults.
- **Evidence**: threat model notes and audit-ready logs.
- **Failure Cost**: irreversible loss or abuse.
- **Tradeoffs**: reduced convenience to prevent abuse.
- **Non-negotiables**: no secret leakage; no silent privilege escalation.

### Data
- **Goal**: data is portable, verifiable, and attributable.
- **Constraints**: avoid lock-in; preserve provenance.
- **Evidence**: export paths; immutable references where applicable.
- **Failure Cost**: users cannot exit or verify ownership.
- **Tradeoffs**: simpler schemas to ensure portability.
- **Non-negotiables**: no untracked transformations of user data.

### Content Moderation / Governance
- **Goal**: protect safety without central capture.
- **Constraints**: transparent rules and appeals.
- **Evidence**: documented policy and community governance notes.
- **Failure Cost**: harm, abuse, or censorship capture.
- **Tradeoffs**: slower decisions to maintain fairness.
- **Non-negotiables**: no hidden moderation rules.

### Engineering / Architecture
- **Goal**: evolvable architecture that can decentralize over time.
- **Constraints**: avoid irreversible coupling to centralized services.
- **Evidence**: documented interfaces and migration paths.
- **Failure Cost**: inability to self-host or federate.
- **Tradeoffs**: accept slower optimization to keep portability.
- **Non-negotiables**: no hard dependency that prevents future node operation.

**Omitted domains**: Medical, Finance (regulated), ML/AI, Gaming, Supply Chain.
These are not material to the current scope.

## Product & Project Standards
**Goals**
- Enable creation, attribution, and permissioned use of artifacts.
- Enable formation of “The Dao of …” communities.
- Make contribution visible and rewarded.

**Success Metrics (early, MVP-level)**
- A new user can publish an artifact and join a Dao in one session.
- A contributor can receive value for usage with clear attribution.
- Users can export or verify their own data/ownership.

**Scope (current)**
- Philosophy-first, MVP-ready, and community centric.
- Avoid speculative token mechanics until core value exchange works.

**Milestones**
- M1: clear manifesto + basic product narrative.
- M2: MVP core loop (create → join → exchange → reputation).

**Acceptance Criteria (for MVP definition)**
- Ownership is explicit and traceable.
- Exchange requires consent.
- Community formation is organic and reversible.

**Key Risks**
- Drift into centralized extraction.
- Speculation overtaking contribution.
- Moderation failure causing harm.
- Privacy leakage eroding trust.

## 12 Golden Rules (Why / How / Check)
1) **User ownership first**
   - Why: prevents platform capture.
   - How: design flows around user keys and explicit consent.
   - Check: no feature transfers ownership by default.
2) **Consent before usage**
   - Why: value exchange must be ethical.
   - How: permission required for access or reuse.
   - Check: usage paths include consent record.
3) **Contribution over attention**
   - Why: avoid manipulation and empty engagement.
   - How: reward meaningful artifacts and actions.
   - Check: ranking signals exclude raw clicks.
4) **Communities arise, not assigned**
   - Why: organic growth builds trust.
   - How: let users gather around shared purpose.
   - Check: no forced grouping logic.
5) **No dark patterns**
   - Why: autonomy and trust are core.
   - How: avoid loops that exploit anxiety or addiction.
   - Check: UX review for coercive nudges.
6) **Portability by design**
   - Why: prevents lock-in.
   - How: export paths and public proofs.
   - Check: data can be exported or verified.
7) **Transparency in governance**
   - Why: power must be visible.
   - How: publish rules and decision logs.
   - Check: policy changes are documented.
8) **Safety without central capture**
   - Why: harms must be contained without censorship abuse.
   - How: transparent policies and appeal channels.
   - Check: moderation rules are public.
9) **Minimal central dependency**
   - Why: future nodes must be possible.
   - How: keep interfaces and storage portable.
   - Check: no single-service lock-in.
10) **Slow down speculation**
   - Why: speculation corrupts purpose.
   - How: prioritize stable, non-speculative value flows.
   - Check: no token incentives that reward hype.
11) **Evidence before claims**
   - Why: philosophy must be grounded.
   - How: tie claims to observable behavior or data.
   - Check: no “vision-only” claims without roadmap.
12) **Small, reversible steps**
   - Why: safeguard against irreversible mistakes.
   - How: prefer incremental changes.
   - Check: each change has a rollback plan.

## Scope Boundaries
- Allowed: all files inside this repository.
- Forbidden: changes outside the repository; secrets/keys/PII; irreversible history changes.

## Permission Model
- Single owner: the user approves philosophy-level changes.
- No extra approval required for routine edits unless explicitly requested.

## Execution Rules
- Keep changes small and auditable; one intent per change set.
- Update narrative docs when behavior or philosophy changes.
- If adding governance or moderation rules, make them explicit and visible.
- **Iteration Loop (required)**: use `references/acceptance-criteria.md` for pass/fail review when it exists; if missing, record the gap in Risks & Open Questions.
- **Step Gate (required)**: after Plan/Change/Verify/Reflect steps, pause and ask for `continue` before proceeding.

## Quality Bar
- Proof before “done”: list tests run or state “Not run” with reason.
- If no tests exist, provide manual verification notes.
- Avoid any change that violates Core Principles or Golden Rules.

## Decision & Accountability
- Philosophy changes must be traceable in git history.
- If a rule is violated, document the reason and the rollback plan.
- The user is the final authority on vision shifts.

## Risks & Open Questions
- Highest-risk failure: drifting into centralized extraction or speculative finance.
- Acceptance criteria file (`references/acceptance-criteria.md`) not present yet.
- Governance model depth and moderation workflows remain undefined.
- Legal/compliance requirements for value exchange are TBD.
