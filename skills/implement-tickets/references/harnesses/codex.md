# Codex Harness

Use this adapter when the chat manager or isolated workers run through Codex.

## Preferred Operating Model

Codex can keep the manager loop in the active chat session: recompute tracker
state, launch multiple isolated `codex exec` workers, monitor JSONL sessions,
verify commits, publish draft PRs, and apply review feedback. Use this direct
mode when the user is present and the chat runtime supplies durable tool calls.

Use the bundled supervisor when the user asks for unattended operation, the chat
session cannot persist worker ownership, or Codex workers are launched by
another agent harness.

Codex should not duplicate features it already provides well. Prefer native
subagents, sandboxing, reasoning, and commentary. Retain the CLI for durable
tracker cursors, feedback deduplication, dependency gates, integration proof,
event history, recovery, and publication authority. See the
[capability contract](capabilities.md).

## Worker Launch

Create an immutable prompt file and a unique run directory outside the
worktree. Launch one process per issue:

```sh
codex exec \
  --json \
  --sandbox workspace-write \
  -c 'approval_policy="never"' \
  --ephemeral \
  -C "$WORKTREE" \
  -o "$RUN_DIR/last-message.txt" \
  - < "$RUN_DIR/prompt.md" > "$RUN_DIR/events.jsonl" 2>&1
```

The bundled supervisor uses self-contained local clones rather than shared Git
worktrees. A linked worktree stores its writable index, refs, and object data in
the primary repository's `.git` directory, which sits outside Codex's
`workspace-write` boundary. Do not solve that by passing the whole shared
`.git` directory through `--add-dir`; that would let one worker mutate other
branches and checkouts.

Use the installed CLI's current help if option placement differs. Do not pass a
model override unless the user or project explicitly selected one. If selected,
record it in run state and pass the supported Codex model/config option.

The process receives write access only to its issue worktree (and explicitly
named matching cross-repository worktrees). General network access remains off;
the manager owns GitHub and publication.

## Monitoring

Track issue, branch, worktree, PID plus start time, JSONL path, last-message
path, launch time, latest event, and log age. Read new JSONL incrementally and
surface progress at least once per minute while a worker runs.

React to each worker's successful exit immediately. Verify and publish that
worker's draft review without waiting for other active Codex sessions. Keep the
sibling sessions running and continue monitoring them during the completed
item's verification and publication. A bounded concurrency group controls
worker load; it is not a publication batch.

An autonomous background supervisor does not by itself deliver messages into a
completed Codex chat turn. While the run is nonterminal, keep the manager turn
open with bounded event waits and send important events through commentary as
they arrive. Do not send a final answer merely because the run is waiting for a
review, feedback, checks, or a human merge. New user messages interrupt the
wait; answer them, then continue monitoring the same run.

Track the last event sequence actually shown to the user. After reconnect,
compaction, or interruption, run `events --after <sequence>` and surface missed
important events before reporting current status. Advance the relay cursor only
after the event has been sent in commentary. If the host cannot keep an active
monitoring turn, state that limitation explicitly; do not promise proactive
chat updates from a detached process.

On successful exit, require a final handoff, a completion event, a scoped local
commit, and a clean or intentionally staged worktree. A zero exit without a
commit is not implementation success. Preserve failed run directories and
worktrees.

Require the final handoff to propose one concise Conventional Commit subject
from the actual diff. The supervisor validates that subject and appends the
assistance trailer when it creates the checkpoint. Generic issue-number
subjects are not acceptable.

## Follow-Up Workers

Reuse the same branch/worktree for review feedback, CI repair, or target-branch
reconciliation. The follow-up prompt includes the original issue title/body
verbatim, exact new feedback or failure evidence, current PR/head context, and
the same publication restrictions. Only one Codex worker may own the worktree
at a time.

Poll comments, reactions, unresolved inline threads, requested-change reviews,
and checks for every waiting review without requiring another user prompt.
Apply the authorization and classification policy in the core workflow before
constructing a follow-up prompt. In the same poll, reply on the ticket to
authorized questions and launch a follow-up worker for unambiguous
implementation requests. Keep polling while a worker runs. Snapshot the
authorized feedback IDs given to that worker; newly authorized IDs belong to
a later follow-up. Surface important supervisor events under an `Events`
lead-in and keep unchanged polls silent. If a required manager write is
denied, report the exact pending action immediately. Treat HTTP 402 as a
stall and resume once after credits are restored.

After `run_completed`, relay the final event and return a final answer. The
supervisor must already have stopped workers, removed verified run-owned
checkouts, retained audit artifacts outside them, released its lifetime lock,
and exited.
