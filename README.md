# Aqua

Aqua is a self-hosted web platform: Go services for identity, users, content,
and administration, plus Vue 3 frontends for the public site and admin console.

> This repository contains the reusable application code. Deployment-specific
> configuration, secrets, and the private server toolbox are intentionally not
> included.

## Stack

- **Backend:** Go, Gin, gRPC, protobuf, Wire
- **Frontend:** Vue 3, TypeScript, Vite, npm workspaces
- **Infrastructure:** PostgreSQL, Redis, etcd, Meilisearch

## Repository layout

| Path | Purpose |
| --- | --- |
| `server/` | Go services (`admin`, `blog`, `passport`, `user`), shared packages, protobuf contracts, migrations |
| `frontend/` | npm workspace with `apps/admin`, `apps/blog`, `apps/passport`, and shared packages |
| `postgres/` | Database schema and an example bootstrap script |

## Quick start

### 1. Backend

Copy the example configuration and replace every placeholder:

```bash
cd server
cp config/app.example.yaml config/app.yaml
cp config/admin.example.yaml config/admin.yaml
```

Generate JWT keys for local development (see `server/config/README.md`), then
start the services you need:

```bash
go run ./services/mono/cmd -config-dir=./config
```

The full stack expects PostgreSQL, Redis, etcd, and Meilisearch to be reachable
at the hosts configured in `config/app.yaml`. Apply
`postgres/structure/web-structure.sql` and the example bootstrap first.

### 2. Frontend

```bash
cd frontend
npm install
cp .env.development.example .env.development
npm run dev:blog
# or: npm run dev:admin / npm run dev:passport
```

Build all apps with `npm run build:admin`, `npm run build:blog`, and
`npm run build:passport`.

## Configuration and secrets

- Never commit `*.pem`, database passwords, SMTP credentials, API tokens, or
  provider keys.
- `server/config/*.example.yaml` and `frontend/.env*.example` are templates.
- In production, mount secrets at `/run/secrets/` and keep real configuration
  outside the repository.
- The private deployment repository is expected to supply compose files,
  reverse-proxy configuration, certificates, and environment-specific values.

## Development checks

```bash
cd server && make check
cd frontend && npm run check
```

## License

MIT. See [LICENSE](./LICENSE).
