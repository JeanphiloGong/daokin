# DaoKin API (Chi)

Minimal Go + Chi API skeleton for DaoKin.

## Layering

- `internal/domain`: domain models and invariants
- `internal/app/ports/in`: input ports (use-case contracts for interfaces)
- `internal/app/ports/out`: output ports (repository contracts)
- `internal/app/usecases`: application use-cases
- `internal/interfaces/http`: HTTP transport interface
- `internal/interfaces/grpc`: gRPC transport placeholder
- `internal/infra/persistence`: persistence adapters (current: memory)
- `internal/infra/messaging`: messaging adapter placeholder

## Quick start

```bash
cd daokin-api

go run ./cmd/server
```

## Endpoints

- `GET /healthz`
- `GET /readyz`
- `GET /v1/ping`
- `POST /v1/auth/challenge`
- `POST /v1/auth/verify`
- `POST /v1/artifacts`
- `GET /v1/artifacts/{id}`
- `GET /v1/artifacts/{id}/attribution`
- `POST /v1/daos/{id}/join`
- `POST /v1/daos/{id}/leave`
- `POST /v1/permissions`
- `POST /v1/permissions/{id}/revoke`
- `POST /v1/transfers`
- `GET /v1/users/{wallet}/export`

## Auth notes (MVP)

- `POST /v1/auth/challenge` issues a nonce message for wallet signing.
- `POST /v1/auth/verify` expects a `personal_sign` (EIP-191) signature of the returned message.
- Write endpoints currently require `Authorization: Bearer <access_token>` and wallet must match the verified session.
- Permission and transfer flows are recorded off-chain for MVP attribution trail queries.

## Config

Environment variables (defaults shown):

- `DAOKIN_ENV=dev`
- `DAOKIN_ADDR=:8080`
- `DAOKIN_READ_TIMEOUT=10s`
- `DAOKIN_WRITE_TIMEOUT=10s`
- `DAOKIN_IDLE_TIMEOUT=60s`
- `DAOKIN_SHUTDOWN_TIMEOUT=10s`
