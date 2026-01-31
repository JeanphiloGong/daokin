# DaoKin Contracts

Minimal contracts to bootstrap DaoKin MVP.

## Contracts

- `ArtifactRegistry`: publish artifacts (URI + content hash) on-chain.
- `TipJar`: native or ERC20 tips with a memo.

## Quick start

```bash
npm install
npm run build
```

## Notes

- Store artifact content off-chain (IPFS/Arweave), keep hashes on-chain.
- Use stablecoins (USDC/DAI) for value exchange in MVP.
