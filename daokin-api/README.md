# DaoKin API (Chi)

Minimal Go + Chi API skeleton for DaoKin.

## Quick start

```bash
cd daokin-api

go run ./cmd/api
```

## Endpoints

- `GET /healthz`
- `GET /readyz`
- `GET /v1/ping`

## Config

Environment variables (defaults shown):

- `DAOKIN_ENV=dev`
- `DAOKIN_ADDR=:8080`
- `DAOKIN_READ_TIMEOUT=10s`
- `DAOKIN_WRITE_TIMEOUT=10s`
- `DAOKIN_IDLE_TIMEOUT=60s`
- `DAOKIN_SHUTDOWN_TIMEOUT=10s`
