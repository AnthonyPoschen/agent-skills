#!/usr/bin/env bash
# Fail when a tracked file is over a byte ceiling.
# A baseline entry "<path> <bytes>" may list a file over the ceiling only up
# to that recorded size. Paths in the baseline are relative to the repository
# root. A baselined file that grows past its recorded size fails. Recorded
# sizes may only shrink. A path that is not in the baseline at the merge base
# is a new path key, including when another path is removed so the entry count
# stays the same. A new path key needs the exact label baseline-growth-approved
# and a `# owner: <who> — <why>` line directly before it. Entry-count growth
# is also checked with baseline-only-shrinks.sh.
#
# Byte counts are the working-tree files for those tracked paths. This does
# not read the staged index blob or the blob at HEAD.
#
# Usage:
#   max-file-size.sh --max-bytes N [--baseline FILE] [--base REF] [--labels "a,b"] <paths...>
#
# --labels defaults to $PR_LABELS. --base defaults to the merge-base of HEAD
# and origin/master (then master). Paths are files or directories. Only
# tracked files are checked. --baseline is a work-tree path; the merge-base
# copy is that file's path from the repository root, including when this
# script is run from a subdirectory.
set -euo pipefail

usage() {
  sed -n '2,22p' "$0" | sed 's/^# \{0,1\}//'
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

# Work-tree path as given (including from a subdirectory) to the path Git
# stores from the repository root. git show REF:path is always root-relative.
repo_path() {
  local file="$1"
  local root abs
  root=$(git rev-parse --show-toplevel)
  abs=$(cd "$(dirname "$file")" && pwd)/$(basename "$file")
  case "$abs" in
    "$root"/*)
      printf '%s\n' "${abs#"$root"/}"
      ;;
    *)
      echo "path is outside the repository: $file" >&2
      exit 2
      ;;
  esac
}

# Print path<TAB>bytes<TAB>owned for each entry. owned is 1 when the line
# directly above the entry is `# owner: <who> — <why>`. Exit 2 on a bad line.
parse_baseline() {
  awk '
    function is_owner(line) {
      sub(/\r$/, "", line)
      gsub(/^[[:space:]]+|[[:space:]]+$/, "", line)
      return line ~ /^#[[:space:]]*owner:[[:space:]]+[^[:space:]].*[[:space:]]—[[:space:]]+[^[:space:]].*$/
    }
    {
      raw = $0
      sub(/\r$/, "", raw)
      if (raw ~ /^[[:space:]]*$/ || raw ~ /^[[:space:]]*#/ || raw ~ /^[[:space:]]*\/\//) {
        prev_owner = is_owner(raw) ? 1 : 0
        next
      }
      line = raw
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
      print path "\t" bytes "\t" (prev_owner ? 1 : 0)
      prev_owner = 0
    }
  ' "$1"
}

load_sizes() {
  local file="$1"
  local -n dest=$2
  local track_owner="${3:-0}"
  local parsed path bytes own
  parsed=$(parse_baseline "$file") || exit $?
  while IFS=$'\t' read -r path bytes own; do
    [[ -z "${path:-}" ]] && continue
    dest["$path"]=$bytes
    if [[ "$track_owner" -eq 1 ]]; then
      owned["$path"]=$own
    fi
  done <<<"$parsed"
}

declare -A recorded=()
declare -A base_recorded=()
declare -A owned=()

size_failed=0
shrink_failed=0

if [[ -n "$baseline" ]]; then
  if [[ ! -f "$baseline" ]]; then
    echo "missing file: $baseline" >&2
    exit 2
  fi
  load_sizes "$baseline" recorded 1
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

  repo_baseline=$(repo_path "$baseline")
  base_file=$(mktemp)
  if git cat-file -e "${base_ref}:${repo_baseline}" 2>/dev/null; then
    git show "${base_ref}:${repo_baseline}" >"$base_file"
    load_sizes "$base_file" base_recorded
  fi
  rm -f "$base_file"

  label_source="$labels"
  if [[ "$labels_set" -eq 0 ]]; then
    label_source="${PR_LABELS:-}"
  fi
  has_growth_label=0
  IFS=',' read -r -a label_list <<<"$label_source"
  for label in "${label_list[@]}"; do
    trimmed=$(printf '%s' "$label" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    if [[ "$trimmed" == "baseline-growth-approved" ]]; then
      has_growth_label=1
    fi
  done

  new_paths=()
  new_unowned=0
  for path in "${!recorded[@]}"; do
    if [[ -n "${base_recorded[$path]+x}" && "${recorded[$path]}" -gt "${base_recorded[$path]}" ]]; then
      echo "$path: recorded size ${base_recorded[$path]} -> ${recorded[$path]}" >&2
      size_failed=1
    elif [[ -z "${base_recorded[$path]+x}" ]]; then
      new_paths+=("$path")
      if [[ "${owned[$path]}" -ne 1 ]]; then
        new_unowned=1
      fi
    fi
  done
  if [[ ${#new_paths[@]} -gt 0 && ( "$has_growth_label" -ne 1 || "$new_unowned" -ne 0 ) ]]; then
    for path in "${new_paths[@]}"; do
      echo "new baseline path: $path" >&2
    done
    echo "Growing the baseline hides new violations. Fix the code instead." >&2
    if [[ "$has_growth_label" -ne 1 ]]; then
      echo "Missing the exact PR label baseline-growth-approved." >&2
    fi
    if [[ "$new_unowned" -ne 0 ]]; then
      echo "Each new baseline path needs a \`# owner: <who> — <why>\` line directly before it (em dash; who and why non-empty)." >&2
    fi
    shrink_failed=1
  fi
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
