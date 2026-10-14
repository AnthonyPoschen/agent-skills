#!/usr/bin/env bash
# Two-way sync between agent-skills `skills/` (origin master) and the box
# workflows folder. Three-way compare against a state file of last-synced
# hashes. Repo changes are copied to the box; box changes become PRs.
# Run with --help for options. Tests: tests/sync-skills-test.sh
set -euo pipefail

STATE_FILE="${SKILLS_SYNC_STATE:-/home/box/agent-data/skills-sync/state}"
BOX_DIR="${SKILLS_SYNC_BOX:-/home/box/agent-data/workflows}"
REPO_DIR="${SKILLS_SYNC_REPO:-/home/box/agent-data/skills-sync/agent-skills}"
REPO_URL="${SKILLS_SYNC_REPO_URL:-https://github.com/AnthonyPoschen/agent-skills.git}"
BRANCH="${SKILLS_SYNC_BRANCH:-master}"
PR_CMD="${SKILLS_SYNC_PR_CMD:-}"   # stub hook: called as <cmd> <branch> <title> <body-file>
LOCK_FILE=""
DRY_RUN=0
RUN_ID="$(date -u +%Y%m%dT%H%M%SZ)-$$"

usage() {
  cat <<USAGE
Usage: sync-skills.sh [--repo DIR] [--repo-url URL] [--box DIR] [--state FILE]
                      [--branch NAME] [--lock FILE] [--pr-cmd CMD] [--dry-run]

  --repo DIR     Local clone of agent-skills used for fetch/push (cloned if missing).
                 Default: $REPO_DIR
  --repo-url URL Clone URL when --repo does not exist yet.
  --box DIR      Box skills folder. Default: $BOX_DIR
  --state FILE   Last-synced state file (keep outside the repo). Default: $STATE_FILE
  --branch NAME  Repo branch treated as canonical. Default: $BRANCH
  --lock FILE    Lock file. Default: <state>.lock
  --pr-cmd CMD   Command that opens/updates a PR: CMD <branch> <title> <body-file>.
                 Default: gh (draft PR, edits the open one if it exists).
  --dry-run      Print the plan; change nothing (no box writes, pushes, state).

Output: one line per skill action; conflicts print a line starting "FLAG".
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --repo) REPO_DIR="$2"; shift 2;;
    --repo-url) REPO_URL="$2"; shift 2;;
    --box) BOX_DIR="$2"; shift 2;;
    --state) STATE_FILE="$2"; shift 2;;
    --branch) BRANCH="$2"; shift 2;;
    --lock) LOCK_FILE="$2"; shift 2;;
    --pr-cmd) PR_CMD="$2"; shift 2;;
    --dry-run) DRY_RUN=1; shift;;
    -h|--help) usage; exit 0;;
    *) echo "unknown argument: $1 (see --help)" >&2; exit 2;;
  esac
done
LOCK_FILE="${LOCK_FILE:-$STATE_FILE.lock}"

die() { echo "ERROR: $*" >&2; exit 1; }
log() { echo "$*"; }

[[ -d "$BOX_DIR" ]] || die "box dir $BOX_DIR missing; pass --box"
mkdir -p "$(dirname "$STATE_FILE")" "$(dirname "$LOCK_FILE")"
exec 9>"$LOCK_FILE"
flock -n 9 || die "another sync holds $LOCK_FILE; retry later"

TMP="$(mktemp -d)"
cleanup() {
  [[ -n "${WT:-}" && -d "${WT:-}" ]] && git -C "$REPO_DIR" worktree remove --force "$WT" >/dev/null 2>&1 || true
  rm -rf "$TMP"
}
trap cleanup EXIT

# --- repo side: fetch canonical branch into a scratch tree -------------------
if [[ ! -d "$REPO_DIR/.git" ]]; then
  [[ $DRY_RUN -eq 1 ]] && REPO_DIR="$TMP/clone"
  git clone -q "$REPO_URL" "$REPO_DIR" || die "clone of $REPO_URL failed; check network/auth"
fi
git -C "$REPO_DIR" fetch -q origin "+refs/heads/$BRANCH:refs/remotes/origin/$BRANCH" \
  || die "fetch origin $BRANCH failed in $REPO_DIR"
MASTER_SHA="$(git -C "$REPO_DIR" rev-parse "origin/$BRANCH")"
REPO_SNAP="$TMP/repo"; mkdir -p "$REPO_SNAP"
git -C "$REPO_DIR" archive "origin/$BRANCH" | tar -x -C "$REPO_SNAP"
mkdir -p "$REPO_SNAP/skills"

