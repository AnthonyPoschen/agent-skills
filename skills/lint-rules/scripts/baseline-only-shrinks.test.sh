#!/usr/bin/env bash
# Cases for baseline-only-shrinks.sh: no growth, missing label, missing owner
# lines, both present, and removal or reorder.
set -euo pipefail

script=$(cd "$(dirname "$0")" && pwd)/baseline-only-shrinks.sh
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
cd "$tmp"
git init -q
git config user.email test@example.com
git config user.name test

printf 'keep\nalpha\nbeta\n' > allow.txt
git add allow.txt
git commit -qm base
base=$(git rev-parse HEAD)

pass() {
  if ! "$script" --base "$base" "$@"; then
    echo "FAIL (expected pass): $*" >&2
    exit 1
  fi
}

fail() {
  if "$script" --base "$base" "$@" >/dev/null; then
    echo "FAIL (expected fail): $*" >&2
    exit 1
  fi
}

# no growth
pass allow.txt

# growth without the label
printf 'keep\nalpha\nbeta\nnew\n' > allow.txt
fail --labels "" allow.txt
fail --labels "baseline-owner" allow.txt

# label without an owner line directly before the added entry
fail --labels "baseline-growth-approved" allow.txt
printf 'keep\nalpha\nbeta\n# owner: ada — new path\n\nnew\n' > allow.txt
fail --labels "baseline-growth-approved" allow.txt
printf 'keep\nalpha\nbeta\n# owner: ada - hyphen is not an em dash\nnew\n' > allow.txt
fail --labels "baseline-growth-approved" allow.txt

# owner line without the exact label
printf 'keep\nalpha\nbeta\n# owner: ada — new path\nnew\n' > allow.txt
fail --labels "baseline-owner,owner" allow.txt

# both
pass --labels "other, baseline-growth-approved" allow.txt

# removal and reorder
printf 'beta\nkeep\n' > allow.txt
pass --labels "" allow.txt

echo "baseline-only-shrinks: ok"
