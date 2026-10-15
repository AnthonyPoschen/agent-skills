#!/usr/bin/env bash
# Offline tests for scripts/sync-skills.sh: temp dirs, a local bare remote,
# and a stub PR opener. Run: bash tests/sync-skills-test.sh
set -uo pipefail
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SYNC="$ROOT/scripts/sync-skills.sh"
T="$(mktemp -d)"; trap 'rm -rf "$T"' EXIT
export GIT_CONFIG_GLOBAL="$T/gitconfig"
git config --global user.name test; git config --global user.email t@t; git config --global init.defaultBranch master
pass=0; fail=0
ok()  { echo "PASS $1"; pass=$((pass+1)); }
bad() { echo "FAIL $1: $2"; fail=$((fail+1)); }
skill() { mkdir -p "$1"; printf -- '---\nname: %s\ndescription: %s\n---\n# %s\n' "$(basename "$1")" "${2:-test}" "${2:-body}" > "$1/SKILL.md"; }

setup() {
  rm -rf "$T/w"; mkdir -p "$T/w"; cd "$T/w"
  git init -q --bare remote.git
  git clone -q remote.git seed 2>/dev/null; cd seed
  cp "$ROOT/package.json" .; mkdir -p scripts; cp "$ROOT/scripts/validate-skills.py" scripts/
  printf '# excluded\nbotonly\n' > sync-exclude.txt
  skill skills/alpha; skill skills/beta
  git add -A; git commit -qm init; git push -q origin master; cd ..
  mkdir box; cp -a seed/skills/alpha seed/skills/beta box/
  skill box/botonly "bot"
  echo 'printf "%s|%s\n" "$1" "$2" >> '"$T/w/prs.log"'; echo stub-pr' > prstub.sh; chmod +x prstub.sh
  : > prs.log
}
run() { bash "$SYNC" --repo "$T/w/clone" --repo-url "$T/w/remote.git" --box "$T/w/box" \
          --state "$T/w/state/state" --pr-cmd "$T/w/prstub.sh" "$@" > "$T/w/out" 2>&1; }
repo_commit() { (cd "$T/w/seed" && git pull -q && eval "$1" && git add -A && git commit -qm change && git push -q origin master); }
snap() { (cd "$T/w" && find box state -type f -exec sha256sum {} + 2>/dev/null | sort; git -C remote.git for-each-ref) | sha256sum; }
prs() { wc -l < "$T/w/prs.log"; }

# 1 no-op
setup; run && run; [[ $(prs) == 0 ]] && ! grep -qE '^(to-|del-|FLAG)' out && ok "no-op run" || bad "no-op run" "$(cat out)"
# 2 repo-only change
setup; run; repo_commit 'skill skills/alpha v2'; run
grep -q '# v2' box/alpha/SKILL.md && [[ $(prs) == 0 ]] && ok "repo-only change copied to box" || bad "repo-only change" "$(cat out)"
# 3 box-only change
setup; run; skill box/beta boxv2; run
git -C remote.git rev-parse -q --verify skills-sync/beta >/dev/null && grep -q '^skills-sync/beta|' prs.log \
  && git -C remote.git show skills-sync/beta:skills/beta/SKILL.md | grep -q boxv2 \
  && git -C remote.git log -1 --format=%B skills-sync/beta | grep -q '^Skills-Sync: ' \
  && ! git -C remote.git show master:skills/beta/SKILL.md | grep -q boxv2 \
  && ok "box-only change opens PR (trailer present, master untouched)" || bad "box-only change" "$(cat out)"
# 3b merge that PR -> settles
git -C seed pull -q; git -C seed merge -q origin/skills-sync/beta 2>/dev/null || (cd seed && git fetch -q && git merge -q origin/skills-sync/beta); git -C seed push -q origin master
: > prs.log; run; run; [[ $(prs) == 0 ]] && ! grep -qE '^(to-|FLAG)' out && ok "merged box PR settles (no further action)" || bad "merged PR settle" "$(cat out)"
# 4 conflict
setup; run; skill box/alpha boxside; repo_commit 'skill skills/alpha reposide'; run
grep -q '^FLAG conflict alpha' out && grep -q boxside box/alpha/SKILL.md && git -C remote.git show master:skills/alpha/SKILL.md | grep -q reposide \
  && grep -q '^skills-sync/alpha|skills-sync: conflict' prs.log && ok "conflict flags, opens PR, leaves both sides" || bad "conflict" "$(cat out)"
# 5a repo deletion, box unchanged
setup; run; repo_commit 'git rm -rq skills/beta'; run
[[ ! -d box/beta ]] && ok "repo deletion carried to box" || bad "repo deletion" "$(cat out)"
# 5b box deletion, repo unchanged
setup; run; rm -rf box/beta; run
grep -q '^skills-sync/beta|skills-sync: box-delete' prs.log && ! git -C remote.git cat-file -e skills-sync/beta:skills/beta/SKILL.md 2>/dev/null \
  && ok "box deletion arrives as PR" || bad "box deletion" "$(cat out)"
