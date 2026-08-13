# Core Workflow

## Supervisor

Run orchestration through a persistent local supervisor. It must persist current
state, a human-readable status table, an inbox, an append-only event stream, and
per-worker logs. Rotate old event logs by size or age.

The supervisor must emit events for dispatch, checkpoint, review creation,
feedback, rebase, CI change, merge, failure, stall, inbox acceptance, and inbox
application. When no event occurs for five minutes, emit one compact all-worker
status table with process health, log activity, working-tree summary, diff line
counts, latest commit, and review state.

The chat manager does not directly modify a worker branch while the supervisor
runs. It writes requests to the inbox. For each request, report accepted and
applied or queued events. Safe changes wait for a checkpoint. Cancel, urgent,
or safety requests may stop a worker immediately.

## Worker Lifecycle

1. Create the isolated worktree and deterministic source branch.
2. Launch one worker using the selected worker adapter and standard prompt.
3. Require an early checkpoint commit, explicit branch push, and review creation.
4. Verify the review source, target, body, checks, and feedback independently.
5. End the worker after the checkpoint. Relaunch only for a defined follow-up.
6. Preserve the worktree and inspect logs/diff before replacing a dead or stalled
   worker.

## Feedback And Rebase

Poll work items and reviews for new feedback. Classify feedback as actionable,
question, approval, or informational. Dispatch one owner for actionable items.
After a target-branch update, request a checkpoint then rebase every active
branch. Mechanical conflicts may be handled automatically. Pause for decisions
that need human intent.

## Live Validation

When the supervisor or adapter is new, treat it as experimental. Preserve state,
worker logs, and event evidence for every fault. Fix the implementation or skill
before relying on the failed behavior again. Do not add a separate test suite
unless the user requests one.
