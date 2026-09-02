#!/usr/bin/env bash
set -euo pipefail

marker='git-worktree-cleanup-hook'

die() {
  printf '%s\n' "$*" >&2
  exit 1
}

usage() {
  cat <<'EOF'
Usage:
  worktree-lifecycle.sh sync <branch>
  worktree-lifecycle.sh cleanup-merged [--apply] [--base <ref>]
  worktree-lifecycle.sh remove <branch>
  worktree-lifecycle.sh install-post-merge-hook
EOF
}

primary_checkout() {
  git worktree list --porcelain | sed -n 's/^worktree //p' | sed -n '1p'
}

setup_repository() {
  primary=$(primary_checkout)
  [ -n "$primary" ] || die 'Could not find the primary Git checkout.'
  primary=$(cd "$primary" && pwd -P)

  git_root="$HOME/git"
  worktree_root="$HOME/worktree"
  case "$primary" in
    "$git_root"/*) relative=${primary#"$git_root"/} ;;
    *) die "Primary checkout is outside $git_root: $primary" ;;
  esac
  mirror_root="$worktree_root/$relative"
}

registered_worktree() {
  local target=$1
  local path
  while IFS= read -r path; do
    [ -d "$path" ] || continue
    [ "$(cd "$path" && pwd -P)" = "$target" ] && return 0
  done < <(git -C "$primary" worktree list --porcelain | sed -n 's/^worktree //p')
  return 1
}

target_for_branch() {
  local branch=$1
  git check-ref-format --branch "$branch" >/dev/null || die "Invalid branch name: $branch"
  printf '%s/%s\n' "$mirror_root" "$branch"
}

require_clean() {
  local target=$1
  [ -z "$(git -C "$target" status --porcelain)" ]
}

require_branch() {
  local target=$1 branch=$2 actual
  actual=$(git -C "$target" branch --show-current)
  [ "$actual" = "$branch" ] || die "Worktree at $target is on ${actual:-a detached HEAD}, not $branch"
}

sync_branch() {
  local branch=$1 target upstream remote
  target=$(target_for_branch "$branch")
  [ -d "$target" ] && registered_worktree "$target" || die "No registered mirrored worktree for $branch at $target"
  require_branch "$target" "$branch"
  require_clean "$target" || die "Worktree is dirty; not updating: $target"
  upstream=$(git -C "$target" rev-parse --abbrev-ref --symbolic-full-name '@{upstream}' 2>/dev/null) || die "Branch has no upstream: $branch"
  remote=$(git -C "$target" config --get "branch.$branch.remote" 2>/dev/null || true)
  [ -n "$remote" ] || die "Branch has no configured remote: $branch"
  git -C "$target" fetch --prune "$remote"
  git -C "$target" merge --ff-only "$upstream"
  printf 'Synchronized %s at %s\n' "$branch" "$target"
}

default_base() {
  local remote base
  remote=$(git -C "$primary" remote | sed -n '1p')
  [ -n "$remote" ] || die 'No remote found; pass --base <ref>.'
  git -C "$primary" fetch --prune "$remote"
  base=$(git -C "$primary" symbolic-ref --quiet --short "refs/remotes/$remote/HEAD" 2>/dev/null || true)
  [ -n "$base" ] || die "Remote $remote has no HEAD; pass --base <ref>."
  printf '%s\n' "$base"
}

cleanup_merged() {
  local apply=false base='' path branch expected
  while [ "$#" -gt 0 ]; do
    case "$1" in
      --apply) apply=true ;;
      --base) shift; [ "$#" -gt 0 ] || die '--base needs a ref'; base=$1 ;;
      *) die "Unknown cleanup option: $1" ;;
    esac
    shift
  done
  [ -n "$base" ] || base=$(default_base)
  git -C "$primary" rev-parse --verify --quiet "$base^{commit}" >/dev/null || die "Base does not resolve to a commit: $base"

  while IFS= read -r path; do
    if [ ! -d "$path" ]; then
      printf 'Skip missing worktree: %s; run git worktree prune after inspection\n' "$path"
      continue
    fi
    path=$(cd "$path" && pwd -P)
    [ "$path" = "$primary" ] && continue
    branch=$(git -C "$path" branch --show-current)
    if [ -z "$branch" ]; then
      printf 'Skip detached worktree: %s\n' "$path"
      continue
    fi
    expected=$(target_for_branch "$branch")
    if [ "$path" != "$expected" ]; then
      printf 'Skip non-mirrored worktree: %s\n' "$path"
      continue
    fi
    if ! require_clean "$path"; then
      printf 'Skip dirty worktree: %s\n' "$path"
      continue
    fi
    if ! git -C "$primary" merge-base --is-ancestor "$branch" "$base"; then
      printf 'Keep unmerged branch %s at %s\n' "$branch" "$path"
      continue
    fi
    if [ "$apply" = false ]; then
      printf 'Would remove merged branch %s at %s\n' "$branch" "$path"
      continue
    fi
    git -C "$primary" worktree remove "$path"
    git -C "$primary" branch -d "$branch"
    printf 'Removed merged branch %s at %s\n' "$branch" "$path"
  done < <(git -C "$primary" worktree list --porcelain | sed -n 's/^worktree //p')
}

remove_branch() {
  local branch=$1 target
  target=$(target_for_branch "$branch")
  [ -d "$target" ] && registered_worktree "$target" || die "No registered mirrored worktree for $branch at $target"
  require_branch "$target" "$branch"
  require_clean "$target" || die "Worktree is dirty; not removing: $target"
  git -C "$primary" worktree remove "$target"
  printf 'Removed worktree for %s at %s; kept the branch.\n' "$branch" "$target"
}

install_post_merge_hook() {
  local hooks_dir hook script_path
  hooks_dir=$(git -C "$primary" rev-parse --git-path hooks)
  case "$hooks_dir" in
    /*) ;;
    *) hooks_dir="$primary/$hooks_dir" ;;
  esac
  hook="$hooks_dir/post-merge"
  if [ -e "$hook" ] && ! grep -q "$marker" "$hook"; then
    die "Refusing to overwrite existing hook: $hook"
  fi
  script_path=$(cd "$(dirname "$0")" && pwd -P)/$(basename "$0")
  mkdir -p "$hooks_dir"
  cat > "$hook" <<EOF
#!/bin/sh
# $marker
"$script_path" cleanup-merged --apply >/dev/null 2>&1 || true
EOF
  chmod +x "$hook"
  printf 'Installed %s\n' "$hook"
}

command=${1:-}
[ -n "$command" ] || { usage; exit 1; }
shift
setup_repository
case "$command" in
  sync) [ "$#" -eq 1 ] || { usage; exit 1; }; sync_branch "$1" ;;
  cleanup-merged) cleanup_merged "$@" ;;
  remove) [ "$#" -eq 1 ] || { usage; exit 1; }; remove_branch "$1" ;;
  install-post-merge-hook) [ "$#" -eq 0 ] || { usage; exit 1; }; install_post_merge_hook ;;
  *) usage; exit 1 ;;
esac
