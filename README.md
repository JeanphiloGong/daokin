# DaoKin（道友）

**DaoKin — Where kin become Daos.**  
**同道者聚，诸道自生。**

DaoKin is a public experimental project exploring a user-owned model for creative work and communities. It treats an **Artifact** as a contribution that can be attributed, permissioned, and connected to value, while a **Dao** is a circle formed around a shared path or purpose.

The project starts deliberately small: establish a verifiable artifact record and an attributable value-exchange primitive first, then evolve toward more portable and self-governed communities. The current repository is an early MVP foundation, not a production platform or a finished decentralized network.

> 中文简介：DaoKin（道友）探索“用户拥有的作品与共同体”。作品应当能够被归属、授权和结算价值，社区则围绕共同的道路与目的自然形成。当前仓库聚焦最小可验证闭环，先验证作品登记与可追溯的价值交换，再逐步推进自治与可迁移。

## At a Glance

| Area | Current position |
| --- | --- |
| Product focus | User-owned artifacts, attribution, permissioned use, and purpose-driven circles |
| Current stage | Early prototype / MVP foundation |
| Web client | SvelteKit scaffold; UI exploration has not yet been connected to domain flows |
| API | Go + Chi service skeleton with health, readiness, and ping endpoints |
| On-chain layer | Minimal Solidity contracts for artifact registration and tipping |
| Contract status | Not audited; no production or mainnet deployment is claimed |
| License | TBD |

## Product Model

### Artifacts

An Artifact is off-chain content, such as a document, design, or piece of research, paired with an on-chain reference: a URI and a content hash. The reference is intended to make authorship and provenance easier to verify without placing the content itself on-chain.

### Daos

“The Dao of …” names a circle organized around a domain, practice, or shared purpose, such as **The Dao of Creation**. Dao membership, governance, and portability are product directions currently captured in the project plan; they are not implemented in this repository yet.

### Permissioned use

Using or reusing an Artifact should be explicit, attributable, and compatible with the creator’s terms. Consent records, licenses, and value splitting belong to the planned product layer rather than the current contract minimum.

## What Exists Today

### Web client: `daokin-web/`

SvelteKit application scaffold for the future product interface. It currently provides the frontend development, checking, testing, and build toolchain; the default page has not yet been replaced with the DaoKin domain experience.

### API: `daokin-api/`

Go 1.22 service built with Chi. The current HTTP boundary includes:

- `GET /healthz` — liveness response
- `GET /readyz` — readiness response
- `GET /v1/ping` — versioned connectivity check

The server also includes request IDs, real-IP extraction, request logging, panic recovery, request timeouts, configurable bind address, and graceful shutdown. Domain APIs and persistence are future work.

### Contracts: `contracts/`

Hardhat project targeting Solidity `0.8.20` with two minimal contracts:

- `ArtifactRegistry` records a creator, URI, content hash, and timestamp, then emits an `ArtifactPublished` event.
- `TipJar` transfers native currency or an ERC-20 token directly to a recipient and emits a `Tipped` event with an optional memo.

These contracts are intentionally small experiments. They do not yet implement permissions, licensing, governance, revenue splits, or upgradeability.

The modules currently build independently. Application wiring between the web client, API, storage, wallet, and contracts is part of the next implementation stage.

## Repository Structure

```text
daokin/
├── daokin-web/       # SvelteKit client scaffold
├── daokin-api/       # Go + Chi HTTP service skeleton
├── contracts/        # Hardhat project and Solidity contracts
├── MISSION.md        # Long-term vision and manifesto
├── MVP_PLAN.md       # MVP scope, acceptance criteria, and roadmap
└── AGENTS.md         # Product and engineering guardrails
```

## Quick Start

The three modules run independently during this prototype stage. Use separate terminals when running more than one of them.

### Web client

```bash
cd daokin-web
npm install
npm run dev
```

Useful checks:

```bash
npm run check
npm run build
```

### API

Requires Go 1.22 or later.

```bash
cd daokin-api
go run ./cmd/api
```

By default the server listens on `:8080`:

```bash
curl http://localhost:8080/healthz
curl http://localhost:8080/readyz
curl http://localhost:8080/v1/ping
```

Configuration is documented in [`daokin-api/README.md`](daokin-api/README.md) and can be supplied through environment variables such as `DAOKIN_ADDR` and `DAOKIN_ENV`.

### Contracts

```bash
cd contracts
npm install
npm run build
```

The contracts have no deployment target or audit report in this repository. Do not use them with real funds.

## Development Direction

The next product slices are planned around the core loop:

1. Replace the web starter page with artifact and circle flows.
2. Define API contracts and persistence for artifacts, ownership proofs, and membership intent.
3. Connect wallet identity and off-chain storage such as IPFS or Arweave.
4. Add explicit permission, attribution, and settlement records.
5. Add contract tests and a security review before considering any public deployment.
6. Explore portable communities and governance only after the core ownership and exchange flows are clear.

The project intentionally avoids token speculation, attention-optimizing feeds, forced community formation, and irreversible central dependencies. See [`MVP_PLAN.md`](MVP_PLAN.md) for scope and acceptance criteria, and [`MISSION.md`](MISSION.md) for the long-term rationale.

## Project Status and Safety

- Breaking changes are expected while the MVP boundaries are being tested.
- The smart contracts are not audited.
- There is no claim of mainnet deployment, production usage, or real-value settlement.
- Legal and compliance requirements for value exchange remain open questions.
- The repository does not yet declare a license; `License: TBD` should be resolved before treating it as a reusable open-source package.

## Project Rules

Read [`AGENTS.md`](AGENTS.md) before changing product principles, ownership semantics, privacy assumptions, or governance rules. The subprojects also contain their own engineering guidance:

- [`daokin-web/AGENTS.md`](daokin-web/AGENTS.md)
- [`daokin-api/AGENTS.md`](daokin-api/AGENTS.md)
