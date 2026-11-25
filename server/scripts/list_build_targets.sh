#!/usr/bin/env bash
# List buildable web server targets (services with a cmd main package).
# mono is always emitted last when present.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
SERVICES_DIR="${ROOT}/services"

names=()
has_mono=0

for cmd_dir in "${SERVICES_DIR}"/*/cmd; do
	[[ -d "${cmd_dir}" ]] || continue
	svc="$(basename "$(dirname "${cmd_dir}")")"
	case "${svc}" in
		_*|.*) continue ;;
	esac

	has_main=0
	shopt -s nullglob
	for f in "${cmd_dir}"/*.go; do
		if grep -q '^package main[[:space:]]*$' "${f}"; then
			has_main=1
			break
		fi
	done
	shopt -u nullglob
	[[ "${has_main}" -eq 1 ]] || continue

	if [[ "${svc}" == "mono" ]]; then
		has_mono=1
	else
		names+=("${svc}")
	fi
done

if [[ ${#names[@]} -gt 0 ]]; then
	# Stable order for CI / foreach builds.
	IFS=$'\n' names_sorted=($(printf '%s\n' "${names[@]}" | LC_ALL=C sort -u))
	unset IFS
	printf '%s' "${names_sorted[*]}"
	[[ "${has_mono}" -eq 1 ]] && printf ' mono'
	printf '\n'
elif [[ "${has_mono}" -eq 1 ]]; then
	echo mono
fi
