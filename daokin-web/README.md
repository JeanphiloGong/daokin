# DaoKin Web

`daokin-web` is the SvelteKit frontend for the DaoKin MVP flow. It owns the user-facing prototype for wallet auth, artifact creation, Dao joining, payment/tip interaction, and attribution viewing.

The shared API contract is maintained in [`../docs/architecture/api-contract-v1.md`](../docs/architecture/api-contract-v1.md). This README is the frontend module entrypoint.

## Boundary

- Owns routes, components, client state, frontend API client code, UI tests, and e2e tests.
- Currently defaults to mock API mode for the MVP demo page.
- Does not own backend authorization behavior, contract events, or project-wide governance policy.

## Key Paths

- `src/routes/`: SvelteKit page and layout surfaces.
- `src/lib/features/mvp/`: MVP feature types, stores, API clients, and components.
- `e2e/`: Playwright browser tests.
- `static/`: static assets served by the web app.

## Quick Start

```bash
npm install
npm run dev
```

The MVP page defaults to mock API mode. To run against `daokin-api`, create `daokin-web/.env.local` from `.env.example` and set:

```env
PUBLIC_DAOKIN_API_MODE=http
PUBLIC_DAOKIN_API_BASE_URL=/api
```

If the browser calls `daokin-api` on another origin directly, that API origin must allow the web dev origin through CORS.

## Build And Check

```bash
npm run build
npm run check
npm run lint
```

## Tests

```bash
npm run test:unit -- --run
npm run test:e2e
```

Playwright browser binaries are required for browser and e2e tests:

```bash
npx playwright install
```

## Related Docs

- System overview: [`../docs/system-overview.md`](../docs/system-overview.md)
- API contract: [`../docs/architecture/api-contract-v1.md`](../docs/architecture/api-contract-v1.md)
- MVP roadmap: [`../docs/roadmap/mvp.md`](../docs/roadmap/mvp.md)
