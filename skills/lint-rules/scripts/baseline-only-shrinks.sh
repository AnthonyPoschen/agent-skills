#!/usr/bin/env bash
# Lint baselines and allowlists may only shrink.
# Compare entry counts in each FILE to the same path at the merge base.
# A higher count fails unless an owner is recorded: a PR label whose name
# contains "owner", or an `owner:` annotation on every added entry.
#
# Usage:
#   baseline-only-shrinks.sh [--base REF] [--labels "a,b"] FILE [FILE...]
#
# --labels defaults to $PR_LABELS (comma-separated). --base defaults to the
# merge-base of HEAD and origin/master (then master).
set -euo pipefail

usage() {
  sed -n '2,12p' "$0" | sed 's/^# \{0,1\}//'
}

base_ref=""
labels="${PR_LABELS:-}"
files=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    --base)
      base_ref="${2:-}"
      shift 2
      ;;
    --labels)
      labels="${2:-}"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    --)
      shift
      files+=("$@")
      break
      ;;
    -*)
      echo "unknown option: $1" >&2
      exit 2
      ;;
    *)
      files+=("$1")
      shift
      ;;
  esac
done

if [[ ${#files[@]} -eq 0 ]]; then
  echo "pass at least one baseline or allowlist file" >&2
  exit 2
fi

if [[ -z "$base_ref" ]]; then
  if git rev-parse --verify --quiet origin/master >/dev/null; then
    base_ref=$(git merge-base HEAD origin/master)
  elif git rev-parse --verify --quiet master >/dev/null; then
    base_ref=$(git merge-base HEAD master)
  else
    echo "pass --base REF (no origin/master or master to merge against)" >&2
    exit 2
  fi
fi

has_owner_label=0
IFS=',' read -r -a label_list <<<"$labels"
for label in "${label_list[@]}"; do
  trimmed=$(printf '%s' "$label" | tr '[:upper:]' '[:lower:]' | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
  if [[ -n "$trimmed" && "$trimmed" == *owner* ]]; then
    has_owner_label=1
  fi
done

count_entries() {
  awk '
    /^[[:space:]]*$/ { next }
    /^[[:space:]]*#/ { next }
    /^[[:space:]]*\/\// { next }
    { n++ }
    END { print n + 0 }
  ' "$1"
}

failed=0
for file in "${files[@]}"; do
  if [[ ! -f "$file" ]]; then
    echo "missing file: $file" >&2
    exit 2
  fi
  head_count=$(count_entries "$file")
  base_file=$(mktemp)
  if git cat-file -e "${base_ref}:${file}" 2>/dev/null; then
    git show "${base_ref}:${file}" >"$base_file"
  else
    : >"$base_file"
  fi
  base_count=$(count_entries "$base_file")
  rm -f "$base_file"

  if [[ "$head_count" -le "$base_count" ]]; then
    echo "$file: $base_count -> $head_count (ok)"
    continue
  fi

  if [[ "$has_owner_label" -eq 1 ]]; then
    echo "$file: $base_count -> $head_count (grew; owner label present)"
    continue
  fi

  unowned=0
  added=0
  while IFS= read -r line; do
    body=${line#+}
    if printf '%s\n' "$body" | awk '
      /^[[:space:]]*$/ { exit 1 }
      /^[[:space:]]*#/ { exit 1 }
      /^[[:space:]]*\/\// { exit 1 }
      { exit 0 }
    '; then
      added=$((added + 1))
      if ! printf '%s\n' "$body" | grep -q 'owner:'; then
        unowned=$((unowned + 1))
      fi
    fi
  done < <(git diff --no-color -U0 "$base_ref" -- "$file" | grep -E '^\+' | grep -vE '^\+\+\+' || true)

  if [[ "$added" -gt 0 && "$unowned" -eq 0 ]]; then
    echo "$file: $base_count -> $head_count (grew; owner: on each added entry)"
    continue
  fi

  echo "$file: $base_count -> $head_count" >&2
  echo "Growing the baseline hides new violations. Fix the code instead." >&2
  echo "A lint baseline or allowlist may only shrink unless an owner label is present (a PR label whose name contains \"owner\", or an owner: annotation on every added entry)." >&2
  failed=1
done

exit "$failed"
