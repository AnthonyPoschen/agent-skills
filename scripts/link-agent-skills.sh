#!/bin/sh
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
ROOT=$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)
SKILLS_SRC=${AGENT_SKILLS_SRC:-"$ROOT/skills"}
SKILLS_DEST=${AGENT_SKILLS_DIR:-"$HOME/.agents/skills"}

if [ ! -d "$SKILLS_SRC" ]; then
  echo "Missing skills source directory: $SKILLS_SRC" >&2
  exit 1
fi

mkdir -p "$SKILLS_DEST"

linked=0
skipped=0
pruned=0

# A merged or renamed skill leaves a link to a folder that no longer exists.
# Remove only dangling links into this source, never anything else.
for target in "$SKILLS_DEST"/*; do
  [ -L "$target" ] || continue
  current_target=$(readlink "$target")
  case "$current_target" in
    "$SKILLS_SRC"/*)
      if [ ! -e "$current_target" ]; then
        rm "$target"
        pruned=$((pruned + 1))
      fi
      ;;
  esac
done
for skill_dir in "$SKILLS_SRC"/*; do
  [ -d "$skill_dir" ] || continue

  skill_name=$(basename "$skill_dir")
  target="$SKILLS_DEST/$skill_name"

  if [ -L "$target" ]; then
    current_target=$(readlink "$target")
    if [ "$current_target" != "$skill_dir" ]; then
      ln -sfn "$skill_dir" "$target"
      linked=$((linked + 1))
    fi
    continue
  fi

  if [ -e "$target" ]; then
    echo "Skipping existing non-symlink skill: $target" >&2
    skipped=$((skipped + 1))
    continue
  fi

  ln -s "$skill_dir" "$target"
  linked=$((linked + 1))
done

echo "Linked $linked skill(s) into $SKILLS_DEST. Pruned $pruned stale link(s)."
if [ "$skipped" -gt 0 ]; then
  echo "Skipped $skipped existing non-symlink skill(s)." >&2
fi
