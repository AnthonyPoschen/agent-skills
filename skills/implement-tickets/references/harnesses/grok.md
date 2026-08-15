# Grok Harness

Use this adapter when the chat manager or isolated workers run through Grok
Build. Grok is a capable interactive harness: native background subagents,
worktree isolation, monitors, and a scheduler that can wake a new manager turn.
Prefer those primitives. Keep the supervisor CLI only for durable tracker
cursors, event history, and publication bookkeeping. Do not run
`ticket-orchestrator supervise` as the Grok worker launcher; its process
adapters launch Codex or OpenCode, not Grok subagents.

## Preferred Operating Model

Keep the manager loop in this Grok chat session. At start, confirm standing
authority to push branches, post review comments, and close leftover open
items after verified integration. The manager:

1. Selects the ready set from the work-source adapter.
2. Claims each item, prepares one isolated checkout, and launches one worker.
3. Publishes each successful worker independently.
4. Watches draft reviews for authorized feedback, CI failures, and human merges.
5. Closes an item only after verified integration if the tracker left it open.
6. Recomputes readiness and fills free concurrency slots without waiting for
   another user prompt.

Use the bundled supervisor when the user asks for unattended bookkeeping that
must survive a Grok process exit. `ticket-orchestrator init --harness grok`
is valid. `sync`, `status`, `events`, and `publish` work. `dispatch` and
`supervise` do not launch Grok workers; this chat session still owns
`spawn_subagent`. See the [capability contract](capabilities.md).

Do not encode the manager loop as a Grok workflow script. `parallel()` is a
completion barrier, workflows have no clock or poll primitive, and a process
restart makes a live run terminal. A workflow may fan out a one-shot research
or verification panel; it cannot own scheduling, PR watch, or merge follow-through.

## Role Split

The chat manager owns tracker I/O, Git publication, review comments, item
closure after integration, dependency gates, and user conversation.

A worker implements one item in its checkout and exits. It does not push,
create or edit reviews, merge, close items, or delete worktrees. Default
`general-purpose` workers inherit the parent's MCP servers, including GitHub.
The worker prompt must forbid tracker and review tools. When the project
defines a worker agent, use it and keep tracker MCP out of that agent. Copy
[Grok worker agent](../../assets/grok-worker-agent.md) to
`.grok/agents/ticket-worker.md` when the project does not already define one.

Grok subagents cannot spawn children. The manager is the only fan-out point.

## Worker Launch

Write an immutable prompt file outside the checkout. Create or reuse the
item's isolated checkout and branch first, then launch one background worker
per ready item up to the concurrency cap:

```text
spawn_subagent
  prompt: contents of the immutable prompt file, including the project's
          implementation skill invocation and the worker contract
  description: implement issue <number>
  subagent_type: ticket-worker when that agent exists, otherwise general-purpose
  background: true
  capability_mode: all
  cwd: the prepared checkout
```

`cwd` and `isolation: worktree` are mutually exclusive. Prefer a manager-prepared
checkout plus `cwd` so the worker starts on the issue branch and never receives
the primary worktree. Use `isolation: worktree` only when no checkout exists
yet; record the returned path as the item checkout and never apply that
worktree onto the primary tree.

Pass `model` only when the user or project selected a worker model. Allowed
explicit slugs are the models this Grok session may spawn. Omit `model` to
inherit the manager session.

Write run state outside the repository before the first launch. Default to
`$HOME/tmp/implement-tickets/<repo>-<run-id>/` so checkouts stay off tmpfs
`/tmp` and do not occupy the shared `$HOME/tmp` root. Do not put run-owned
clones in `/tmp`. Keep other tools free to use `$HOME/tmp` beside this
directory. If `ticket-orchestrator init` prints a state directory, use that
instead. Record item number, branch, checkout path, subagent ID, prompt
path, handoff path, launch time, and log age before launching the next worker.
One worker owns a branch at a time. The scheduled poll prompt must include
that state path.

## In-Session Reaping

`get_command_or_subagent_output` with several IDs and a positive timeout waits
for every listed task. That is a publication batch. Do not use it to wait out
a concurrency group.

