#!/usr/bin/env bash
# Offline tests for scripts/vendor-apstack-skills.sh. Run:
#   bash tests/vendor-apstack-skills-test.sh
set -uo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
VENDOR="$ROOT/scripts/vendor-apstack-skills.sh"
T="$(mktemp -d)"; trap 'rm -rf "$T"' EXIT
pass=0; fail=0
ok()  { echo "PASS $1"; pass=$((pass+1)); }
bad() { echo "FAIL $1: $2"; fail=$((fail+1)); }
SHA="0123456789abcdef0123456789abcdef01234567"
SHA2="fedcba9876543210fedcba9876543210fedcba98"

skill() { # skill <dir> [body]
  mkdir -p "$1"
  printf -- '---\nname: %s\ndescription: %s\n---\n# %s\n' "$(basename "$1")" "${2:-test}" "${2:-body}" > "$1/SKILL.md"
}
manual() { # manual <dir>
  mkdir -p "$1/references"
  printf -- '---\nname: %s\ndescription: manual only\ndisable-model-invocation: true\n---\n# leave me\n' "$(basename "$1")" > "$1/SKILL.md"
  printf 'note\n' > "$1/references/note.md"
}

setup() {
  rm -rf "$T/src" "$T/repo"
  mkdir -p "$T/src/skills" "$T/repo"
  skill "$T/src/skills/alpha" "alpha-v1"
  manual "$T/src/skills/manual-one"
  skill "$T/repo/.agents/skills/repo-owned" "owned"
}
run() {
  bash "$VENDOR" --source "$T/src" --sha "$1" --dest "$T/repo/.agents/skills" >"$T/out" 2>"$T/err"
}
snap() { (cd "$T/repo" && find .agents -type f -print0 | LC_ALL=C sort -z | xargs -0 sha256sum) | sha256sum; }

# 1 first sync copies upstream skills, leaves repo-owned skills, keeps manual-only bytes
setup
if run "$SHA" && [[ -f "$T/repo/.agents/skills/alpha/SKILL.md" ]] \
  && grep -q 'alpha-v1' "$T/repo/.agents/skills/alpha/SKILL.md" \
  && cmp -s "$T/src/skills/manual-one/SKILL.md" "$T/repo/.agents/skills/manual-one/SKILL.md" \
  && cmp -s "$T/src/skills/manual-one/references/note.md" "$T/repo/.agents/skills/manual-one/references/note.md" \
  && grep -q 'disable-model-invocation: true' "$T/repo/.agents/skills/manual-one/SKILL.md" \
  && grep -q 'owned' "$T/repo/.agents/skills/repo-owned/SKILL.md" \
  && grep -q "$SHA" "$T/repo/.agents/skills/.apstack-vendor.json" \
  && grep -q '"alpha"' "$T/repo/.agents/skills/.apstack-vendor.json" \
  && grep -q '"manual-one"' "$T/repo/.agents/skills/.apstack-vendor.json" \
  && ! grep -q 'repo-owned' "$T/repo/.agents/skills/.apstack-vendor.json"; then
  ok "first sync"
else
  bad "first sync" "$(cat "$T/out" "$T/err")"
fi

# 2 update replaces a tracked skill and leaves the repo-owned skill
setup
run "$SHA" || { bad "update setup" "$(cat "$T/err")"; }
skill "$T/src/skills/alpha" "alpha-v2"
if run "$SHA2" && grep -q 'alpha-v2' "$T/repo/.agents/skills/alpha/SKILL.md" \
  && grep -q "$SHA2" "$T/repo/.agents/skills/.apstack-vendor.json" \
  && grep -q 'owned' "$T/repo/.agents/skills/repo-owned/SKILL.md" \
  && cmp -s "$T/src/skills/manual-one/SKILL.md" "$T/repo/.agents/skills/manual-one/SKILL.md"; then
  ok "update"
else
  bad "update" "$(cat "$T/out" "$T/err")"
fi

# 3 upstream deletion removes only the tracked skill
setup
run "$SHA" || { bad "delete setup" "$(cat "$T/err")"; }
rm -rf "$T/src/skills/alpha"
if run "$SHA" && [[ ! -e "$T/repo/.agents/skills/alpha" ]] \
  && [[ -d "$T/repo/.agents/skills/manual-one" ]] \
  && [[ -d "$T/repo/.agents/skills/repo-owned" ]] \
  && ! grep -q '"alpha"' "$T/repo/.agents/skills/.apstack-vendor.json" \
  && grep -q '"manual-one"' "$T/repo/.agents/skills/.apstack-vendor.json"; then
  ok "upstream delete"
else
  bad "upstream delete" "$(cat "$T/out" "$T/err"; find "$T/repo/.agents/skills" -maxdepth 2 -type d)"
fi

# 4 name collision with a repo-owned skill fails and writes nothing
setup
skill "$T/repo/.agents/skills/alpha" "local-alpha"
before="$(snap)"
if run "$SHA"; then
  bad "collision" "expected non-zero exit"
elif grep -q 'COLLISION' "$T/err" && [[ "$(snap)" == "$before" ]] \
  && grep -q 'local-alpha' "$T/repo/.agents/skills/alpha/SKILL.md" \
  && [[ ! -e "$T/repo/.agents/skills/manual-one" ]] \
  && [[ ! -e "$T/repo/.agents/skills/.apstack-vendor.json" ]]; then
  ok "collision"
else
  bad "collision" "$(cat "$T/out" "$T/err")"
fi

# 5 idempotency: a second run with the same SHA changes no bytes
setup
run "$SHA" || { bad "idempotency setup" "$(cat "$T/err")"; }
s1="$(snap)"
if run "$SHA" && [[ "$(snap)" == "$s1" ]] && grep -q 'changed=no' "$T/out"; then
  ok "idempotency"
else
  bad "idempotency" "$(cat "$T/out" "$T/err")"
fi

echo "---- $pass passed, $fail failed"
[[ $fail -eq 0 ]]