# Exclude list comes from canonical master so it is reviewed like any change.
# Falls back to the copy beside this script if master has none yet.
declare -A EXCL=()
EXCL_FILE="$REPO_SNAP/sync-exclude.txt"
[[ -f "$EXCL_FILE" ]] || EXCL_FILE="$(cd "$(dirname "$0")/.." && pwd)/sync-exclude.txt"
if [[ -f "$EXCL_FILE" ]]; then
  while IFS= read -r l; do l="${l%%#*}"; l="${l//[[:space:]]/}"; [[ -n "$l" ]] && EXCL[$l]=1; done \
    < "$EXCL_FILE"
fi

# Content hash of a skill folder: relative paths + file contents, order-stable.
skill_hash() {
  local d="$1"
  [[ -f "$d/SKILL.md" ]] || { echo "-"; return; }
  (cd "$d" && find . -type f ! -path '*/.git/*' -print0 | LC_ALL=C sort -z \
    | xargs -0 -r sha256sum) | sha256sum | cut -c1-16
}

declare -A BASE=()
if [[ -f "$STATE_FILE" ]]; then
  while read -r k n h; do [[ "$k" == skill ]] && BASE[$n]="$h"; done < "$STATE_FILE"
fi

valid_skill() {   # SKILL.md with frontmatter carrying name + description
  local f="$1/SKILL.md"
  [[ -f "$f" ]] || { echo "missing SKILL.md in $1" >&2; return 1; }
  [[ "$(head -n1 "$f")" == "---" ]] || { echo "$f: no frontmatter" >&2; return 1; }
  awk 'NR>1&&/^---$/{e=1;exit} END{exit !e}' "$f" || { echo "$f: frontmatter not closed" >&2; return 1; }
  awk 'NR>1&&/^---$/{exit} /^name:/{n=1} /^description:/{d=1} END{exit !(n&&d)}' "$f" \
    || { echo "$f: frontmatter needs name and description" >&2; return 1; }
}

