#!/usr/bin/env bash
set -euo pipefail

# Guardrail: dao layer must not depend on HTTP transport-layer map helpers.
# This blocks `gin.H` leakage from dao/logic into handler responses.
#
# Optional scope: AQUA_CHECK_TARGET=<service>

TARGET="${AQUA_CHECK_TARGET:-}"
SEARCH_ROOT="services"
GLOB='**/dao/*.go'
if [[ -n "$TARGET" ]]; then
  SEARCH_ROOT="services/${TARGET}"
  if [[ ! -d "$SEARCH_ROOT" ]]; then
    echo "dao transport leak check skipped: ${SEARCH_ROOT} not found"
    exit 0
  fi
fi

violations="$(
  rg -n "gin\.H|\*gin\.Context|gin\.Context" "$SEARCH_ROOT" --glob "$GLOB" || true
)"

if [[ -n "$violations" ]]; then
  echo "dao transport leak check failed: forbidden gin transport types found in dao layer"
  echo "$violations"
  exit 1
fi

echo "dao transport leak check passed"
