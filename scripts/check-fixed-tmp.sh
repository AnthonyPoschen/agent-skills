#!/usr/bin/env bash
# Fail on literal fixed /tmp/<name> paths in scripts and skills (shared-box collision hazard).
# Allowed: mktemp, $TMPDIR, /tmp/agents/ (per-agent root), or a line ending in "# tmp-ok".
set -euo pipefail
cd "${1:-$(dirname "$0")/..}"
hits="$(grep -rnE '/tmp/[A-Za-z]' --include='*.sh' --include='*.md' --include='*.py' scripts skills tests 2>/dev/null \
  | grep -vE '/tmp/agents|# tmp-ok|check-fixed-tmp' || true)"
if [[ -n "$hits" ]]; then echo "fixed /tmp paths found (use mktemp / agent_mktemp):"; echo "$hits"; exit 1; fi
echo "check-fixed-tmp: ok"
