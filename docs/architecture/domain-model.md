# Domain Model (MVP)

Status: v1 handoff draft for implementation.
Scope: User, Dao, Artifact, Permission, Transfer, Membership.
Alignment: `docs/roadmap/mvp.md` (AC-1/2/3 mandatory, AC-4 baseline).

## 1) Shared Conventions

- IDs: ULID string (`*_id`) for all off-chain entities.
- Timestamps: RFC3339 UTC (`created_at`, `updated_at`, etc.).
- Wallet address: EVM hex address; persist lowercase for stable indexing.
- Amount: `amount_atomic` as decimal string (no floating point).
- Chain linkage is optional and replaceable via `chain_ref`.

`chain_ref` (embedded object when on-chain data exists):

| Field | Type | Required | Notes |
|---|---|---|---|
| `chain_id` | integer | yes | EVM chain id |
| `contract_alias` | string | yes | `artifact_registry` or `tipjar` |
| `contract_address` | string | yes | Deployed contract |
| `tx_hash` | string | yes | 0x tx hash |
| `log_index` | integer | yes | Event log index |
| `block_number` | integer | yes | Source block |
| `block_timestamp` | string | yes | RFC3339 UTC |
| `event_name` | string | yes | Emitted event name |

## 2) Entities

### 2.1 User

| Field | Type | Required | Notes |
|---|---|---|---|
| `user_id` | string | yes | Internal identity |
| `primary_wallet` | string | yes | Auth principal |
| `status` | enum | yes | `active`, `suspended` |
| `created_at` | string | yes | Creation time |
| `updated_at` | string | yes | Last update time |

State machine:
- `active -> suspended`
- `suspended -> active`

Constraints:
- `primary_wallet` unique.
- New user starts with zero memberships (no forced Dao assignment, AC-3.2).

### 2.2 Dao

| Field | Type | Required | Notes |
|---|---|---|---|
| `dao_id` | string | yes | Internal id |
| `slug` | string | yes | Unique URL-safe id |
| `name` | string | yes | Display name |
| `description` | string | no | Optional text |
| `creator_user_id` | string | yes | Creator |
| `status` | enum | yes | `active`, `archived` |
| `created_at` | string | yes | Creation time |
| `updated_at` | string | yes | Last update time |

State machine:
- `active -> archived` (terminal for MVP)

Constraints:
- `slug` unique.
- Archiving a Dao does not auto-delete historical memberships.

### 2.3 Artifact

| Field | Type | Required | Notes |
|---|---|---|---|
| `artifact_id` | string | yes | Internal id |
| `creator_user_id` | string | yes | Owner user id |
| `creator_wallet` | string | yes | AC-1.1 required |
| `dao_id` | string | no | Optional circle context |
| `title` | string | no | Optional display |
| `content_uri` | string | yes | AC-1.1 required |
| `content_hash` | string | yes | AC-1.1 required, `0x` bytes32 string |
| `status` | enum | yes | `active`, `withdrawn` |
| `provenance_status` | enum | yes | `offchain_recorded`, `onchain_confirmed` |
| `chain_artifact_id` | string | no | On-chain uint256 artifact id (decimal string) |
| `chain_ref` | object | no | Required when `provenance_status=onchain_confirmed` |
| `created_at` | string | yes | AC-1.1 required |
| `updated_at` | string | yes | Last update time |

State machine:
- `status`: `active -> withdrawn`
- `provenance_status`: `offchain_recorded -> onchain_confirmed`

Constraints:
- `creator_wallet`, `content_uri`, `content_hash`, `created_at` are immutable after creation.
- If on-chain event is linked, `chain_ref` must be complete and immutable.

### 2.4 Permission

