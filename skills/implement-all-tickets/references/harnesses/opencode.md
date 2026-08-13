# OpenCode Harness

Use this adapter when background workers run through OpenCode.

## Initial Setup

Ask the user to choose the primary chat provider/model and background worker
provider/model. Preserve existing provider availability. Do not add
`enabled_providers` or `disabled_providers` settings unless the user explicitly
requests provider restrictions.

Create a project `opencode.json` that keeps the primary chat model unchanged and
defines a dedicated background worker agent. Substitute the selected worker
model:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "agent": {
    "worker": {
      "description": "Implements isolated work items for orchestration.",
      "mode": "subagent",
      "model": "<selected-worker-provider>/<selected-worker-model>"
    }
  }
}
```

## Worker Launch

Launch every background session with the dedicated worker agent:

```sh
opencode run --agent worker --auto "<worker prompt>"
```

The primary chat manager may use its selected primary model. Record both model
selections in supervisor state and status output. Do not use `general`, `build`,
or an ad-hoc model override for background work.

## Logs And State

Store supervisor state in the primary worktree:

```text
.opencode/orchestration/<run-id>/
```

Store each worker log inside its worktree:

```text
.opencode/orchestration/worker-<iid>-<task-id>.log
```

Ignore runtime directories in Git. They are removed with authorized worktree
cleanup.
