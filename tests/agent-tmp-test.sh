#!/usr/bin/env bash
set -euo pipefail
H="$(cd "$(dirname "$0")/.." && pwd)/scripts/agent-tmp.sh"
R="$(mktemp -d)"; trap 'rm -rf "$R"' EXIT; export AGENT_TMP_ROOT="$R"
a="$(AGENT_NAME=coder bash -c ". '$H'; echo \$TMPDIR")"
b="$(AGENT_NAME=chief bash -c ". '$H'; echo \$TMPDIR")"
[[ "$a" != "$b" ]] || { echo "FAIL same dir"; exit 1; }
[[ "$(stat -c %a "$a")" == 700 ]] || { echo "FAIL mode"; exit 1; }
fa="$(AGENT_NAME=coder bash -c ". '$H'; pr_body_file")"; fb="$(AGENT_NAME=chief bash -c ". '$H'; pr_body_file")"
[[ "$fa" == "$a"/* && "$fb" == "$b"/* && "$fa" != "$fb" ]] || { echo "FAIL pr body paths"; exit 1; }
files="$(AGENT_NAME=coder bash -c ". '$H'; for i in \$(seq 200); do agent_mktemp x & done; wait")"
[[ "$(sort -u <<<"$files" | wc -l)" == 200 ]] || { echo "FAIL collision"; exit 1; }
AGENT_NAME='../evil x' bash -c ". '$H'; [[ \$TMPDIR == '$R'/.._evil_x ]]" || { echo "FAIL sanitize"; exit 1; }
echo "agent-tmp-test: ok"
