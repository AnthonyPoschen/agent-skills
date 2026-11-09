#!/usr/bin/env bash
# Fail when a tracked file is over a byte ceiling.
# A baseline entry "<path> <bytes>" may list a file over the ceiling only up
# to that recorded size. A baselined file that grows past its recorded size
# fails. Recorded sizes may only shrink. Adding a baseline entry is checked
# with baseline-only-shrinks.sh (same label and owner-line rules).
#
# Usage:
#   max-file-size.sh --max-bytes N [--baseline FILE] [--base REF] [--labels "a,b"] <paths...>
#
# --labels defaults to $PR_LABELS. --base defaults to the merge-base of HEAD
# and origin/master (then master). Paths are files or directories. Only
# tracked files are checked.
set -euo pipefail

usage() {
  sed -n '2,13p' "$0" | sed 's/^# \{0,1\}//'
}

why='files this large are re-read in pieces by agents dozens of times per run; split by responsibility'

max_bytes=""
baseline=""
base_ref=""
labels_set=0
labels=""
paths=()

while [[ $# -gt 0 ]]; do
  case "$1" in
    --max-bytes)
      max_bytes="${2:-}"
      shift 2
      ;;
    --baseline)
      baseline="${2:-}"
      shift 2
      ;;
    --base)
      base_ref="${2:-}"
      shift 2
      ;;
    --labels)
      labels="${2:-}"
      labels_set=1
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    --)
      shift
      paths+=("$@")
      break
      ;;
    -*)
      echo "unknown option: $1" >&2
      exit 2
      ;;
    *)
      paths+=("$1")
      shift
      ;;
  esac
done

if [[ ! "$max_bytes" =~ ^[0-9]+$ ]]; then
  echo "pass --max-bytes N (a non-negative integer)" >&2
  exit 2
fi

if [[ ${#paths[@]} -eq 0 ]]; then
  echo "pass at least one path" >&2
  exit 2
fi

for path in "${paths[@]}"; do
  if [[ ! -e "$path" ]]; then
    echo "missing path: $path" >&2
    exit 2
  fi
done

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "run inside a git work tree (only tracked files are checked)" >&2
  exit 2
fi

script_dir=$(cd "$(dirname "$0")" && pwd)

resolve_base() {
  if [[ -n "$base_ref" ]]; then
    return
  fi
  if git rev-parse --verify --quiet origin/master >/dev/null; then
    base_ref=$(git merge-base HEAD origin/master)
  elif git rev-parse --verify --quiet master >/dev/null; then
    base_ref=$(git merge-base HEAD master)
  else
    echo "pass --base REF (no origin/master or master to merge against)" >&2
    exit 2
  fi
}

# Print path<TAB>bytes for each entry. Exit 2 on a malformed line.
parse_baseline() {
  awk '
    {
      sub(/\r$/, "")
      if ($0 ~ /^[[:space:]]*$/ || $0 ~ /^[[:space:]]*#/ || $0 ~ /^[[:space:]]*\/\//) next
      line = $0
      gsub(/^[[:space:]]+|[[:space:]]+$/, "", line)
      if (!match(line, /[[:space:]]+[0-9]+$/)) {
        print "bad baseline entry: " line > "/dev/stderr"
        exit 2
      }
      bytes = substr(line, RSTART)
      gsub(/^[[:space:]]+/, "", bytes)
      path = substr(line, 1, RSTART - 1)
      gsub(/[[:space:]]+$/, "", path)
      sub(/^\.\//, "", path)
      if (path == "") {
        print "bad baseline entry: " line > "/dev/stderr"
        exit 2
      }
      if (seen[path]++) {
        print "duplicate baseline entry: " path > "/dev/stderr"
        exit 2
      }
      print path "\t" bytes
    }
  ' "$1"
}

load_sizes() {
  local file="$1"
  local -n dest=$2
  local parsed path bytes
  parsed=$(parse_baseline "$file") || exit $?
  while IFS=$'\t' read -r path bytes; do
    [[ -z "${path:-}" ]] && continue
    dest["$path"]=$bytes
  done <<<"$parsed"
}

declare -A recorded=()
declare -A base_recorded=()

size_failed=0
shrink_failed=0

if [[ -n "$baseline" ]]; then
  if [[ ! -f "$baseline" ]]; then
    echo "missing file: $baseline" >&2
    exit 2
  fi
  load_sizes "$baseline" recorded
  resolve_base
  shrink_cmd=("$script_dir/baseline-only-shrinks.sh" --base "$base_ref")
  if [[ "$labels_set" -eq 1 ]]; then
    shrink_cmd+=(--labels "$labels")
  fi
  shrink_cmd+=("$baseline")
  set +e
  "${shrink_cmd[@]}"
  shrink_status=$?
  set -e
  if [[ "$shrink_status" -eq 2 ]]; then
    exit 2
  fi
  if [[ "$shrink_status" -ne 0 ]]; then
    shrink_failed=1
  fi

  base_file=$(mktemp)
  if git cat-file -e "${base_ref}:${baseline}" 2>/dev/null; then
    git show "${base_ref}:${baseline}" >"$base_file"
    load_sizes "$base_file" base_recorded
  fi
  rm -f "$base_file"

  for path in "${!recorded[@]}"; do
    if [[ -n "${base_recorded[$path]+x}" && "${recorded[$path]}" -gt "${base_recorded[$path]}" ]]; then
      echo "$path: recorded size ${base_recorded[$path]} -> ${recorded[$path]}" >&2
      size_failed=1
    fi
  done
fi

root=$(git rev-parse --show-toplevel)
while IFS= read -r -d '' file; do
  [[ -f "$root/$file" ]] || continue
  bytes=$(wc -c <"$root/$file" | tr -d '[:space:]')
  if [[ -n "${recorded[$file]+x}" ]]; then
    if [[ "$bytes" -gt "${recorded[$file]}" ]]; then
      echo "$file: $bytes bytes > recorded ${recorded[$file]}" >&2
      size_failed=1
    fi
  elif [[ "$bytes" -gt "$max_bytes" ]]; then
    echo "$file: $bytes bytes > max $max_bytes (not in baseline)" >&2
    size_failed=1
  fi
done < <(git ls-files -z --full-name -- "${paths[@]}")

if [[ "$size_failed" -eq 1 ]]; then
  echo "$why" >&2
fi

if [[ "$size_failed" -eq 1 || "$shrink_failed" -eq 1 ]]; then
  exit 1
fi
exit 0
