#!/bin/sh
set -eu

: "${APP_DB_PASSWORD_FILE:?APP_DB_PASSWORD_FILE is required}"
: "${POSTGRES_USER:?POSTGRES_USER is required}"

app_password=$(cat "$APP_DB_PASSWORD_FILE")

# Create the application role and database on first initialization.
psql -v ON_ERROR_STOP=1 \
  --username "$POSTGRES_USER" \
  --dbname postgres \
  -v app_password="$app_password" <<'SQL'
CREATE ROLE web LOGIN PASSWORD :'app_password';
CREATE DATABASE web OWNER web;
SQL

# Load the Aqua schema into the application database.
psql -v ON_ERROR_STOP=1 \
  --username "$POSTGRES_USER" \
  --dbname web \
  --file /schema/web-structure.sql
