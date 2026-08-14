# OpenCode Harness

Use this adapter when the chat manager or background workers run through
OpenCode. Unlike a tool-rich Codex chat manager, an OpenCode chat session should
not be the only owner of a long-running queue. Start the bundled supervisor so
worker processes, logs, feedback cursors, and dependency gates survive chat
turns and restarts.

## Initial Setup

Ask the user to choose the primary chat provider/model and background worker
provider/model. Preserve existing provider availability. Do not add
`enabled_providers` or `disabled_providers` unless the user explicitly requests
provider restrictions.

The repository owns the worker definition and permissions in `opencode.json` or
`.opencode/opencode.json`. The supervisor does not inject
`OPENCODE_CONFIG_CONTENT`, because an inline final-scope configuration would
silently override the project's reviewed policy. See the readable starting
template at [worker policy template](../../assets/opencode-worker-policy.json).

Start from the template's deny-by-default policy, then grant only the project
commands the worker needs. Keep `worker` in `primary` mode so headless
`opencode run --agent worker` cannot fall back to a broader built-in agent. The
supervisor removes tracker, SSH-agent, and common cloud/application credentials
from the worker environment, but project configuration remains the authority
for tool permissions.

OpenCode permissions are tool controls, not an operating-system sandbox. An
allowed build command executes with the host user's process and network
authority. For hostile repositories, run OpenCode inside a container or OS
sandbox with a worktree-only mount, filtered provider egress, no host
credentials, and a disposable home. Do not claim equivalence with Codex's
native `workspace-write` sandbox.

Projects define the worker in `opencode.json` for direct use, model selection,
and permissions. Keep the primary chat model unchanged:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "agent": {
    "worker": {
      "description": "Implements one isolated orchestration work item.",
      "mode": "primary",
      "model": "<selected-worker-provider>/<selected-worker-model>"
    }
  }
}
```

Do not overwrite an existing project configuration. Merge the agent entry and
preserve permissions, MCPs, providers, commands, and the primary model.

## Supervisor

Build the standalone supervisor once. Go is the only build dependency; the
resulting binary has no runtime package dependency:

```sh
go build -o "$HOME/.local/bin/ticket-orchestrator" \
  /path/to/implement-tickets/scripts/orchestrate.go
```

Initialize a run from the target repository. The exact tracker options depend
on the selected work-source adapter:

```sh
ticket-orchestrator init \
  --repo "$PWD" \
  --source github \
  --harness opencode \
  --target master \
  --ready-label ready-for-agent \
  --concurrency 3 \
  --worker-agent worker
```

Add repeated `--issue`, `--issue-range`, or one `--issue-query` argument when
the run should own less than the full ready queue.

Run one cycle while testing configuration:

```sh
ticket-orchestrator supervise \
  --state /path/printed/by/init/state.json \
  --once
```

Then run the persistent loop under the user's process manager, terminal
multiplexer, or service supervisor. The orchestration script does not daemonize
itself or hide its PID.

Use `status`, `inbox`, `stop`, and `supervise --once` for operator interaction.
The runtime directory lives under the user's state directory by default, not in
the repository, so worker logs do not dirty Git.

## Worker Launch

The supervisor launches a dedicated worker session with raw JSON events:

```sh
opencode run \
  --agent worker \
  --format json \
  --title "orchestration-<item>" \
  --dir "$WORKTREE" \
  "$(<"$PROMPT_FILE")"
```

Pass `--model provider/model` only when the user selected an explicit worker
model that is not already defined by the worker agent. Never pass `--auto` to a
headless orchestration worker. The supervisor may accept other launcher
arguments after preflight.

## OpenCode 2 Migration

OpenCode 2 has a stronger ordered action/resource policy, rejects unanswered
permissions during non-interactive runs, canonicalizes edit paths, and rejects
symlink escapes. A future adapter should translate the policy to V2
`permissions` rules and use `shell` and `subagent` action names. V2 still states
that an allowed shell command has the host user's filesystem, process, and
network authority, so retain the same OS-isolation recommendation.

OpenCode workers receive the same offline publication contract as Codex
workers. They leave tested worktree diffs but do not write Git metadata, push,
query the tracker, create reviews, merge, or clean worktrees. The supervisor
owns those operations.

## Recovery

The supervisor wraps each OpenCode process with a durable exit-code file. On
restart it distinguishes a live worker from a completed or PID-reused process,
then resumes polling without launching a second owner. If a worker dies or its
log stalls, preserve its worktree/log and dispatch at most one recovery worker
that inspects existing state first.

Direct user feedback goes through the supervisor inbox. Safe requests wait for
the current worker to exit or checkpoint; cancel and safety requests may stop
the process group immediately.
