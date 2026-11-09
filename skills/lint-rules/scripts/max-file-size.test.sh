#!/usr/bin/env bash
# Cases for max-file-size.sh: under the ceiling, over it, baseline recorded
# size, growth of a recorded size, and baseline entry growth via
# baseline-only-shrinks.sh.
set -euo pipefail

script=$(cd "$(dirname "$0")" && pwd)/max-file-size.sh
tmp=$(mktemp -d)
trap 'rm -rf "$tmp"' EXIT
cd "$tmp"
git init -q
git config user.email test@example.com
git config user.name test

why='files this large are re-read in pieces by agents dozens of times per run; split by responsibility'

head -c 100 </dev/zero > small.txt
head -c 200 </dev/zero > big.txt
printf 'big.txt 200\n' > size.txt
git add small.txt big.txt size.txt
git commit -qm base
base=$(git rev-parse HEAD)

run() {
  set +e
  out=$("$script" "$@" 2>&1)
  status=$?
  set -e
}

pass() {
  run "$@"
  if [[ "$status" -ne 0 ]]; then
    echo "FAIL (expected pass, got $status): $*" >&2
    printf '%s\n' "$out" >&2
    exit 1
  fi
}

fail_size() {
  run "$@"
  if [[ "$status" -ne 1 ]]; then
    echo "FAIL (expected exit 1, got $status): $*" >&2
    printf '%s\n' "$out" >&2
    exit 1
  fi
  local count
  count=$(grep -cF "$why" <<<"$out" || true)
  if [[ "$count" -ne 1 ]]; then
    echo "FAIL (why line count $count, want 1): $*" >&2
    printf '%s\n' "$out" >&2
    exit 1
  fi
}

usage() {
  run "$@"
  if [[ "$status" -ne 2 ]]; then
    echo "FAIL (expected exit 2, got $status): $*" >&2
    printf '%s\n' "$out" >&2
    exit 1
  fi
}

# under and exactly at the ceiling
pass --max-bytes 200 --base "$base" .
pass --max-bytes 100 --base "$base" small.txt

# over the ceiling, not listed
fail_size --max-bytes 100 --base "$base" .

# listed, still at its recorded size, over the ceiling
pass --max-bytes 100 --base "$base" --baseline size.txt .

# untracked file over the ceiling is ignored
head -c 400 </dev/zero > untracked.txt
pass --max-bytes 100 --base "$base" --baseline size.txt .

# baselined file grows past its recorded size
head -c 201 </dev/zero > big.txt
fail_size --max-bytes 100 --base "$base" --baseline size.txt .
head -c 200 </dev/zero > big.txt

# raising the recorded size is growth, even with the approval label
printf 'big.txt 201\n' > size.txt
fail_size --max-bytes 100 --base "$base" --labels "baseline-growth-approved" --baseline size.txt .

# lowering the recorded size with the file still passes
head -c 100 </dev/zero > big.txt
printf 'big.txt 100\n' > size.txt
pass --max-bytes 100 --base "$base" --baseline size.txt .

# a new oversized entry needs the shrink rules (label and owner line)
git checkout -q -- big.txt size.txt
head -c 250 </dev/zero > extra.txt
git add extra.txt
printf 'big.txt 200\nextra.txt 250\n' > size.txt
run --max-bytes 100 --base "$base" --labels "" --baseline size.txt .
if [[ "$status" -ne 1 ]]; then
  echo "FAIL (new entry without label should exit 1, got $status)" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi
grep -qF 'Growing the baseline hides new violations' <<<"$out" || {
  echo "FAIL (expected baseline-only-shrinks message)" >&2
  printf '%s\n' "$out" >&2
  exit 1
}

printf 'big.txt 200\n# owner: ada — grandfather extra\nextra.txt 250\n' > size.txt
pass --max-bytes 100 --base "$base" --labels "baseline-growth-approved" --baseline size.txt .

# two files over the ceiling share one why line
git checkout -q -- big.txt
git rm -q -f -- extra.txt
head -c 250 </dev/zero > other.txt
git add other.txt
printf 'small.txt 100\n' > size.txt
fail_size --max-bytes 100 --base "$base" --baseline size.txt .

usage --baseline size.txt .
usage --max-bytes no .
usage --max-bytes 100 missing.txt

