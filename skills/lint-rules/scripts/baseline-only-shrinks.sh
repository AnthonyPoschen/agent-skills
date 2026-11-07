#!/usr/bin/env bash
# Lint baselines and allowlists may only shrink.
# Compare entry counts in each FILE to the same path at the merge base.
# A higher count fails unless both are true: the pull request has the exact
# label baseline-growth-approved, and the line directly before each added
# entry is `# owner: <who> — <why>` (em dash; who and why non-empty).
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

has_growth_label=0
IFS=',' read -r -a label_list <<<"$labels"
for label in "${label_list[@]}"; do
  trimmed=$(printf '%s' "$label" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
  if [[ "$trimmed" == "baseline-growth-approved" ]]; then
    has_growth_label=1
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

# Print "added unowned" for entries in the head file whose text is not
# matched by an entry at the base. An entry is owned only when the line
# immediately above it is `# owner: <who> — <why>`.
added_unowned() {
  awk '
    function is_entry(line) {
      return line !~ /^[[:space:]]*$/ && line !~ /^[[:space:]]*#/ && line !~ /^[[:space:]]*\/\//
    }
    function is_owner(line) {
      sub(/\r$/, "", line)
      gsub(/^[[:space:]]+|[[:space:]]+$/, "", line)
      return line ~ /^#[[:space:]]*owner:[[:space:]]+[^[:space:]].*[[:space:]]—[[:space:]]+[^[:space:]].*$/
    }
    NR == FNR {
      if (is_entry($0)) base[$0]++
      next
    }
    {
      if (is_entry($0)) {
        if (base[$0] > 0) {
          base[$0]--
        } else {
          added++
          if (!is_owner(prev)) unowned++
        }
      }
      prev = $0
    }
    END { printf "%d %d\n", added + 0, unowned + 0 }
  ' "$1" "$2"
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

  if [[ "$head_count" -le "$base_count" ]]; then
    echo "$file: $base_count -> $head_count (ok)"
    rm -f "$base_file"
    continue
  fi

  read -r added unowned < <(added_unowned "$base_file" "$file")
  rm -f "$base_file"

  if [[ "$has_growth_label" -eq 1 && "$added" -gt 0 && "$unowned" -eq 0 ]]; then
    echo "$file: $base_count -> $head_count (grew; baseline-growth-approved and owner lines)"
    continue
  fi

  echo "$file: $base_count -> $head_count" >&2
  echo "Growing the baseline hides new violations. Fix the code instead." >&2
  if [[ "$has_growth_label" -ne 1 ]]; then
    echo "Missing the exact PR label baseline-growth-approved." >&2
  fi
  if [[ "$unowned" -gt 0 || "$added" -eq 0 ]]; then
    echo "Each added entry needs a \`# owner: <who> — <why>\` line directly before it (em dash; who and why non-empty)." >&2
  fi
  failed=1
done

exit "$failed"
