#!/usr/bin/env bash
set -e

if [ -z "${SERVICE_BIN:-}" ]; then
  echo "SERVICE_BIN is not set" >&2
  exit 1
fi

if [ -w "$SERVICE_BIN" ]; then
  chmod +x "$SERVICE_BIN"
fi
exec "$SERVICE_BIN"
