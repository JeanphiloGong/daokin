---
id: SPEC-API-001
title: API Contract v1
type: spec
level: system
domain: shared
status: active
owner: project-owner
created_at: 2026-04-07
updated_at: 2026-04-30
last_verified_at: 2026-04-30
review_by: 2026-07-30
version: v2
source_of_truth: true
related_issues: [1, 2, 3, 4, 5]
related_docs:
  - ../roadmap/mvp.md
  - domain-model.md
  - contract-event-map.md
supersedes: []
superseded_by: []
tags: [api, mvp, contract]
---

# API Contract v1

This contract describes the HTTP API currently implemented by `daokin-api` for the MVP loop. It is the shared backend/frontend handoff for M2 execution and should stay aligned with the router and handlers in `daokin-api/internal/interfaces/http/`.

## Conventions

- Base path: `/v1`
- Content type: `application/json`
- Auth header: `Authorization: Bearer <access_token>`
- Time format: RFC3339 UTC
- Wallet format: EVM address; clients should send lowercase addresses
- Amount format: `amount_atomic` decimal string
- Success responses return the resource directly, not a `{ "data": ... }` wrapper
- Error responses return `{ "error": "message" }`

## Auth Challenge

`POST /v1/auth/challenge`

Request:

```json
{
  "wallet": "0x1111111111111111111111111111111111111111"
}
```

Response:

```json
{
  "wallet": "0x1111111111111111111111111111111111111111",
  "nonce": "7b3f...",
  "message": "Sign this DaoKin challenge nonce: 7b3f...",
  "created_at": "2026-04-30T10:00:00Z",
  "expires_at": "2026-04-30T10:05:00Z"
}
```

Failure cases:

- `400`: invalid request body or invalid wallet input
- `500`: challenge creation failure

## Auth Verify

`POST /v1/auth/verify`

Request:

```json
{
  "wallet": "0x1111111111111111111111111111111111111111",
  "nonce": "7b3f...",
  "signature": "0x..."
}
```

Response:

```json
{
  "wallet": "0x1111111111111111111111111111111111111111",
  "access_token": "8d9f...",
  "verified_at": "2026-04-30T10:01:00Z"
}
```

Failure cases:

- `400`: invalid request body or missing field
- `401`: challenge missing, expired, nonce mismatch, invalid signature, or wallet mismatch
- `500`: verification failure outside expected auth errors

## Artifact Create

`POST /v1/artifacts`

Auth required for `creator_wallet`.

Request:

```json
{
  "creator_wallet": "0x1111111111111111111111111111111111111111",
  "content_hash": "0xabc123...",
  "content_uri": "daokin://artifact/my-first-artifact"
}
```

Response:

```json
{
  "id": "art_000001",
  "creator_wallet": "0x1111111111111111111111111111111111111111",
  "content_hash": "0xabc123...",
  "content_uri": "daokin://artifact/my-first-artifact",
  "created_at": "2026-04-30T10:02:00Z"
}
```

Failure cases:

- `400`: invalid request body or missing artifact fields
- `401`: missing or invalid bearer token for `creator_wallet`
- `500`: persistence or command failure

## Artifact Get

`GET /v1/artifacts/{id}`

Response:

```json
{
  "id": "art_000001",
  "creator_wallet": "0x1111111111111111111111111111111111111111",
  "content_hash": "0xabc123...",
  "content_uri": "daokin://artifact/my-first-artifact",
  "created_at": "2026-04-30T10:02:00Z"
}
```

Failure cases:

- `404`: artifact not found

## Dao Join

`POST /v1/daos/{id}/join`

Auth required for `wallet`.

Request:

```json
{
  "wallet": "0x1111111111111111111111111111111111111111"
}
```

Response:

```json
{
  "dao_id": "dao-builders",
  "wallet": "0x1111111111111111111111111111111111111111",
  "joined_at": "2026-04-30T10:03:00Z",
  "left_at": null
}
```

Failure cases:

- `400`: invalid request body or empty DAO/wallet
- `401`: missing or invalid bearer token for `wallet`

## Dao Leave

`POST /v1/daos/{id}/leave`

Auth required for `wallet`.

Request:

```json
{
  "wallet": "0x1111111111111111111111111111111111111111",
  "reason": "MVP reversibility check"
}
```

Response:

```json
{
  "dao_id": "dao-builders",
  "wallet": "0x1111111111111111111111111111111111111111",
  "joined_at": "2026-04-30T10:03:00Z",
  "left_at": "2026-04-30T10:30:00Z",
  "leave_reason": "MVP reversibility check"
}
```

Failure cases:

- `400`: invalid request body or empty DAO/wallet
- `401`: missing or invalid bearer token for `wallet`
- `404`: active membership not found
- `409`: membership already left

## Permission Grant

`POST /v1/permissions`

Auth required for `granter_wallet`. The granter must be the artifact creator.

Request:

