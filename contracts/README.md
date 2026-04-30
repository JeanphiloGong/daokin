# DaoKin Contracts

`contracts` is the Hardhat + Solidity package for DaoKin MVP provenance and attributed value transfer.

The cross-module event interpretation is maintained in [`../docs/architecture/contract-event-map.md`](../docs/architecture/contract-event-map.md). This README is the contracts module entrypoint.

## Boundary

- Owns Solidity contracts, Hardhat configuration, contract tests, and local contract usage notes.
- Provides on-chain provenance and tip primitives for the MVP.
- Does not own backend persistence, frontend state, or project-wide governance policy.

## Contracts

- `ArtifactRegistry`: publish artifacts (URI + content hash) with on-chain provenance.
- `TipJar`: native or ERC20 tips attributed to a specific artifact + recipient.
- `ReentrantRecipient` (test helper): adversarial contract for reentrancy checks in tests.

## Quick start

```bash
npm install
npm run build
```

## Tests

```bash
npm test
```

## ABI (MVP surface)

### ArtifactRegistry

Functions:
- `publish(string uri, bytes32 contentHash) returns (uint256 artifactId)`
- `exists(uint256 artifactId) view returns (bool)`
- `totalArtifacts() view returns (uint256)`
- `artifacts(uint256 artifactId) view returns (address creator, string uri, bytes32 contentHash, uint256 createdAt)`

Errors:
- `InvalidURI()`
- `InvalidContentHash()`

Events:
- `ArtifactPublished(uint256 indexed artifactId, address indexed creator, string uri, bytes32 contentHash, uint256 createdAt)`
  - Provenance fields: artifact id, creator, content pointer/hash, and timestamp.

### TipJar

Functions:
- `tipNative(uint256 artifactId, address to, string memo)` payable
- `tipToken(uint256 artifactId, address token, address to, uint256 amount, string memo)`
- `artifactRegistry() view returns (address)`

Errors:
- `InvalidArtifactId()`
- `ArtifactNotFound()`
- `InvalidRegistry()`
- `InvalidRecipient()`
- `InvalidAmount()`
- `InvalidToken()`
- `ReentrancyBlocked()`
- `TransferFailed()`

Events:
- `Tipped(uint256 indexed artifactId, address indexed from, address indexed to, address token, uint256 amount, string memo)`
  - Attribution fields: artifact id + sender + recipient + token/amount.

## Example calls (ethers v6)

```js
const registry = await ethers.deployContract("ArtifactRegistry");
const tipJar = await ethers.deployContract("TipJar", [await registry.getAddress()]);

const contentHash = ethers.keccak256(ethers.toUtf8Bytes("artifact-v1"));
const publishTx = await registry.publish("ipfs://bafy.../artifact.json", contentHash);
await publishTx.wait();

// Tip in native token (ETH on localnet)
await tipJar.tipNative(
  1n,                      // artifactId
  "0xRecipientAddress",
  "thanks for this artifact",
  { value: ethers.parseEther("0.01") }
);

// Tip in ERC20
await tipJar.tipToken(
  1n,                      // artifactId
  "0xTokenAddress",
  "0xRecipientAddress",
  ethers.parseUnits("5", 6),
  "USDC tip"
);
```

## Notes

- Keep artifact data off-chain (IPFS/Arweave), anchor provenance with `uri + contentHash + createdAt`.
- `TipJar` verifies artifact existence against `ArtifactRegistry` before transfer.
- `TipJar` uses a non-reentrancy guard on native and token tip paths.
- MVP value exchange avoids speculative token logic and prioritizes clear attribution.

## Related Docs

- System overview: [`../docs/system-overview.md`](../docs/system-overview.md)
- Contract event map: [`../docs/architecture/contract-event-map.md`](../docs/architecture/contract-event-map.md)
- MVP roadmap and acceptance baseline: [`../docs/roadmap/mvp.md`](../docs/roadmap/mvp.md)
