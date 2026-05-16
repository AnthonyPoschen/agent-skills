#!/bin/sh
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
ROOT=$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)
GIT_DIR=$(git -C "$ROOT" rev-parse --git-dir)

case "$GIT_DIR" in
  /*) ;;
  *) GIT_DIR="$ROOT/$GIT_DIR" ;;
esac

HOOKS_DIR="$GIT_DIR/hooks"
MARKER="agent-skills-link-hook"

install_hook() {
  hook_name=$1
  hook_path="$HOOKS_DIR/$hook_name"

  if [ -e "$hook_path" ] && ! grep -q "$MARKER" "$hook_path"; then
    echo "Refusing to overwrite existing hook: $hook_path" >&2
    echo "Add scripts/link-agent-skills.sh to it manually or move it aside." >&2
    exit 1
  fi

  cat > "$hook_path" <<'HOOK'
#!/bin/sh
# agent-skills-link-hook
REPO_ROOT=$(git rev-parse --show-toplevel 2>/dev/null) || exit 0
"$REPO_ROOT/scripts/link-agent-skills.sh"
HOOK

  chmod +x "$hook_path"
  echo "Installed $hook_path"
}

mkdir -p "$HOOKS_DIR"
install_hook post-merge
install_hook post-rewrite