```json
{
  "artifact_id": "art_000001",
  "granter_wallet": "0x1111111111111111111111111111111111111111",
  "grantee_wallet": "0x2222222222222222222222222222222222222222",
  "scope": "view"
}
```

Response:

```json
{
  "id": "perm_000001",
  "artifact_id": "art_000001",
  "granter_wallet": "0x1111111111111111111111111111111111111111",
  "grantee_wallet": "0x2222222222222222222222222222222222222222",
  "scope": "view",
  "status": "active",
  "granted_at": "2026-04-30T10:04:00Z"
}
```

Failure cases:

- `400`: invalid request body or missing field
- `401`: missing or invalid bearer token for `granter_wallet`
- `403`: granter is not the artifact creator
- `404`: artifact not found
- `422`: invalid rule, such as granting to self

## Permission Revoke

`POST /v1/permissions/{id}/revoke`

Auth required for `granter_wallet`. The granter must match the original permission granter.

Request:

```json
{
  "granter_wallet": "0x1111111111111111111111111111111111111111",
  "reason": "Terms changed"
}
```

Response:

```json
{
  "id": "perm_000001",
  "artifact_id": "art_000001",
  "granter_wallet": "0x1111111111111111111111111111111111111111",
  "grantee_wallet": "0x2222222222222222222222222222222222222222",
  "scope": "view",
  "status": "revoked",
  "granted_at": "2026-04-30T10:04:00Z",
  "revoked_at": "2026-04-30T10:20:00Z",
  "revoke_reason": "Terms changed"
}
```

Failure cases:

- `400`: invalid request body or missing field
- `401`: missing or invalid bearer token for `granter_wallet`
- `403`: requester is not the permission granter
- `404`: permission not found
- `409`: permission already revoked

## Transfer Record

`POST /v1/transfers`

Auth required for `payer_wallet`. If the payer is not the artifact creator, the payer must have an active `view` permission for the artifact.

Request:

```json
{
  "artifact_id": "art_000001",
  "payer_wallet": "0x2222222222222222222222222222222222222222",
  "token": "USDC",
  "amount_atomic": "1000000",
  "tx_hash": "0xabc..."
}
```

Response:

```json
{
  "id": "tx_000001",
  "artifact_id": "art_000001",
  "payer_wallet": "0x2222222222222222222222222222222222222222",
  "recipient_wallet": "0x1111111111111111111111111111111111111111",
  "token": "USDC",
  "amount_atomic": "1000000",
  "tx_hash": "0xabc...",
  "created_at": "2026-04-30T10:05:00Z"
}
```

Failure cases:

- `400`: invalid request body or amount
- `401`: missing or invalid bearer token for `payer_wallet`
- `403`: non-owner payer lacks active `view` permission
- `404`: artifact not found

## Artifact Attribution

`GET /v1/artifacts/{id}/attribution`

Response:

```json
{
  "artifact_id": "art_000001",
  "permissions": [
    {
      "id": "perm_000001",
      "artifact_id": "art_000001",
      "granter_wallet": "0x1111111111111111111111111111111111111111",
      "grantee_wallet": "0x2222222222222222222222222222222222222222",
      "scope": "view",
      "status": "active",
      "granted_at": "2026-04-30T10:04:00Z"
    }
  ],
  "transfers": [
    {
      "id": "tx_000001",
      "artifact_id": "art_000001",
      "payer_wallet": "0x2222222222222222222222222222222222222222",
      "recipient_wallet": "0x1111111111111111111111111111111111111111",
      "token": "USDC",
      "amount_atomic": "1000000",
      "tx_hash": "0xabc...",
      "created_at": "2026-04-30T10:05:00Z"
    }
  ],
  "total_records": 2
}
```

Failure cases:

- `400`: empty artifact id
- `500`: repository query failure

## User Export

`GET /v1/users/{wallet}/export`

Auth required for `{wallet}`.

Response:

```json
{
  "wallet": "0x1111111111111111111111111111111111111111",
  "artifacts": [],
  "memberships": [],
  "permissions": [],
  "transfers": []
}
```

The export includes records created by, owned by, or associated with the wallet.

Failure cases:

- `400`: empty wallet
- `401`: missing or invalid bearer token for `{wallet}`

## Acceptance Coverage

- AC-1.1 provenance fields: `POST /v1/artifacts` and `GET /v1/artifacts/{id}`
- AC-1.2 data export path: `GET /v1/users/{wallet}/export`
- AC-2.1 consent record: `POST /v1/permissions` and `POST /v1/permissions/{id}/revoke`
- AC-2.2 value attribution: `POST /v1/transfers` and `GET /v1/artifacts/{id}/attribution`
- AC-3.1 join/leave reversibility: `POST /v1/daos/{id}/join` and `POST /v1/daos/{id}/leave`
- AC-3.2 no forced grouping: auth and artifact creation do not create DAO memberships
- AC-4.2 audit baseline: request logs include request id and authenticated actor wallet on protected write paths
