#!/usr/bin/env bash
set -euo pipefail

# Optional scope: AQUA_CHECK_TARGET=<service>

ROOT="${ROOT:-.}"
TARGET="${AQUA_CHECK_TARGET:-}"

python3 - "$ROOT" "$TARGET" <<'PY'
from pathlib import Path
import re
import sys

ROOT = Path(sys.argv[1])
TARGET = sys.argv[2].strip()

if TARGET:
    wire_files = [ROOT / "services" / TARGET / "wire.go"]
    wire_files = [p for p in wire_files if p.is_file()]
    if not wire_files:
        print(f"wire build order check skipped: no services/{TARGET}/wire.go")
        sys.exit(0)
else:
    wire_files = sorted((ROOT / "services").glob("*/wire.go"))
    if not wire_files:
        print("wire build order check failed: no services/*/wire.go found")
        sys.exit(1)

ORDER_LABELS = {
    0: "bootstrap runtime",
    1: "bootstrap shared providers",
    2: "config/token providers",
    3: "infra providers",
    4: "service providers",
    5: "final assembler",
}


def parse_import_aliases(text: str) -> dict[str, str]:
    aliases: dict[str, str] = {}
    m = re.search(r"import\s*\((.*?)\)\s*", text, re.S)
    if not m:
        return aliases
    for raw in m.group(1).splitlines():
        line = raw.strip()
        if not line or line.startswith("//"):
            continue
        mm = re.match(r'(?:(\w+)\s+)?"([^"]+)"', line)
        if not mm:
            continue
        alias, path = mm.groups()
        if alias is None:
            alias = path.rsplit("/", 1)[-1]
        aliases[alias] = path
    return aliases


def extract_wire_build_args(text: str) -> list[str]:
    idx = text.find("wire.Build(")
    if idx < 0:
        return []
    i = idx + len("wire.Build(")
    depth = 1
    buf = []
    while i < len(text) and depth > 0:
        ch = text[i]
        if ch == "(":
            depth += 1
        elif ch == ")":
            depth -= 1
            if depth == 0:
                break
        buf.append(ch)
        i += 1
    inside = "".join(buf)
    args = []
    cur = []
    depth2 = 0
    for ch in inside:
        if ch == "(":
            depth2 += 1
        elif ch == ")":
            depth2 -= 1
        if ch == "," and depth2 == 0:
            arg = "".join(cur).strip()
            if arg:
                args.append(arg)
            cur = []
            continue
        cur.append(ch)
    tail = "".join(cur).strip()
    if tail:
        args.append(tail)
    return args


def classify_arg(arg: str, aliases: dict[str, str]) -> int:
    arg = arg.strip()
    if arg == "bootstrap.RuntimeBaseSet" or arg == "boot.RuntimeSet":
        return 0
    if arg == "provideModule":
        return 5
    if arg == "bootstrap.FileStoreSet":
        return 1

    if "." not in arg:
        return 4

    prefix = arg.split(".", 1)[0]
    import_path = aliases.get(prefix, "")

    if import_path == "github.com/ipuppet/aqua/server/pkg/infra/config" or import_path == "github.com/ipuppet/aqua/server/pkg/transport/token":
        return 2

    if import_path.startswith("github.com/ipuppet/aqua/server/services/"):
        return 4

    if import_path.startswith("github.com/ipuppet/aqua/server/pkg/"):
        return 3

    return 4


violations: list[str] = []

for path in wire_files:
    text = path.read_text(encoding="utf-8")
    aliases = parse_import_aliases(text)
    args = extract_wire_build_args(text)
    if not args:
        violations.append(f"{path}: missing wire.Build(...)")
        continue
    if args[0] != "bootstrap.RuntimeBaseSet":
        violations.append(f"{path}: first arg must be bootstrap.RuntimeBaseSet, got {args[0]}")
    if args[-1] != "provideModule":
        violations.append(f"{path}: last arg must be provideModule, got {args[-1]}")

    prev = -1
    for idx, arg in enumerate(args):
        cat = classify_arg(arg, aliases)
        if cat < prev:
            violations.append(
                f"{path}: arg#{idx+1} {arg} breaks order ({ORDER_LABELS.get(prev)} -> {ORDER_LABELS.get(cat)})"
            )
        prev = cat

if violations:
    print("wire build order check failed:")
    for v in violations:
        print(f"- {v}")
    sys.exit(1)

print("wire build order check passed")
PY