# Replacing or renaming a path is a new path key even when the entry count
# does not grow. Drop the old file so the failure is the new key, not an
# unlisted file that is still over the ceiling.
git reset -q --hard "$base"
git rm -q big.txt
head -c 1000 </dev/zero > new.txt
git add new.txt
printf 'new.txt 1000\n' > size.txt
run --max-bytes 100 --base "$base" --labels "" --baseline size.txt .
if [[ "$status" -ne 1 ]] || ! grep -qF 'new baseline path: new.txt' <<<"$out"; then
  echo "FAIL (replacing old.txt with new.txt should deny, got $status)" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi
printf '# owner: ada — grandfather new.txt\nnew.txt 1000\n' > size.txt
run --max-bytes 100 --base "$base" --labels "" --baseline size.txt .
if [[ "$status" -ne 1 ]] || ! grep -qF 'baseline-growth-approved' <<<"$out"; then
  echo "FAIL (owned replacement without the label should deny, got $status)" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi
run --max-bytes 100 --base "$base" --labels "baseline-growth-approved" --baseline size.txt .
if [[ "$status" -ne 0 ]]; then
  echo "FAIL (approved replacement should pass, got $status)" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi

git reset -q --hard "$base"
git mv big.txt renamed.txt
head -c 1000 </dev/zero > renamed.txt
git add renamed.txt
printf 'renamed.txt 1000\n' > size.txt
run --max-bytes 100 --base "$base" --labels "baseline-growth-approved" --baseline size.txt .
if [[ "$status" -ne 1 ]] || ! grep -qF 'new baseline path: renamed.txt' <<<"$out"; then
  echo "FAIL (rename and growth without an owner line should deny, got $status)" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi
printf '# owner: ada — renamed and grew\nrenamed.txt 1000\n' > size.txt
pass --max-bytes 100 --base "$base" --labels "baseline-growth-approved" --baseline size.txt .

# An empty committed baseline can be initialized only with label and owner.
git reset -q --hard "$base"
: > size.txt
git add size.txt
git commit -qm empty-baseline
empty=$(git rev-parse HEAD)
head -c 200 </dev/zero > big.txt
printf '# owner: ada — initial ceiling\nbig.txt 200\n' > size.txt
pass --max-bytes 100 --base "$empty" --labels "baseline-growth-approved" --baseline size.txt .
printf 'big.txt 200\n' > size.txt
run --max-bytes 100 --base "$empty" --labels "baseline-growth-approved" --baseline size.txt .
if [[ "$status" -ne 1 ]]; then
  echo "FAIL (empty baseline init without an owner line should deny, got $status)" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi
printf '# owner: ada — initial ceiling\nbig.txt 200\n' > size.txt
run --max-bytes 100 --base "$empty" --labels "" --baseline size.txt .
if [[ "$status" -ne 1 ]]; then
  echo "FAIL (empty baseline init without the label should deny, got $status)" >&2
  printf '%s\n' "$out" >&2
  exit 1
fi

# From a subdirectory, --baseline names the work-tree file. The merge-base
# object is that file from the repository root, not a root file of the same
# name. Here the root baseline is empty and the nested one has one entry.
git reset -q --hard "$base"
: > size.txt
mkdir -p sub
head -c 200 </dev/zero > sub/big.txt
printf 'sub/big.txt 200\n' > sub/size.txt
git add -A
git commit -qm nested
nested=$(git rev-parse HEAD)
(
  cd sub
  if ! "$script" --max-bytes 100 --base "$nested" --baseline size.txt .; then
    echo "FAIL (unchanged nested baseline should pass from a subdirectory)" >&2
    exit 1
  fi
  printf 'sub/big.txt 200\nsub/other.txt 400\n' > size.txt
  head -c 400 </dev/zero > other.txt
  git add other.txt
  set +e
  out=$("$script" --max-bytes 100 --base "$nested" --labels "" --baseline size.txt . 2>&1)
  status=$?
  set -e
  if [[ "$status" -ne 1 ]] || ! grep -qF 'new baseline path: sub/other.txt' <<<"$out"; then
    echo "FAIL (new path from a subdirectory should deny, got $status)" >&2
    printf '%s\n' "$out" >&2
    exit 1
  fi
)

echo "max-file-size: ok"
