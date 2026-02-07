# Contract Event Map (ArtifactRegistry / TipJar)

Status: MVP mapping spec for backend ingestion and traceability.
Goal: map contract events to backend canonical fields without hard lock-in.

## 1) Source Contracts and Events

From `contracts/contracts/ArtifactRegistry.sol`:

- `ArtifactPublished(uint256 artifactId, address creator, string uri, bytes32 contentHash)`

From `contracts/contracts/TipJar.sol`:

- `Tipped(address from, address to, address token, uint256 amount, string memo)`

## 2) Canonical Backend Chain Fields

Use the same chain metadata for all ingested events:

- `chain_id`
- `contract_alias`
- `contract_address`
- `tx_hash`
- `log_index`
- `block_number`
- `block_timestamp`
- `event_name`

These are stored inside domain `chain_ref` or transfer metadata.

## 3) ArtifactRegistry -> Artifact Mapping

| Contract source | Backend field | Rule |
|---|---|---|
| `artifactId` | `artifact.chain_artifact_id` | Store as decimal string |
| `creator` | `artifact.creator_wallet` | Must equal existing creator wallet |
| `uri` | `artifact.content_uri` | Must match off-chain record |
| `contentHash` | `artifact.content_hash` | Normalize to 0x bytes32 |
| tx receipt | `artifact.chain_ref.tx_hash` | Required when linked |
| log metadata | `artifact.chain_ref.*` | Fill alias/address/log/block/time |
| derived | `artifact.provenance_status` | Set to `onchain_confirmed` |

Ingestion conflict rule:
- If `creator`, `uri`, or `content_hash` mismatches backend record, reject linkage and create an audit alert.

## 4) TipJar -> Transfer Mapping

| Contract source | Backend field | Rule |
|---|---|---|
| `from` | `transfer.sender_wallet` | Lowercase for indexing |
| `to` | `transfer.recipient_wallet` | Must match transfer intent recipient |
| `token` | `transfer.token_address` | Zero address means native token |
| `amount` | `transfer.amount_atomic` | Decimal string |
| `memo` | `transfer.memo` | Should contain `transfer_intent_id` |
| tx receipt | `transfer.tx_hash` | Required for AC-2.2 |
| log metadata | `transfer.chain_ref.*` | Fill alias/address/log/block/time |
| derived | `transfer.tx_status` | `initiated` then `confirmed` after finality |

## 5) Correlation Strategy (Required for AC-2.2)

`Tipped` does not include `artifact_id`, so backend must correlate using a pre-created transfer intent.

Required flow:
1. Before wallet signs tx, backend creates `transfer_intent_id` with `artifact_id`, optional `permission_id`, recipient, token, amount.
2. Client puts `transfer_intent_id` into `memo`, format: `dk:v1:intent:<transfer_intent_id>`.
3. Indexer parses memo and resolves intent.
4. Backend writes `Transfer` record with mandatory link: `tx_hash + artifact_id + recipient_wallet`.

Failure rule:
- If intent cannot be resolved, do not create final Transfer record; store the event in a dead-letter queue for manual reconciliation.

## 6) Reorg and Finality Handling

- On first event observation: create/update record as `tx_status=initiated`.
- After `N` confirmations (configurable): set `tx_status=confirmed`.
- If log is removed by reorg: set `tx_status=reorged` and keep audit trail.

## 7) Replaceability Design

- Use contract aliases (`artifact_registry`, `tipjar`) instead of hard-coded addresses in domain logic.
- Event decoding is behind adapter interface; domain services consume normalized event objects only.
- New contracts can be introduced by adding an adapter and map table, without changing API fields.

## 8) Acceptance Criteria Traceability

- AC-1.1: `ArtifactPublished` provides creator + uri/hash pointer linked to API artifact record.
- AC-2.2: `Tipped` plus transfer intent provides `tx_hash + artifact_id + recipient_wallet`.
- AC-3/AC-4 impact: no direct membership semantics, but all ingestion actions must emit request/job audit logs.
