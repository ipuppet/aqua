#!/usr/bin/env bash
set -euo pipefail

# Enforce no direct cross-service imports.
# Rule:
# - Files under services/<svc>/... may import services/<svc>/... (self only)
# - Importing services/<other>/... is forbidden; use protobuf contracts under gen/go instead.
#
# Optional scope: AQUA_CHECK_TARGET=<service>

TARGET="${AQUA_CHECK_TARGET:-}"
SEARCH_ROOT="services"
if [[ -n "$TARGET" ]]; then
  SEARCH_ROOT="services/${TARGET}"
  if [[ ! -d "$SEARCH_ROOT" ]]; then
    echo "service import boundary check skipped: ${SEARCH_ROOT} not found"
    exit 0
  fi
fi

violations=""

while IFS= read -r line; do
  file_path="${line%%:*}"
  rest="${line#*:}"
  rest="${rest#*:}"

  file_service="$(echo "$file_path" | sed -E 's#^services/([^/]+)/.*#\1#')"
  import_service="$(echo "$rest" | sed -E 's#.*aqua/services/([^/"]+).*#\1#')"

  if [[ -z "$file_service" || -z "$import_service" ]]; then
    continue
  fi

  # mono is the modular-monolith host and may import every service package.
  if [[ "$file_service" == "mono" ]]; then
    continue
  fi

  if [[ "$file_service" != "$import_service" ]]; then
    violations+="$line"$'\n'
  fi
done < <(rg -n 'aqua/services/[^/"]+' "$SEARCH_ROOT" --glob '*.go' || true)

if [[ -n "$violations" ]]; then
  echo "service import boundary violations detected:"
  echo "$violations"
  exit 1
fi

echo "service import boundary check passed"
