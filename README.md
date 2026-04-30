# DaoKin（道友）

**DaoKin — Where kin become Daos.**  
**同道者聚，诸道自生。**

DaoKin（中文名：道友）是一个开源项目，探索“用户拥有的作品与共同体”：
- 作品（Artifacts）归创作者所有，并可被授权使用与结算价值
- 圈子以 “The Dao of …” 的形式组织，逐步走向自治与可迁移

本仓库目前以 **MVP / 原型** 为目标：先跑通核心闭环，再逐步去中心化（而不是一开始就做全网分布式）。  
更宏观的愿景与宣言已移至 `MISSION.md`。

## Repo Structure
- `daokin-web/`: SvelteKit web app (UI prototype)
- `daokin-api/`: Go + Chi API skeleton
- `contracts/`: Hardhat + Solidity minimal contracts

## Core Concepts
- **Artifact**: off-chain content (e.g. IPFS/Arweave) with an on-chain proof (hash/URI).
- **Dao (circle)**: a community organized as “The Dao of …” (e.g. The Dao of Creation).
- **Permissioned use**: using/reusing an artifact should be explicit and attributable.

## Status
- Early-stage prototype; breaking changes expected.
- Smart contracts are **not audited**.

## Quick Start
### Web
```bash
cd daokin-web
npm install
npm run dev
```

### API
```bash
cd daokin-api
go run ./cmd/server
```

### Contracts
```bash
cd contracts
npm install
npm run build
```

## Docs
- Documentation landing and reading paths: `docs/README.md`
- System overview: `docs/system-overview.md`
- Mission / long-term vision: `MISSION.md`
- Project rules / guardrails: `AGENTS.md`
- Module entrypoints: `daokin-api/README.md`, `daokin-web/README.md`, `contracts/README.md`
- MVP roadmap, acceptance criteria, and verification template: `docs/roadmap/mvp.md`

## License
TBD.