After each launch or wake:

1. Snapshot every active worker with `timeout_ms` omitted or `0`.
2. Start verification and publication for each successful exit immediately.
3. If workers remain, wait on one still-running ID with a short bounded
   timeout, then snapshot all again.
4. Answer interrupting user messages, then resume the same run.

A zero-exit worker without a handoff, proposed commit subject, and scoped
checkout diff is not implementation success. Preserve that checkout.

Replace a dead or stalled worker with `kill_command_or_subagent` on its
subagent ID, then launch one replacement that inspects the existing checkout
first. The checkpoint commit uses `Assisted-by: Grok/<model>` when the model
is known, otherwise `Assisted-by: Grok/AI`.

## Review Watch

After the first draft review exists, start both of these and keep them until
the run is terminal:

- A `persistent` `monitor` that prints one line per important review event
  (`FEEDBACK`, `CHECKS_FAILED`, `MERGED`, `CLOSED`) and stays quiet otherwise.
  Watch every in-review pull or merge request recorded in run state, not only
  the first published review. Filter with `grep --line-buffered`. Do not
  stream raw `gh` or MCP payloads.
- A `scheduler_create` poll, default interval `2m`, `foreground` true so the
  fire continues this manager conversation. Set `durable` only when the user
  asked for unattended continuation across sessions. The prompt must stay
  short: name the run-state path and say "run one implement-tickets poll
  cycle from that state." The wake steps in this section are part of that
  cycle. Read the core-workflow feedback policy when a comment needs
  classification.

The monitor gives low-latency wakes. The scheduler is the autonomy backbone so
a finished worker, a new review comment, or a human merge does not wait for
the user to prompt again. Delete both when the run becomes terminal.

A scheduled poll is a full wake.

On each wake:

1. Poll every in-review item for comments, unresolved threads,
   requested-change reviews, reactions, and checks.
2. Apply the core-workflow authorization and classification policy.
3. Take the matching tracker action in this turn: post a ticket reply,
   launch a follow-up worker, push a completed conflict repair, or close
   an item after verified integration.
4. Queue newly authorized IDs when a worker already owns that branch.
   Snapshot the IDs given to a running follow-up so later IDs become a
   later worker.

Do not answer a ticket question only in chat. Do not wait for another user
prompt to reply or to implement unambiguous authorized feedback.

If a required manager write is denied (push, review comment, or close after
verified integration), report the exact pending action in that turn. Do not
treat the denial as an unchanged poll.

An unchanged wake emits no user-visible message. After five minutes without an
important event, emit one compact heartbeat. Treat HTTP 402 or an equivalent
usage-balance error as a stall: preserve the checkout and resume once after
the user says credits are restored.

For a follow-up on the same item, prefer `resume_from` with the completed
worker's ID, or spawn a new worker with `cwd` set to the same checkout. On
resume, send only the new evidence. Include the original title and body only
on a cold start.

## Integration, Closure, And Unblocked Work

A human merges. The manager never merges, never enables auto-merge, and never
approves its own review.

When a watch event or poll reports `MERGED` or an equivalent integration:

1. Fetch the target and confirm the merge commit is an ancestor.
2. Perform the work-source close step if the item is still open.
3. Mark the item integrated, reclaim the concurrency slot, and recompute the
   ready set from current tracker state.
4. Claim and launch the next ready unblocked items up to the cap.
5. Do not wait for remaining siblings, a wave, or another user message.

A human-closed unmerged review is `review_closed`, not integration. Do not
close its work item as if the work landed.

## Session Limits

A final manager answer cannot receive later tool output in that same turn.
Scheduler and monitor wakes start new turns; do not promise that a detached
process will append to an ended turn.

On reconnect, compaction, or a scheduled fire, load run state, snapshot
workers, replay missed important events, then continue. If no durable
scheduler exists and the host cannot keep a monitoring turn, say that
limitation; do not claim unattended PR follow-through.

After `run_completed`, stop workers, delete the run's monitors and scheduled
polls, remove verified run-owned checkouts, retain audit logs outside those
checkouts, and return the compact completion table.
