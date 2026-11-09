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
  (cd "$T/repo" && bash "$VENDOR" --source "$T/src" --sha "$1" --dest .agents/skills) >"$T/out" 2>"$T/err"
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

expect_fail() { # expect_fail <name> <needle> -- script args
  local name="$1" needle="$2"
  shift 2
  if bash "$VENDOR" "$@" >"$T/out" 2>"$T/err"; then
    bad "$name" "expected non-zero exit"
  elif grep -q "$needle" "$T/err"; then
    ok "$name"
  else
    bad "$name" "$(cat "$T/out" "$T/err")"
  fi
}

# 6 sha must be 40 hex characters
setup
expect_fail "sha short" "40-character" --source "$T/src" --sha abcdef --dest "$T/repo/.agents/skills"
expect_fail "sha branch name" "40-character" --source "$T/src" --sha main --dest "$T/repo/.agents/skills"
expect_fail "sha 39 chars" "40-character" --source "$T/src" --sha "${SHA:0:39}" --dest "$T/repo/.agents/skills"
expect_fail "sha 41 chars" "40-character" --source "$T/src" --sha "${SHA}a" --dest "$T/repo/.agents/skills"

# 7 dest must not be .., absolute, or the repo root
setup
expect_fail "dest dotdot" "must not contain" --source "$T/src" --sha "$SHA" --dest ..
expect_fail "dest absolute" "must be relative" --source "$T/src" --sha "$SHA" --dest /var/apstack-skills
expect_fail "dest dot" "repo root" --source "$T/src" --sha "$SHA" --dest .
WF="$ROOT/.github/workflows/vendor-apstack-skills.yml"
if grep -q 'github.event.repository.default_branch' "$WF" \
  && grep -F -q '""|"."|/*|*..*' "$WF"; then
  ok "workflow dest and checkout guards"
else
  bad "workflow dest and checkout guards" "workflow yaml missing ref or dest case"
fi

# 8 invalid manifest skill name
setup
mkdir -p "$T/repo/.agents/skills"
printf '%s\n' '{"source":"x","sha":"'"$SHA"'","skills":["Not_Valid"]}' > "$T/repo/.agents/skills/.apstack-vendor.json"
if run "$SHA"; then
  bad "invalid manifest name" "expected non-zero exit"
elif grep -q 'invalid skill name' "$T/err" \
  && [[ ! -e "$T/repo/.agents/skills/alpha" ]]; then
  ok "invalid manifest name"
else
  bad "invalid manifest name" "$(cat "$T/out" "$T/err")"
fi

# 9 empty source must not delete tracked skills
setup
run "$SHA" || { bad "empty source setup" "$(cat "$T/err")"; }
rm -rf "$T/src/skills"
mkdir -p "$T/src/skills"
mkdir -p "$T/src/skills/not-a-skill"
if run "$SHA"; then
  bad "empty source" "expected non-zero exit"
elif grep -q 'no valid skills' "$T/err" \
  && [[ -d "$T/repo/.agents/skills/alpha" ]] \
  && [[ -d "$T/repo/.agents/skills/manual-one" ]] \
  && grep -q '"alpha"' "$T/repo/.agents/skills/.apstack-vendor.json"; then
  ok "empty source"
else
  bad "empty source" "$(cat "$T/out" "$T/err")"
fi

# 10 skill directory that is a symlink
setup
rm -rf "$T/src/skills/alpha"
skill "$T/src/real-alpha" "linked"
ln -s "$T/src/real-alpha" "$T/src/skills/alpha"
if run "$SHA"; then
  bad "skill dir symlink" "expected non-zero exit"
elif grep -q 'must not be a symlink' "$T/err" \
  && [[ ! -e "$T/repo/.agents/skills/alpha" ]] \
  && [[ ! -e "$T/repo/.agents/skills/.apstack-vendor.json" ]]; then
  ok "skill dir symlink"
else
  bad "skill dir symlink" "$(cat "$T/out" "$T/err")"
fi

# 11 symlink that resolves outside the skill
setup
ln -s ../../outside "$T/src/skills/alpha/escape"
if run "$SHA"; then
  bad "symlink escape" "expected non-zero exit"
