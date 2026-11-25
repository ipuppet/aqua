# PostgreSQL

- `structure/web-structure.sql` contains the Aqua schema.
- `initdb.d/0-init.sql.example` is a minimal bootstrap script for a fresh
  database. Copy it to `0-init.sql`, replace the password placeholder, and
  mount the directory into the PostgreSQL container's `docker-entrypoint-initdb.d`.
