#!/bin/sh
# Register the shared Cua Driver MCP server with local agent clients.
# The server definition lives in config/mcp/cua-driver.json. This script
# resolves the installed binary and writes each client's local config.
set -eu

SCRIPT_DIR=$(CDPATH= cd -- "$(dirname -- "$0")" && pwd)
ROOT=$(CDPATH= cd -- "$SCRIPT_DIR/.." && pwd)
SPEC="$ROOT/config/mcp/cua-driver.json"

if [ ! -f "$SPEC" ]; then
  echo "Missing MCP definition: $SPEC" >&2
  exit 1
fi

resolve_binary() {
  if command -v cua-driver >/dev/null 2>&1; then
    command -v cua-driver
    return
  fi

  for candidate in "$HOME/.local/bin/cua-driver" /usr/bin/cua-driver; do
    if [ -x "$candidate" ]; then
      printf '%s\n' "$candidate"
      return
    fi
  done

  echo "cua-driver is not installed. Install it, then rerun this script." >&2
  exit 1
}

BINARY=$(resolve_binary)
export CUA_DRIVER_BINARY="$BINARY"
export CUA_DRIVER_SPEC="$SPEC"

echo "Using $BINARY"

if command -v grok >/dev/null 2>&1; then
  grok mcp add --scope user cua-driver -- "$BINARY" mcp
  echo "Registered grok user MCP server cua-driver"
fi

if command -v codex >/dev/null 2>&1; then
  if codex mcp list 2>/dev/null | grep -q 'cua-driver'; then
    codex mcp remove cua-driver
  fi
  codex mcp add cua-driver -- "$BINARY" mcp
  echo "Registered codex MCP server cua-driver"
fi

if command -v claude >/dev/null 2>&1; then
  claude mcp remove --scope user cua-driver >/dev/null 2>&1 || true
  claude mcp add --scope user --transport stdio cua-driver -- "$BINARY" mcp
  echo "Registered claude user MCP server cua-driver"
fi

python3 - "$BINARY" <<'PY'
import json
import os
import sys
from pathlib import Path

binary = sys.argv[1]
home = Path.home()

def on_path(name):
    return any(
        os.access(Path(entry) / name, os.X_OK)
        for entry in os.environ.get("PATH", "").split(os.pathsep)
        if entry
    )

def load(path):
    if not path.exists():
        return None
    with path.open() as handle:
        return json.load(handle)

def save(path, data):
    path.parent.mkdir(parents=True, exist_ok=True)
    with path.open("w") as handle:
        json.dump(data, handle, indent=2)
        handle.write("\n")

opencode_path = home / ".config" / "opencode" / "opencode.json"
opencode = load(opencode_path)
if opencode is not None:
    servers = opencode.setdefault("mcp", {})
    servers["cua-driver"] = {
        "type": "local",
        "command": [binary, "mcp"],
        "enabled": True,
    }
    save(opencode_path, opencode)
    print(f"Registered opencode MCP server in {opencode_path}")

if (home / ".cursor").exists() or on_path("cursor"):
    cursor_path = home / ".cursor" / "mcp.json"
    cursor = load(cursor_path) or {"mcpServers": {}}
    cursor.setdefault("mcpServers", {})["cua-driver"] = {
        "command": binary,
        "args": ["mcp"],
        "type": "stdio",
    }
    save(cursor_path, cursor)
    print(f"Registered cursor MCP server in {cursor_path}")

if (home / ".gemini").exists() or on_path("gemini"):
    gemini_path = home / ".gemini" / "config" / "mcp_config.json"
    gemini = load(gemini_path) or {"mcpServers": {}}
    gemini.setdefault("mcpServers", {})["cua-driver"] = {
        "command": binary,
        "args": ["mcp"],
    }
    save(gemini_path, gemini)
    print(f"Registered gemini MCP server in {gemini_path}")
PY
