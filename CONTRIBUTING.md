# Contributing

## Before opening a pull request

- Run `cd server && make check`.
- Run `cd frontend && npm install && npm run check`.
- Keep generated protobuf and OpenAPI files in sync.
- Do not add real credentials, domains, customer data, or deployment-specific
  configuration. Use the provided `*.example.*` files.

## Commit style

Use short, imperative commit messages, for example:

- `feat(blog): add tag filtering`
- `fix(passport): reject untrusted redirect hosts`
- `docs: clarify local setup`
