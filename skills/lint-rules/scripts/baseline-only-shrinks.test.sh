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

# zero-byte baseline at the merge base: approved initialization passes,
# and the same initialization without the label or owner line fails
: > allow.txt
git add allow.txt
git commit -qm empty
empty=$(git rev-parse HEAD)
printf '# owner: ada — first entry\nfirst\n' > allow.txt
if ! "$script" --base "$empty" --labels "baseline-growth-approved" allow.txt; then
  echo "FAIL (approved init of an empty baseline should pass)" >&2
  exit 1
fi
if "$script" --base "$empty" --labels "" allow.txt >/dev/null; then
  echo "FAIL (empty baseline init without the label should fail)" >&2
  exit 1
fi
printf 'first\n' > allow.txt
if "$script" --base "$empty" --labels "baseline-growth-approved" allow.txt >/dev/null; then
  echo "FAIL (empty baseline init without an owner line should fail)" >&2
  exit 1
fi

# baseline path missing at the merge base, including from a subdirectory
git rm -qf -- allow.txt
git commit -qm remove
missing=$(git rev-parse HEAD)
mkdir -p sub
printf '# owner: ada — first entry\nnested\n' > sub/allow.txt
git add sub/allow.txt
if ! "$script" --base "$missing" --labels "baseline-growth-approved" sub/allow.txt; then
  echo "FAIL (approved init of a missing baseline should pass)" >&2
  exit 1
fi
(
  cd sub
  if ! "$script" --base "$missing" --labels "baseline-growth-approved" allow.txt; then
    echo "FAIL (approved init from a subdirectory should see the repo path)" >&2
    exit 1
  fi
)
printf 'nested\n' > sub/allow.txt
if "$script" --base "$missing" --labels "baseline-growth-approved" sub/allow.txt >/dev/null; then
  echo "FAIL (missing baseline init without an owner line should fail)" >&2
  exit 1
fi

# unchanged file from a subdirectory still compares against that repo path
git add sub/allow.txt
git commit -qm nested
nested=$(git rev-parse HEAD)
(
  cd sub
  if ! "$script" --base "$nested" --labels "" allow.txt; then
    echo "FAIL (unchanged baseline from a subdirectory should pass)" >&2
    exit 1
  fi
  printf 'nested\nextra\n' > allow.txt
  if "$script" --base "$nested" --labels "" allow.txt >/dev/null; then
    echo "FAIL (growth from a subdirectory should fail)" >&2
    exit 1
  fi
)

echo "baseline-only-shrinks: ok"
