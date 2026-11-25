# Deployment example

This directory is a **single-host Docker Compose example** for Aqua. It shows one
reasonable way to run the public application code without committing any real
secrets or domain-specific configuration.

It is intentionally small:

- one PostgreSQL, Redis, etcd, and Meilisearch container;
- one `api` container running the `mono` build (admin, blog, passport, user);
- one nginx container serving the three frontend builds and proxying the API;
- one shared volume for uploaded images.

TLS is terminated by an external reverse proxy (Caddy, Traefik, Cloudflare,
another nginx, etc.). This example only listens on HTTP.

## Layout

```text
deploy/
├── compose.yaml
├── .env.example
├── config/
│   ├── app.example.yaml
│   └── admin.example.yaml
├── frontend/Dockerfile
├── nginx/
│   ├── proxy.inc
│   └── templates/default.conf.template
├── postgres/initdb.d/00-app.sh
└── secrets/                 # created locally, ignored by git
```

## Before you start

You need:

- Docker Engine with Compose v2;
- a DNS name such as `example.com`;
- four subdomains pointing at the host: `admin.example.com`,
  `blog.example.com`, `passport.example.com`, and `api.example.com`;
- a TLS reverse proxy in front of `HTTP_PORT`.

The frontend production build uses `https://` for API URLs, so this example
assumes a TLS proxy is in front of the nginx container. The nginx container
itself only listens on HTTP.

## 1. Configure the domain

```bash
cd deploy
cp .env.example .env
```

Edit `.env`:

```dotenv
AQUA_DOMAIN=example.com
HTTP_PORT=8080
```

If you use a different domain, also replace `api.example.com` in
`config/admin.example.yaml` and `no-reply@example.com` in
`config/app.example.yaml`.

## 2. Create secrets

Create the ignored secret files. Use different random values for each one:

```bash
mkdir -p secrets
openssl rand -base64 32 > secrets/postgres_root.txt
openssl rand -base64 32 > secrets/postgres_app.txt
openssl rand -base64 32 > secrets/redis.txt
openssl rand -base64 32 > secrets/meili.txt

openssl genpkey -algorithm RSA \
  -pkeyopt rsa_keygen_bits:2048 \
  -out secrets/web_jwt_private.pem
openssl rsa -pubout \
  -in secrets/web_jwt_private.pem \
  -out secrets/web_jwt_public.pem
```

`postgres_app.txt` is the password for the `web` database user. It must match
`database.password` in `config/app.yaml`. The Redis and Meilisearch files must
match `redis.password` and `meilisearch.api_key` respectively.

The JWT files are mounted at `/run/secrets/web_jwt_private` and
`/run/secrets/web_jwt_public`, matching `config/app.example.yaml`.

## 3. Create the runtime configuration

```bash
cp config/app.example.yaml config/app.yaml
cp config/admin.example.yaml config/admin.yaml
```

Replace every `REPLACE_*` value:

- `database.password` must equal the contents of `secrets/postgres_app.txt`;
- `redis.password` must equal the contents of `secrets/redis.txt`;
- `meilisearch.api_key` must equal the contents of `secrets/meili.txt`;
- `internal_service.token` and `file.key` should be independent random values.

`config/app.yaml` and `config/admin.yaml` are ignored by git and mounted
read-only into the API container. Do not commit them.

## 4. Start the stack

```bash
docker compose up -d --build
```

Check status and logs:

```bash
docker compose ps
docker compose logs -f api
```

Open:

- `https://admin.<AQUA_DOMAIN>` — admin console;
- `https://blog.<AQUA_DOMAIN>` — public site;
- `https://passport.<AQUA_DOMAIN>` — login and registration;
- `https://api.<AQUA_DOMAIN>/ping` — API health check.

Put your TLS reverse proxy in front of `HTTP_PORT` and forward those hostnames
to the same port. The application sets `VITE_AUTH_COOKIE_DOMAIN` to
`.${AQUA_DOMAIN}` at frontend build time so the passport cookie is shared by the
subdomains.

## What happens on first boot

PostgreSQL runs `postgres/initdb.d/00-app.sh`, which:

1. creates the `web` role using the password from `secrets/postgres_app.txt`;
2. creates the `web` database;
3. loads `../postgres/structure/web-structure.sql`.

The schema creates the tables but does **not** create a default administrator.
Register a user through the passport frontend, then grant the appropriate RBAC
roles/permissions using your own bootstrap procedure. Do not add default
credentials to the repository. Configure the `email` block and Turnstile first
if you want self-service registration and password reset to work.

## Updating

```bash
git pull
cd deploy
docker compose up -d --build
```

The example uses named volumes, so database and upload data survive rebuilds.
To remove everything, including data:

```bash
docker compose down -v
```

## Production notes

- This is an example, not a hardened production deployment.
- Terminate TLS outside this Compose stack and redirect HTTP to HTTPS there.
- Restrict the API CORS policy to your real frontend origins. The example adds
  CORS headers for the subdomain layout used by the bundled nginx.
- Replace the image tags with versions you have tested.
- Back up `postgres_data` and `storage_data`.
- If you run split services instead of `mono`, change the `SERVICE` build arg and
  the nginx upstreams accordingly.
