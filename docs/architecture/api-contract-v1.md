# API Contract v1 (MVP)

Status: implementation contract for backend/frontend handoff.
Alignment: `references/acceptance-criteria.md` AC-1/2/3 mandatory, AC-4 audit baseline.

## 1) Global Conventions

- Base path: `/v1`
- Content type: `application/json`
- Auth: `Authorization: Bearer <access_token>`
- Request id: server returns `X-Request-Id` on every response
- Time format: RFC3339 UTC
- Wallet format: lowercase EVM address
- Amount format: `amount_atomic` decimal string

## 2) Response and Error Shape

Success:

```json
{
  "data": {}
}
```

Error:

```json
{
  "error": {
    "code": "INVALID_ARGUMENT",
    "message": "content_hash must be bytes32 hex",
    "request_id": "req_01J...",
    "details": {}
  }
}
```

Common error codes:

| HTTP | Code | Meaning |
|---|---|---|
| 400 | `INVALID_ARGUMENT` | Bad field value or missing field |
| 401 | `UNAUTHORIZED` | Missing/invalid token |
| 403 | `FORBIDDEN` | Authenticated but not allowed |
| 404 | `NOT_FOUND` | Resource missing |
| 409 | `CONFLICT` | State conflict (already joined, duplicate permission) |
| 422 | `RULE_VIOLATION` | Business rule failed |
| 429 | `RATE_LIMITED` | Too many requests |
| 500 | `INTERNAL` | Internal failure |

Domain-specific codes (subset):
- `CHALLENGE_EXPIRED`, `SIGNATURE_INVALID`, `WALLET_MISMATCH`
- `DAO_NOT_FOUND`, `MEMBERSHIP_NOT_ACTIVE`
- `ARTIFACT_NOT_FOUND`, `PERMISSION_NOT_FOUND`
- `PERMISSION_CONFLICT`, `FORBIDDEN_NOT_OWNER`, `FORBIDDEN_NOT_GRANTER`

## 3) Endpoint Contracts

### 3.1 Auth challenge

`POST /v1/auth/challenge`

Request:

```json
{
  "wallet": "0xabc...",
  "chain_id": 1,
  "statement": "Sign in to DaoKin"
}
```

Response:

```json
{
  "data": {
    "challenge_id": "chl_01J...",
    "wallet": "0xabc...",
    "nonce": "n_01J...",
    "message": "DaoKin login nonce: n_01J...",
    "expires_at": "2026-02-07T10:10:00Z"
  }
}
```

Errors: `INVALID_ARGUMENT`, `RATE_LIMITED`.

### 3.2 Auth verify

`POST /v1/auth/verify`

Request:

```json
{
  "challenge_id": "chl_01J...",
  "wallet": "0xabc...",
  "signature": "0x..."
}
```

Response:

```json
{
  "data": {
    "access_token": "jwt_or_paseto",
    "token_type": "Bearer",
    "expires_at": "2026-02-07T22:10:00Z",
    "user": {
      "user_id": "usr_01J...",
      "primary_wallet": "0xabc...",
      "status": "active",
      "created_at": "2026-02-07T10:00:00Z"
    }
  }
}
```

Errors: `CHALLENGE_EXPIRED`, `SIGNATURE_INVALID`, `WALLET_MISMATCH`, `UNAUTHORIZED`.

### 3.3 Artifact create

`POST /v1/artifacts` (auth required)

Request:

```json
{
  "title": "My first artifact",
  "content_uri": "ipfs://bafy...",
  "content_hash": "0x1234...",
  "dao_id": "dao_01J...",
  "publish_onchain": true
}
```

Response:

```json
{
  "data": {
    "artifact_id": "art_01J...",
    "chain_artifact_id": "42",
    "creator_wallet": "0xabc...",
    "content_uri": "ipfs://bafy...",
    "content_hash": "0x1234...",
    "status": "active",
    "provenance_status": "offchain_recorded",
    "created_at": "2026-02-07T10:20:00Z"
  }
}
```

Errors: `UNAUTHORIZED`, `INVALID_ARGUMENT`, `DAO_NOT_FOUND`, `CONFLICT`.

### 3.4 Artifact get

`GET /v1/artifacts/{artifact_id}`

Response:

```json
{
  "data": {
    "artifact_id": "art_01J...",
    "chain_artifact_id": "42",
    "creator_wallet": "0xabc...",
    "content_uri": "ipfs://bafy...",
    "content_hash": "0x1234...",
    "status": "active",
    "provenance_status": "onchain_confirmed",
    "chain_ref": {
      "chain_id": 1,
      "contract_alias": "artifact_registry",
      "contract_address": "0xreg...",
      "tx_hash": "0xtx...",
      "log_index": 3,
      "block_number": 21900001,
      "block_timestamp": "2026-02-07T10:21:00Z"
    },
    "created_at": "2026-02-07T10:20:00Z"
  }
}
```

