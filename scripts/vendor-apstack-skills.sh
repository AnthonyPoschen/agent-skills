#!/usr/bin/env bash
# Copy skills/ from a pinned agent-skills checkout into a repo skills folder.
# Only skills recorded in <dest>/.apstack-vendor.json are updated or removed.
# A skill already in dest and absent from that manifest is a hard error.
# Tests: tests/vendor-apstack-skills-test.sh
set -euo pipefail

SOURCE=""
SHA=""
DEST=".agents/skills"
UPSTREAM="https://github.com/AnthonyPoschen/agent-skills"

usage() {
  cat <<USAGE
Usage: vendor-apstack-skills.sh --source DIR --sha SHA [--dest RELATIVE-DIR]

  --source DIR  Checkout of agent-skills whose skills/ tree is copied.
  --sha SHA     Full 40-character commit SHA of that checkout. Recorded in
                the vendor manifest. The script does not fetch git.
  --dest DIR    Destination skills directory. Default: .agents/skills
                Cursor loads project skills from .agents/skills/ and from
                .cursor/skills/. Cloud Agents in repos that already commit
                .agents/skills/ (genos) pick that path up.

Refuses to overwrite or delete a skill that is not already listed in
<dest>/.apstack-vendor.json. Copies SKILL.md bytes unchanged, including
disable-model-invocation. Prints changed=yes or changed=no.
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --source) SOURCE="$2"; shift 2;;
    --sha) SHA="$2"; shift 2;;
    --dest) DEST="$2"; shift 2;;
    -h|--help) usage; exit 0;;
    *) echo "unknown argument: $1 (see --help)" >&2; exit 2;;
  esac
done

die() { echo "ERROR: $*" >&2; exit 1; }

[[ -n "$SOURCE" && -n "$SHA" ]] || die "--source and --sha are required"
[[ -d "$SOURCE/skills" ]] || die "source skills dir missing: $SOURCE/skills"
[[ "$SHA" =~ ^[0-9a-fA-F]{40}$ ]] || die "sha must be a full 40-character hex commit, got: $SHA"
[[ -n "$DEST" && "$DEST" != "/" ]] || die "dest must be a skills directory"
case "$DEST" in
  *..*) die "dest must not contain ..";;
esac
SHA="$(printf '%s' "$SHA" | tr 'A-F' 'a-f')"

skill_hash() {
  local d="$1"
  [[ -d "$d" ]] || { echo "-"; return; }
  (cd "$d" && find . -type f ! -path '*/.git/*' -print0 | LC_ALL=C sort -z \
    | xargs -0 -r sha256sum) | sha256sum | cut -c1-16
}

declare -A TRACKED=()
MANIFEST="$DEST/.apstack-vendor.json"
if [[ -f "$MANIFEST" ]]; then
  while IFS= read -r name; do
    [[ "$name" =~ ^[a-z0-9][a-z0-9-]*$ ]] || die "vendor manifest has invalid skill name: $name"
    TRACKED["$name"]=1
  done < <(python3 - "$MANIFEST" <<'PY'
import json, sys
path = sys.argv[1]
try:
    data = json.loads(open(path, encoding="utf-8").read())
except json.JSONDecodeError as exc:
    raise SystemExit(f"vendor manifest is not JSON: {exc}")
skills = data.get("skills") if isinstance(data, dict) else None
if not isinstance(skills, list) or not all(isinstance(item, str) for item in skills):
    raise SystemExit("vendor manifest skills must be a list of strings")
for item in skills:
    print(item)
PY
)
fi

declare -A SOURCE_SKILLS=()
shopt -s nullglob
for dir in "$SOURCE/skills"/*/; do
  name="$(basename "$dir")"
  [[ "$name" =~ ^[a-z0-9][a-z0-9-]*$ ]] || continue
  [[ -f "$dir/SKILL.md" ]] || continue
  SOURCE_SKILLS["$name"]=1
done

for name in "${!SOURCE_SKILLS[@]}"; do
  if [[ -e "$DEST/$name" && -z "${TRACKED[$name]:-}" ]]; then
    echo "COLLISION: $name exists in $DEST and is not in the apstack vendor manifest; refusing to overwrite" >&2
    exit 1
  fi
done

mkdir -p "$DEST"
changed=0
for name in "${!SOURCE_SKILLS[@]}"; do
  if [[ "$(skill_hash "$SOURCE/skills/$name")" == "$(skill_hash "$DEST/$name")" ]]; then
    continue
  fi
  stage="$DEST/.staging-$name"
  rm -rf "$stage"
  cp -a "$SOURCE/skills/$name" "$stage"
  rm -rf "${DEST:?}/$name"
  mv "$stage" "$DEST/$name"
  changed=1
done

for name in "${!TRACKED[@]}"; do
  if [[ -z "${SOURCE_SKILLS[$name]:-}" && -e "$DEST/$name" ]]; then
    rm -rf "${DEST:?}/$name"
    changed=1
  fi
done

names="$(printf '%s\n' "${!SOURCE_SKILLS[@]}" | LC_ALL=C sort)"
new_manifest="$(mktemp)"
NAMES="$names" SHA="$SHA" UPSTREAM="$UPSTREAM" python3 - "$new_manifest" <<'PY'
import json, os, sys
skills = [line for line in os.environ["NAMES"].splitlines() if line]
doc = {
    "source": os.environ["UPSTREAM"],
    "sha": os.environ["SHA"],
    "skills": skills,
}
with open(sys.argv[1], "w", encoding="utf-8") as handle:
    json.dump(doc, handle, indent=2)
    handle.write("\n")
PY
if [[ ! -f "$MANIFEST" ]] || ! cmp -s "$new_manifest" "$MANIFEST"; then
  mv "$new_manifest" "$MANIFEST"
  changed=1
else
  rm -f "$new_manifest"
fi

if [[ "$changed" -eq 1 ]]; then
  echo "changed=yes"
  flag=true
else
  echo "changed=no"
  flag=false
fi
if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
  echo "changed=$flag" >> "$GITHUB_OUTPUT"
fi
