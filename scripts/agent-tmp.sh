# shellcheck shell=bash
# Source me: gives this agent a private TMPDIR so fixed /tmp names cannot collide
# between agents sharing one box.  Agent name: $AGENT_NAME, else $GROK_AGENT_ID, else $USER-$$.
#   . scripts/agent-tmp.sh
#   body="$(pr_body_file)"; printf '%s\n' "$text" >"$body"; gh pr create --draft --body-file "$body"
agent_tmp_init() {
  local name="${AGENT_NAME:-${GROK_AGENT_ID:-${USER:-agent}-$$}}"
  name="${name//[^A-Za-z0-9._-]/_}"
  local root="${AGENT_TMP_ROOT:-/tmp/agents}"
  mkdir -p "$root" && chmod 1777 "$root" 2>/dev/null || true
  mkdir -p "$root/$name" && chmod 700 "$root/$name" || return 1
  export AGENT_TMPDIR="$root/$name" TMPDIR="$root/$name"
}
agent_mktemp() { mktemp "${AGENT_TMPDIR:?run agent_tmp_init}/${1:-tmp}.XXXXXXXX"; }
agent_mktemp_d() { mktemp -d "${AGENT_TMPDIR:?run agent_tmp_init}/${1:-dir}.XXXXXXXX"; }
pr_body_file() { agent_mktemp prbody; }
agent_tmp_init