# 6a deletion vs edit (repo deleted, box edited)
setup; run; skill box/beta edited; repo_commit 'git rm -rq skills/beta'; run
grep -q '^FLAG conflict beta' out && grep -q edited box/beta/SKILL.md && ok "repo delete vs box edit = conflict" || bad "delete vs edit (repo)" "$(cat out)"
# 6b box deleted, repo edited
setup; run; rm -rf box/beta; repo_commit 'skill skills/beta repoedit'; run
grep -q '^FLAG conflict beta' out && [[ ! -d box/beta ]] && ok "box delete vs repo edit = conflict" || bad "delete vs edit (box)" "$(cat out)"
# 7 excluded skill (box-only bot skill, then also edited on repo)
setup; run; skill box/botonly changed; repo_commit 'skill skills/botonly repo'; run
grep -q '^exclude  botonly' out && grep -q changed box/botonly/SKILL.md && [[ $(prs) == 0 ]] \
  && ok "excluded skill never synced" || bad "exclude" "$(cat out)"
# 8 loop guard: run right after a sync does nothing
setup; run; repo_commit 'skill skills/gamma new'; skill box/delta boxnew; run; : > prs.log; s1=$(snap); run
[[ $(prs) -le 1 ]] && ! grep -qE '^(to-box|del-|FLAG)' out && [[ "$(git -C remote.git rev-parse skills-sync/delta)" == "$(cd $T/w && git -C remote.git rev-parse skills-sync/delta)" ]] \
  && [[ "$(snap)" == "$s1" ]] && ok "loop guard: second run makes no new writes" || bad "loop guard" "$(cat out)"
# 9 dry-run changes nothing
setup; run; repo_commit 'skill skills/alpha v3'; skill box/beta b3; s1=$(snap); run --dry-run
grep -q '^to-box   alpha' out && grep -q '^to-pr    beta' out && [[ "$(snap)" == "$s1" ]] && [[ $(prs) == 0 ]] \
  && ok "dry-run prints plan, changes nothing" || bad "dry-run" "$(cat out)"
# 10 idempotence: repeated runs converge to identical state
setup; repo_commit 'skill skills/gamma g'; run; run; s1=$(snap); run; run
[[ "$(snap)" == "$s1" ]] && ! grep -qE '^(to-|del-|FLAG)' out && ok "idempotent" || bad "idempotent" "$(cat out)"
# 11 invalid repo skill is not written to box
setup; run; repo_commit 'printf "no frontmatter\n" > skills/alpha/SKILL.md'; run
grep -q '^name: alpha' box/alpha/SKILL.md && grep -q 'SKIP alpha' out && ok "invalid skill rejected before write" || bad "validate" "$(cat out)"
# 12 lock stops concurrent run
setup; mkdir -p state; exec 8>state/state.lock; flock 8; run; rc=$?; exec 8>&-
[[ $rc -ne 0 ]] && grep -q 'another sync' out && ok "lock blocks concurrent run" || bad "lock" "$(cat out)"

# 13 skill merge: repo folds alpha+beta into a new parent in one commit
setup; run; repo_commit 'git rm -rq skills/alpha skills/beta && skill skills/merged parent && mkdir -p skills/merged/references && echo ref > skills/merged/references/alpha.md'; run; run
[[ ! -d box/alpha && ! -d box/beta && -f box/merged/references/alpha.md && $(prs) == 0 ]] && ! grep -qE '^FLAG' out \
  && ok "skill merge replaces old skills on box without PRs" || bad "skill merge" "$(cat out)"
# 14 link script prunes links to merged-away skills, keeps foreign links
L="$T/link"; mkdir -p "$L/src"; skill "$L/src/alpha"; skill "$L/src/beta"; ln -s /elsewhere "$L/foreign" 2>/dev/null; mkdir -p "$L/dest"; ln -s /elsewhere "$L/dest/foreign"
AGENT_SKILLS_SRC="$L/src" AGENT_SKILLS_DIR="$L/dest" sh "$ROOT/scripts/link-agent-skills.sh" >/dev/null
rm -rf "$L/src/alpha"; skill "$L/src/merged"
AGENT_SKILLS_SRC="$L/src" AGENT_SKILLS_DIR="$L/dest" sh "$ROOT/scripts/link-agent-skills.sh" >/dev/null
[[ ! -L "$L/dest/alpha" && -L "$L/dest/beta" && -L "$L/dest/merged" && -L "$L/dest/foreign" ]] \
  && ok "link script prunes stale links, keeps foreign ones" || bad "link prune" "$(ls -l "$L/dest")"
echo "---- $pass passed, $fail failed"; [[ $fail -eq 0 ]]
