# System Overview

DaoKin is an MVP-stage repository for user-owned artifacts, permissioned use, and community formation around "The Dao of ..." circles. The current system is intentionally small: it proves the core loop before committing to deeper decentralization or wider governance machinery.

This page is a routing overview. The detailed authority for API shape, domain entities, contract events, and delivery state lives in the linked source documents.

## System Shape

- `daokin-web/`: SvelteKit web prototype for the MVP flow.
- `daokin-api/`: Go + Chi API service for auth, artifact records, Dao membership, permission records, transfer records, attribution queries, and user export.
- `contracts/`: Hardhat + Solidity contracts for artifact provenance and attributed tips.
- `docs/`: system-level architecture, governance, roadmap, security, and cross-module contracts.
- `references/`: bounded acceptance baselines, verification templates, and run evidence.

## MVP Core Loop

The MVP loop is:

```text
wallet auth -> create artifact -> join dao -> grant permission -> record transfer -> view attribution -> export user data
```

The loop is governed by [../references/acceptance-criteria.md](../references/acceptance-criteria.md). The delivery plan is tracked in [roadmap/mvp.md](roadmap/mvp.md).

## Module Boundaries

The frontend owns user interaction, local UI state, and the demo flow. It should treat shared API behavior as a contract defined by [architecture/api-contract-v1.md](architecture/api-contract-v1.md).

The backend owns off-chain MVP state, authorization checks, permission enforcement, attribution queries, and user export. It should keep provider-specific infrastructure behind application ports so the service can evolve toward portable storage and future node operation.

The contracts package owns on-chain provenance and value-attribution primitives. Cross-module event interpretation belongs in [architecture/contract-event-map.md](architecture/contract-event-map.md), not in a module-only note.

## Current Authority Docs

- Domain entities and invariants: [architecture/domain-model.md](architecture/domain-model.md)
- HTTP API contract: [architecture/api-contract-v1.md](architecture/api-contract-v1.md)
- Contract event mapping: [architecture/contract-event-map.md](architecture/contract-event-map.md)
- MVP roadmap: [roadmap/mvp.md](roadmap/mvp.md)
- Threat model and security gate: [security/threat-model-mvp.md](security/threat-model-mvp.md)
- Documentation governance: [governance/documentation-governance.md](governance/documentation-governance.md)

## Current Delivery State

The repository is in the MVP implementation phase. The backend has the main MVP route surface and in-memory persistence. The frontend is still a prototype flow and must be checked against real API integration before MVP acceptance. The contracts package has minimal provenance and tip primitives, but the contracts are not audited.

MVP acceptance requires the checks in [../references/acceptance-criteria.md](../references/acceptance-criteria.md), with execution evidence recorded through [../references/verification/mvp-iteration-01.md](../references/verification/mvp-iteration-01.md).

## Operating Constraints

- User ownership, consent, portability, and transparent governance are project constraints, not optional product polish.
- Root docs own cross-module truth. Module README files own local orientation and commands.
- Module-local docs should be created only when the module has durable local details that exceed its README.
- Verification artifacts belong in `references/` unless they become durable architecture, policy, or contract knowledge.