repo_validate() {   # npm run validate where the repo has it (links via temp dir)
  local tree="$1"
  [[ -f "$tree/package.json" ]] && grep -q '"validate"' "$tree/package.json" || return 0
  command -v npm >/dev/null || { echo "npm missing; cannot run npm run validate" >&2; return 1; }
  local links="$TMP/links"; rm -rf "$links"; mkdir -p "$links"
  for d in "$tree"/skills/*/; do ln -s "$d" "$links/$(basename "$d")"; done
  (cd "$tree" && AGENT_SKILLS_DIR="$links" npm run -s validate) >&2
}

default_pr() {   # <branch> <title> <body-file>; idempotent: edit if open
  local num
  num="$(cd "$REPO_DIR" && gh pr list --head "$1" --state open --json number -q '.[0].number')"
  if [[ -n "$num" ]]; then
    (cd "$REPO_DIR" && gh pr edit "$num" --title "$2" --body-file "$3" >/dev/null) && echo "updated PR #$num"
  else
    (cd "$REPO_DIR" && gh pr create --draft --base "$BRANCH" --head "$1" --title "$2" --body-file "$3")
  fi
}

open_pr() {   # <skill> <kind: box-change|box-delete|conflict> <note>
  local s="$1" kind="$2" note="$3" br="skills-sync/$1"
  WT="$TMP/wt-$s"
  git -C "$REPO_DIR" worktree add -q -f -B "$br" "$WT" "origin/$BRANCH"
  rm -rf "${WT:?}/skills/$s"
  [[ -d "$BOX_DIR/$s" && -f "$BOX_DIR/$s/SKILL.md" ]] && cp -a "$BOX_DIR/$s" "$WT/skills/$s"
  if [[ -d "$WT/skills/$s" ]]; then valid_skill "$WT/skills/$s" || { log "SKIP $s: box copy invalid"; git -C "$REPO_DIR" worktree remove --force "$WT"; WT=""; return 0; }; fi
  repo_validate "$WT" || { log "SKIP $s: npm run validate failed with box copy"; git -C "$REPO_DIR" worktree remove --force "$WT"; WT=""; return 0; }
  git -C "$WT" add -A "skills/$s"
  local title="skills-sync: ${kind} ${s}"
  git -C "$WT" -c user.name="skills-sync" -c user.email="skills-sync@users.noreply.github.com" \
    commit -q -m "$title" -m "$note" -m "Skills-Sync: $RUN_ID" || true
  # Push only when the content differs from what the remote branch already has.
  local remote_tree=""
  remote_tree="$(git -C "$REPO_DIR" ls-remote origin "refs/heads/$br" | cut -f1)"
  if [[ -n "$remote_tree" ]] && git -C "$REPO_DIR" fetch -q origin "$br" 2>/dev/null \
     && [[ "$(git -C "$WT" rev-parse HEAD^{tree})" == "$(git -C "$REPO_DIR" rev-parse FETCH_HEAD^{tree})" ]]; then
    log "  branch $br already up to date"
  else
    git -C "$WT" push -q -f origin "$br:refs/heads/$br" || die "push $br failed; check repo auth"
  fi
  local body="$TMP/body-$s"
  { echo "Automated two-way skills sync (run \`$RUN_ID\`)."; echo
    echo "$note"; echo
    [[ "$kind" == conflict ]] && { echo "**CONFLICT:** both the box and \`$BRANCH\` changed \`$s\` since the last sync. This PR carries the box version; neither side was changed. Resolve here (edit/merge or close), and the next sync settles once both sides match."; echo; }
    echo "Merge to accept. Base sha last synced: \`${LAST_SHA:-none}\`."; } > "$body"
  local out
  if [[ -n "$PR_CMD" ]]; then out="$($PR_CMD "$br" "$title" "$body")"; else out="$(default_pr "$br" "$title" "$body")"; fi
  log "  pr: $out"
  git -C "$REPO_DIR" worktree remove --force "$WT"; WT=""
}

copy_to_box() {   # repo -> box, atomic per skill
  local s="$1"
  valid_skill "$REPO_SNAP/skills/$s" || { log "SKIP $s: repo copy invalid"; return 1; }
  rm -rf "$BOX_DIR/.sync-$s"; cp -a "$REPO_SNAP/skills/$s" "$BOX_DIR/.sync-$s"
  rm -rf "${BOX_DIR:?}/$s"; mv "$BOX_DIR/.sync-$s" "$BOX_DIR/$s"
}

LAST_SHA=""; [[ -f "$STATE_FILE" ]] && LAST_SHA="$(awk '$1=="master"{print $2}' "$STATE_FILE")"
declare -A NEWBASE=()
for k in "${!BASE[@]}"; do NEWBASE[$k]="${BASE[$k]}"; done

names="$( { ls -1 "$REPO_SNAP/skills"; ls -1 "$BOX_DIR"; printf '%s\n' "${!BASE[@]}"; } 2>/dev/null \
          | grep -E '^[a-z0-9][a-z0-9-]*$' | sort -u)"
[[ $DRY_RUN -eq 1 ]] && log "DRY-RUN plan (run $RUN_ID, $BRANCH @ ${MASTER_SHA:0:7}, base ${LAST_SHA:0:7})"
flags=0
for s in $names; do
  if [[ -n "${EXCL[$s]:-}" ]]; then log "exclude  $s"; continue; fi
  r="$(skill_hash "$REPO_SNAP/skills/$s")"; b="$(skill_hash "$BOX_DIR/$s")"; base="${BASE[$s]:-}"
  if [[ "$r" == "$b" ]]; then
    if [[ "$r" == "-" ]]; then unset 'NEWBASE[$s]'; else NEWBASE[$s]="$r"; fi
    continue                         # in sync (also how the sync's own writes settle)
  fi
  if [[ -z "$base" || "$base" == "-" ]]; then
    if   [[ "$b" == "-" ]]; then act=repo-new
    elif [[ "$r" == "-" ]]; then act=box-new
    else act=conflict; fi
  elif [[ "$r" == "$base" ]]; then act=$([[ "$b" == "-" ]] && echo box-delete || echo box-change)
  elif [[ "$b" == "$base" ]]; then act=$([[ "$r" == "-" ]] && echo repo-delete || echo repo-change)
  else act=conflict; fi

  case "$act" in
    repo-new|repo-change)
      log "to-box   $s ($act)"
      if [[ $DRY_RUN -eq 0 ]] && copy_to_box "$s"; then NEWBASE[$s]="$r"; fi;;
    repo-delete)
      log "del-box  $s (deleted on $BRANCH, unchanged on box)"
      if [[ $DRY_RUN -eq 0 ]]; then rm -rf "${BOX_DIR:?}/$s"; unset 'NEWBASE[$s]'; fi;;
    box-new|box-change|box-delete)
      log "to-pr    $s ($act)"
      [[ $DRY_RUN -eq 0 ]] && open_pr "$s" "$act" "Box changed \`$s\` ($act); \`$BRANCH\` unchanged since last sync.";;
    conflict)
      flags=1
      log "FLAG conflict $s: repo=$r box=$b base=${base:-none}; both sides left untouched"
      [[ $DRY_RUN -eq 0 ]] && open_pr "$s" conflict "Repo hash \`$r\`, box hash \`$b\`, base \`${base:-none}\`.";;
  esac
done

if [[ $DRY_RUN -eq 0 ]]; then
  { echo "# skills-sync state; written by scripts/sync-skills.sh, do not edit"
    echo "master $MASTER_SHA"
    for k in $(printf '%s\n' "${!NEWBASE[@]}" | sort); do echo "skill $k ${NEWBASE[$k]}"; done
  } > "$STATE_FILE.tmp" && mv "$STATE_FILE.tmp" "$STATE_FILE"
fi
log "done (flags=$flags)"