Errors: `NOT_FOUND`.

### 3.5 Dao join

`POST /v1/daos/{dao_id}/join` (auth required)

Request:

```json
{
  "reason": "Interested in public goods"
}
```

Response:

```json
{
  "data": {
    "membership_id": "mem_01J...",
    "dao_id": "dao_01J...",
    "user_id": "usr_01J...",
    "role": "member",
    "status": "active",
    "joined_at": "2026-02-07T10:30:00Z"
  }
}
```

Errors: `UNAUTHORIZED`, `DAO_NOT_FOUND`, `CONFLICT`.

### 3.6 Dao leave

`POST /v1/daos/{dao_id}/leave` (auth required)

Request:

```json
{
  "left_reason": "Switching focus"
}
```

Response:

```json
{
  "data": {
    "membership_id": "mem_01J...",
    "dao_id": "dao_01J...",
    "user_id": "usr_01J...",
    "status": "left",
    "joined_at": "2026-02-07T10:30:00Z",
    "left_at": "2026-02-07T12:30:00Z",
    "left_reason": "Switching focus"
  }
}
```

Errors: `UNAUTHORIZED`, `DAO_NOT_FOUND`, `MEMBERSHIP_NOT_ACTIVE`.

### 3.7 Permission grant

`POST /v1/permissions` (auth required)

Request:

```json
{
  "artifact_id": "art_01J...",
  "grantee": "0xdef...",
  "scope": "reuse",
  "expires_at": "2026-03-01T00:00:00Z"
}
```

Response:

```json
{
  "data": {
    "permission_id": "prm_01J...",
    "artifact_id": "art_01J...",
    "granter": "0xabc...",
    "grantee": "0xdef...",
    "scope": "reuse",
    "status": "granted",
    "granted_at": "2026-02-07T11:00:00Z",
    "expires_at": "2026-03-01T00:00:00Z"
  }
}
```

Errors: `UNAUTHORIZED`, `ARTIFACT_NOT_FOUND`, `FORBIDDEN_NOT_OWNER`, `PERMISSION_CONFLICT`, `INVALID_ARGUMENT`.

### 3.8 Permission revoke

`POST /v1/permissions/{permission_id}/revoke` (auth required)

Request:

```json
{
  "revoke_reason": "Terms changed"
}
```

Response:

```json
{
  "data": {
    "permission_id": "prm_01J...",
    "status": "revoked",
    "revoked_at": "2026-02-07T11:10:00Z",
    "revoke_reason": "Terms changed"
  }
}
```

Errors: `UNAUTHORIZED`, `PERMISSION_NOT_FOUND`, `FORBIDDEN_NOT_GRANTER`, `CONFLICT`.

### 3.9 Transfer query (artifact scoped)

`GET /v1/artifacts/{artifact_id}/transfers?status=confirmed&cursor=...&limit=20`

Response:

```json
{
  "data": {
    "items": [
      {
        "transfer_id": "trf_01J...",
        "artifact_id": "art_01J...",
        "permission_id": "prm_01J...",
        "sender_wallet": "0xdef...",
        "recipient_wallet": "0xabc...",
        "token_address": "0x0000000000000000000000000000000000000000",
        "amount_atomic": "10000000000000000",
        "tx_hash": "0xtx...",
        "tx_status": "confirmed",
        "memo": "dk:v1:intent:ti_01J...",
        "confirmed_at": "2026-02-07T11:30:00Z"
      }
    ],
    "next_cursor": null
  }
}
```

Errors: `INVALID_ARGUMENT`, `NOT_FOUND`.

## 4) AC Coverage Matrix

- AC-1.1 provenance fields: covered by `POST /v1/artifacts` and `GET /v1/artifacts/{artifact_id}`.
- AC-2.1 consent record: covered by `POST /v1/permissions` and revoke endpoint.
- AC-2.2 attribution query: covered by `GET /v1/artifacts/{artifact_id}/transfers`.
- AC-3.1 join/leave reversibility: covered by join and leave endpoints with audit timestamps.
- AC-3.2 no forced grouping: no endpoint auto-creates membership during auth/onboarding.
- AC-4.2 audit baseline: `X-Request-Id` + authenticated actor in request logs for all write paths.

## 5) Optional v1 Extension (for AC-1.2 export)

To close AC-1.2 fully, add:
- `POST /v1/users/me/export-jobs`
- `GET /v1/users/me/export-jobs/{job_id}`

Output should include only the caller's artifacts and memberships.