| Field | Type | Required | Notes |
|---|---|---|---|
| `permission_id` | string | yes | Internal id |
| `artifact_id` | string | yes | Target artifact |
| `granter` | string | yes | AC-2.1 required |
| `grantee` | string | yes | AC-2.1 required |
| `scope` | enum | yes | `view`, `reuse`, `commercial_use` |
| `status` | enum | yes | `granted`, `revoked`, `expired` |
| `granted_at` | string | yes | AC-2.1 required |
| `expires_at` | string | no | Optional TTL |
| `revoked_at` | string | no | Set when revoked |
| `revoke_reason` | string | no | Optional reason |
| `created_at` | string | yes | Creation time |
| `updated_at` | string | yes | Last update time |

State machine:
- `granted -> revoked`
- `granted -> expired` (time driven)
- `revoked` and `expired` are terminal

Constraints:
- `granter` must equal artifact `creator_wallet`.
- Only one active (`status=granted`, non-expired) record for same (`artifact_id`, `grantee`, `scope`).
- Protected usage checks only pass on active permission (AC-2.1).

### 2.5 Transfer

| Field | Type | Required | Notes |
|---|---|---|---|
| `transfer_id` | string | yes | Internal id |
| `transfer_intent_id` | string | yes | Pre-chain correlation key |
| `artifact_id` | string | yes | AC-2.2 required |
| `permission_id` | string | no | Linked consent record |
| `sender_wallet` | string | yes | Payer |
| `recipient_wallet` | string | yes | AC-2.2 required |
| `token_address` | string | yes | `0x000...000` means native token |
| `amount_atomic` | string | yes | Integer string |
| `memo` | string | no | Tip memo / correlation hint |
| `tx_hash` | string | yes | AC-2.2 required when submitted |
| `tx_status` | enum | yes | `initiated`, `confirmed`, `failed`, `reorged` |
| `chain_ref` | object | no | Required after on-chain observation |
| `created_at` | string | yes | Record created |
| `confirmed_at` | string | no | Set on final confirmation |

State machine:
- `initiated -> confirmed`
- `initiated -> failed`
- `confirmed -> reorged` (rare)

Constraints:
- `artifact_id + tx_hash` must be queryable (AC-2.2).
- `recipient_wallet` defaults to artifact creator wallet in MVP.
- Transfer list must support query by `artifact_id` (AC-2.2 pass check).

### 2.6 Membership

| Field | Type | Required | Notes |
|---|---|---|---|
| `membership_id` | string | yes | Internal id |
| `dao_id` | string | yes | Target Dao |
| `user_id` | string | yes | Member user |
| `role` | enum | yes | `member` (default), `steward` |
| `status` | enum | yes | `active`, `left` |
| `joined_at` | string | yes | AC-3.1 required |
| `left_at` | string | no | AC-3.1 required when left |
| `left_reason` | string | no | AC-3.1 optional reason |
| `created_at` | string | yes | Creation time |
| `updated_at` | string | yes | Last update time |

State machine:
- `active -> left`
- Rejoin creates a new membership record (preserves audit history)

Constraints:
- Max one `active` membership per (`dao_id`, `user_id`).
- Join and leave must both be user-triggerable (AC-3.1).

## 3) Cross-Entity MVP Invariants

1. Artifact provenance fields are always readable via API: `creator_wallet`, `content_uri`, `content_hash`, `created_at` (AC-1.1).
2. Permissioned use requires explicit permission record (`granter`, `grantee`, `scope`, `status`, `granted_at`) (AC-2.1).
3. Transfer attribution record links `tx_hash`, `artifact_id`, `recipient_wallet` and is queryable per artifact (AC-2.2).
4. Onboarding does not auto-create membership records (AC-3.2).
5. All mutating operations emit audit log metadata: request id, actor wallet/user, action, timestamp (AC-4.2).

## 4) Portability and Replaceability Rules

- Domain IDs are chain-agnostic; on-chain ids are stored as optional references.
- Contract names are stored as aliases (`artifact_registry`, `tipjar`) so addresses can be swapped per environment.
- Contract event ingestion is adapter-based; domain schema does not depend on ABI internals.
- If on-chain integration is unavailable, `provenance_status=offchain_recorded` still satisfies core loop, with clear flags.
