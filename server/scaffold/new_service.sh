#!/usr/bin/env bash
# Scaffold a new aqua web microservice from scaffold/templates/service.
# Usage: ./scaffold/new_service.sh <name> [--no-mono] [--no-gen] [--force]

set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
TEMPLATE_ROOT="${ROOT}/scaffold/templates/service"
RESERVED=(root admin blog passport user mono)

NAME=""
WITH_MONO=1
WITH_GEN=1
FORCE=0

usage() {
	cat <<'EOF'
Usage: ./scaffold/new_service.sh <name> [options]

Create a minimal grpc + grpc-gateway microservice (SayHello starter).
Only touches Go sources (proto, services/, mono, codegen) — not Docker/Compose.

Options:
  --no-mono      Do not register the service in services/mono/cmd
  --no-gen       Skip buf generate / wire regenerate
  --force        Overwrite an existing services/<name> tree
  -h, --help     Show this help

Makefile:
  make new-service name=<service>
EOF
}

while [[ $# -gt 0 ]]; do
	case "$1" in
		--no-mono) WITH_MONO=0; shift ;;
		--no-gen) WITH_GEN=0; shift ;;
		--force) FORCE=1; shift ;;
		-h|--help) usage; exit 0 ;;
		-*)
			echo "unknown option: $1" >&2
			usage >&2
			exit 1
			;;
		*)
			if [[ -n "${NAME}" ]]; then
				echo "unexpected argument: $1" >&2
				usage >&2
				exit 1
			fi
			NAME="$1"
			shift
			;;
	esac
done

if [[ -z "${NAME}" ]]; then
	usage >&2
	exit 1
fi

if [[ ! "${NAME}" =~ ^[a-z][a-z0-9]*$ ]]; then
	echo "invalid service name '${NAME}': use lowercase letters/digits, starting with a letter" >&2
	exit 1
fi

for reserved in "${RESERVED[@]}"; do
	if [[ "${NAME}" == "${reserved}" ]]; then
		echo "service name '${NAME}' is reserved" >&2
		exit 1
	fi
done

SERVICE_DIR="${ROOT}/services/${NAME}"
PROTO_DIR="${ROOT}/proto/aqua/${NAME}/v1"
PROTO_FILE="${PROTO_DIR}/${NAME}.proto"

if [[ -e "${SERVICE_DIR}" || -e "${PROTO_FILE}" ]]; then
	if [[ "${FORCE}" -ne 1 ]]; then
		echo "service '${NAME}' already exists (use --force to overwrite)" >&2
		exit 1
	fi
	rm -rf "${SERVICE_DIR}"
	rm -rf "${PROTO_DIR}"
fi

if [[ ! -d "${TEMPLATE_ROOT}" ]]; then
	echo "scaffold template missing: ${TEMPLATE_ROOT}" >&2
	exit 1
fi

PASCAL="$(python3 -c 'import sys; n=sys.argv[1]; print(n[:1].upper()+n[1:])' "${NAME}")"

echo "-> scaffolding service ${NAME} (${PASCAL})"

python3 - "${TEMPLATE_ROOT}" "${ROOT}" "${NAME}" "${PASCAL}" <<'PY'
import sys
from pathlib import Path

template_root = Path(sys.argv[1])
repo_root = Path(sys.argv[2])
name = sys.argv[3]
pascal = sys.argv[4]

replacements = {
    "__SERVICE__": name,
    "__Service__": pascal,
}

for src in template_root.rglob("*"):
    if not src.is_file():
        continue
    rel = src.relative_to(template_root)
    rel_str = str(rel)
    for old, new in replacements.items():
        rel_str = rel_str.replace(old, new)
    # Proto lives under proto/aqua/<service>/v1/, not proto/
    if rel_str.startswith("proto/"):
        # template: proto/__SERVICE__.proto -> proto/aqua/<name>/v1/<name>.proto
        dst = repo_root / "proto" / "aqua" / name / "v1" / f"{name}.proto"
    else:
        dst = repo_root / rel_str
    text = src.read_text(encoding="utf-8")
    for old, new in replacements.items():
        text = text.replace(old, new)
    dst.parent.mkdir(parents=True, exist_ok=True)
    dst.write_text(text, encoding="utf-8")
    print(f"   wrote {dst.relative_to(repo_root)}")
PY

# ---------------------------------------------------------------------------
# mono host registration
# ---------------------------------------------------------------------------
if [[ "${WITH_MONO}" -eq 1 ]]; then
	MONO_MAIN="${ROOT}/services/mono/cmd/main.go"
	python3 - "${MONO_MAIN}" "${NAME}" <<'PY'
import re
import sys
from pathlib import Path

path = Path(sys.argv[1])
name = sys.argv[2]
text = path.read_text(encoding="utf-8")

import_line = f'\t"aqua/services/{name}"'
service_line = f'\t\t{{Name: {name}.ServiceName, Factory: {name}.Assemble}},'

if import_line in text and f"{name}.ServiceName" in text:
    print(f"   mono already registers {name}")
    raise SystemExit(0)

# Insert import (alphabetically among aqua/services/*).
import_block = re.search(r"import \(\n(.*?\n)\)", text, re.S)
if not import_block:
    raise SystemExit("mono main.go: import block not found")
imports = import_block.group(1).splitlines(keepends=True)
svc_imports = [ln for ln in imports if '"aqua/services/' in ln]
other = [ln for ln in imports if ln not in svc_imports]
svc_imports.append(import_line + "\n")
# Keep stable order: existing style is not strictly alpha (root, admin, ...).
# Append before closing of service imports — place after last aqua/services import.
svc_imports_sorted = sorted(set(svc_imports), key=lambda s: s.strip().strip('"').split("/")[-1])
new_imports = "".join(other + svc_imports_sorted)
text = text[: import_block.start(1)] + new_imports + text[import_block.end(1) :]

# Insert RunAll entry before closing of the slice.
run_all = re.search(r"bootstrap\.RunAll\(\[\]bootstrap\.Service\{\n(.*?\n)\t\}\)", text, re.S)
if not run_all:
    raise SystemExit("mono main.go: RunAll block not found")
entries = run_all.group(1)
if f"{name}.ServiceName" not in entries:
    if not entries.endswith("\n"):
        entries += "\n"
    entries = entries + service_line + "\n"
    text = text[: run_all.start(1)] + entries + text[run_all.end(1) :]

path.write_text(text, encoding="utf-8")
print(f"   updated services/mono/cmd/main.go (+{name})")
PY
fi

# ---------------------------------------------------------------------------
# Codegen
# ---------------------------------------------------------------------------
if [[ "${WITH_GEN}" -eq 1 ]]; then
	echo "-> buf generate"
	(
		cd "${ROOT}"
		buf generate
	)
	echo "-> wire generate (${NAME})"
	(
		cd "${ROOT}/services/${NAME}"
		go run -mod=mod github.com/google/wire/cmd/wire
	)
fi

cat <<EOF

✨ service ${NAME} scaffolded

Next:
  make build target=${NAME}
  make run target=${NAME}
  curl -s 'http://127.0.0.1:8080/api/${NAME}/hello?name=aqua'

Optional:
  - Add docs/services/${NAME}.md for service-specific rules
  - Expand logic/dao beyond the SayHello starter

EOF