elif grep -q 'resolves outside' "$T/err" \
  && [[ ! -e "$T/repo/.agents/skills/.apstack-vendor.json" ]]; then
  ok "symlink escape"
else
  bad "symlink escape" "$(cat "$T/out" "$T/err")"
fi

# 12 relative symlink inside the skill is copied (clippy.toml pattern)
setup
mkdir -p "$T/src/skills/alpha/rules/rust/fixture"
printf 'rules\n' > "$T/src/skills/alpha/rules/rust/clippy.toml"
ln -s ../clippy.toml "$T/src/skills/alpha/rules/rust/fixture/clippy.toml"
if run "$SHA" \
  && [[ "$(readlink "$T/repo/.agents/skills/alpha/rules/rust/fixture/clippy.toml")" == "../clippy.toml" ]] \
  && grep -qx 'rules' "$T/repo/.agents/skills/alpha/rules/rust/fixture/clippy.toml"; then
  ok "internal relative symlink"
else
  bad "internal relative symlink" "$(cat "$T/out" "$T/err")"
fi

# 13 dangling symlink in dest counts as a collision
setup
ln -s missing "$T/repo/.agents/skills/alpha"
if run "$SHA"; then
  bad "dangling collision" "expected non-zero exit"
elif grep -q 'COLLISION' "$T/err" \
  && [[ -L "$T/repo/.agents/skills/alpha" ]] \
  && [[ "$(readlink "$T/repo/.agents/skills/alpha")" == missing ]] \
  && [[ ! -e "$T/repo/.agents/skills/.apstack-vendor.json" ]]; then
  ok "dangling collision"
else
  bad "dangling collision" "$(cat "$T/out" "$T/err")"
fi

# 14 tracked skill gains and loses a file; mode-only and symlink retargets count
setup
run "$SHA" || { bad "gain/lose setup" "$(cat "$T/err")"; }
printf 'extra\n' > "$T/src/skills/alpha/extra.txt"
if run "$SHA" && [[ "$(cat "$T/repo/.agents/skills/alpha/extra.txt")" == "extra" ]] \
  && grep -q 'changed=yes' "$T/out"; then
  ok "tracked skill gains a file"
else
  bad "tracked skill gains a file" "$(cat "$T/out" "$T/err")"
fi
rm -f "$T/src/skills/alpha/extra.txt"
if run "$SHA" && [[ ! -e "$T/repo/.agents/skills/alpha/extra.txt" ]] \
  && grep -q 'changed=yes' "$T/out"; then
  ok "tracked skill loses a file"
else
  bad "tracked skill loses a file" "$(cat "$T/out" "$T/err")"
fi
chmod 600 "$T/repo/.agents/skills/alpha/SKILL.md"
if run "$SHA" && grep -q 'changed=yes' "$T/out" \
  && [[ "$(stat -c '%a' "$T/repo/.agents/skills/alpha/SKILL.md")" == "$(stat -c '%a' "$T/src/skills/alpha/SKILL.md")" ]]; then
  ok "mode-only change"
else
  bad "mode-only change" "$(cat "$T/out" "$T/err")"
fi
ln -s SKILL.md "$T/src/skills/alpha/alias"
run "$SHA" || { bad "symlink retarget setup" "$(cat "$T/err")"; }
ln -sf extra-missing "$T/src/skills/alpha/alias"
if run "$SHA" && [[ "$(readlink "$T/repo/.agents/skills/alpha/alias")" == "extra-missing" ]] \
  && grep -q 'changed=yes' "$T/out"; then
  ok "symlink retarget"
else
  bad "symlink retarget" "$(cat "$T/out" "$T/err")"
fi

# 15 unchanged run documents the workflow skip of an open vendor PR
if grep -q 'changed=no; open vendor PR' "$WF" \
  && grep -q "steps.vendor.outputs.changed != 'true'" "$WF"; then
  ok "unchanged open PR skip"
else
  bad "unchanged open PR skip" "workflow missing changed=no skip"
fi

echo "---- $pass passed, $fail failed"
[[ $fail -eq 0 ]]
