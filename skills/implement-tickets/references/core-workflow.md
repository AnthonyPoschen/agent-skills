# Core Workflow

## Sources Of Truth

Use the tracker, Git, and the target branch as orchestration state. The durable
supervisor cache exists to resume work, not to replace those sources:

- The tracker owns work-item state, explicit dependencies, assignments, review
  feedback, and review integration state.
- Git owns branches, commits, worktrees, and proof that work reached the target
  branch.
- The supervisor owns worker PIDs, prompts, logs, cursors, inbox requests, and
  the last observed external state.

Stable project-owned settings live in `.github/implement-tickets.json`, with
their meaning explained through the documentation chain rooted at `AGENTS.md`.
See [project setup and discovery](project-setup.md). Runtime state must record
the resolved values but must not silently replace missing project conventions.

Re-read external state after every merge, closure, assignment change, worker
failure, review update, and target-branch change. Never keep scheduling from a
stale ready-set snapshot.

## Readiness

A work item is ready only when it is open, carries the project's ready label or
state, is unassigned (unless resuming a run-owned assignment), and every
explicit blocker is integrated into the target branch. A draft review, a pushed
branch, or a closed blocker without integrated code does not satisfy a blocker.

Build dependencies only from native tracker relationships or an explicit
`Blocked by` section. Do not infer them from item numbers, roadmap order, shared
files, or likely architecture.

## Supervisor State

For unattended or OpenCode orchestration, run the bundled supervisor. It keeps
state outside the repository by default and persists:

- `state.json`: configuration, target head, item states, review cursors, worker
  ownership, and verification results;
- `events.jsonl`: append-only dispatch, checkpoint, publication, feedback, CI,
  merge, failure, stall, inbox, and recovery events;
- `inbox.jsonl`: direct user requests with accepted, queued, applied, or failed
  state;
- `workers/`: immutable prompts, JSONL/stdout logs, last messages, and exit-code
  files.

Use an exclusive run lock so two supervisors cannot dispatch the same branch.
Hold supervisor ownership for the process lifetime; a per-cycle lock alone lets
old and replacement supervisors alternate between polls.
Record PID start time as well as PID to avoid mistaking a reused PID for a live
worker. Rotate large event logs without deleting the newest evidence.

When no event occurs for five minutes, emit one compact status table with
worker health, log age, worktree state, diff size, latest commit, review state,
checks, and blockers. The chat manager stays responsive by writing requests to
the inbox instead of modifying active worker branches.

## Isolated Checkout Lifecycle

1. Fetch the latest target branch without changing the primary checkout.
2. Derive a stable `issue/<number>-<slug>` or equivalent branch.
3. Reuse a matching run-owned branch/checkout after inspecting its commits,
   dirty state, review, and worker ownership. Refuse collisions.
4. Create a self-contained local clone and a new branch from the fetched target
   with no upstream. Keep its `.git` directory inside the isolated checkout so
   Codex and OpenCode can commit without write access to shared repository
   metadata. Do not use `git worktree` plus a broad writable grant to the
   primary `.git` directory.
5. Assign or claim the work item before launching its worker.
6. Retain each checkout until the entire selected run is terminal. During run
   finalization, remove only verified clean run-owned checkouts whose heads
   match their published review heads.

Never delete a dirty, unpublished, failed, or conflict-paused checkout.

## Worker Contract

Start every worker statement with the project's explicit implementation skill
invocation when one exists. Include the item title and body verbatim, followed
by resolved branch/worktree paths and these rules:

- implement one work item and its smallest supporting changes;
- read every applicable project instruction file;
- treat the fetched target branch as the integration source of truth;
- preserve unrelated changes and already merged behavior;
- run focused checks regularly and the full applicable suite at the end;
- review the final diff and leave tested changes unstaged for the supervisor;
- do not write Git metadata or create commits;
- do not query the tracker, push, create or edit reviews, merge, close items, or
  delete worktrees;
- stop with a precise report when authority or product judgment is required.
- end a successful handoff with one `Commit subject: <type>: <summary>` line
  derived from the actual diff.

Keeping workers offline and publication-free makes retries safe. The manager can
independently reject an unrelated diff without undoing external writes.

## Completion Verification And Publication

Worker completion is event-driven per item, not a batch barrier. After any one
worker exits successfully, immediately begin the following pipeline for that
item while other workers continue running. Never delay verification or draft
review publication until sibling workers, the current concurrency wave, or the
entire ready set finishes. If several workers finish together, process their
pipelines independently in completion order; use safe parallel verification
when resource limits allow.

For each completed worker:

1. Read its final handoff and machine-readable completion event.
2. Confirm the worktree has no unresolved conflicts or unrelated changes.
3. Confirm the configured target is still an ancestor or has been intentionally
   reconciled and the diff contains only the assigned work.
4. Run independent focused checks and the project-wide suite in proportion to
   risk. Validate generated assets, secrets, migrations, and cross-repository
   contracts when present.
5. Validate the implementing AI's proposed Conventional Commit subject, stage
   the verified diff, and create a supervisor-owned checkpoint commit with the
   required `Assisted-by` trailer. For OpenCode only, when a successful worker
   omits its subject, the supervisor may use the item's already-valid
   Conventional Commit title. Record a `commit_subject_derived` event. Never
   fall back to an issue-number subject.
6. Push only the expected branch with an explicit refspec.
7. Create or update one draft review. Verify its source, target, head commit,
   body, and checks after every write.

The review body records the canonical item, summary, verification, affected
repositories, assumptions, risks, and convergence work left for later items.
Never advertise work as ready while required checks fail.

## Feedback, CI, And Integration Drift

Poll review comments, review threads, requested changes, and checks. Ignore bot
noise and feedback already addressed by a later commit. Classify new input as:

- actionable and unambiguous: relaunch one follow-up worker in the same
  worktree with the exact feedback and original item body;
- question: dispatch a response worker when repository evidence can answer it;
  otherwise bring it to the chat manager;
- approval or informational: record it without dispatch;
- scope or product change: pause and ask the user.

Do not resolve human review threads for the reviewer.

After a feedback worker completes, reply in the original review thread. Begin
the reply with `AI-generated:`. State the answer first, then name the published
commit when one exists. A feedback-only answer still requires a thread reply.
The supervisor posts this reply before marking the feedback applied; a reply
failure leaves the work item actionable.

After the target branch changes, do not publish a stale branch. A follow-up
worker may reconcile the latest target while preserving both reviewed behavior
sets. A normal merge is safer for already reviewed parallel work; a rebase is
appropriate when linear history is required. Never make a substantive conflict
choice silently.

For CI failures, collect the failing check names and logs in the manager, then
give the offline follow-up worker the exact failure evidence. Do not grant the
worker broad tracker access merely to inspect CI.

## Failures And Recovery

Treat a missing process, nonzero exit, stale log, malformed worker output, or
unexpected branch mutation as a recoverable run state. Preserve prompts, logs,
worktrees, commits, and dirty files. Before replacement, inspect that evidence
and launch at most one worker that validates existing state before continuing.

Do not label hard work as blocked. Pause only for missing authority, unavailable
infrastructure, secrets, a repeated environmental limitation, or a decision
that would change product intent.

## Completion

An integrated item is complete only after all required reviews merge and the
resulting commit is present on the target branch. A human-closed unmerged review
is a distinct terminal `review_closed` outcome, not successful integration.
Recompute the ready set immediately.

When every selected item is terminal, stop and reap workers before cleanup.
Verify each managed checkout is self-contained, below the run-owned worktree
root, clean, and at the published review head. Remove those checkouts and the
empty worktree root. Keep state, event logs, prompts, handoffs, and worker logs
in the separate run directory. Emit `run_completed`, release lifetime
ownership, and exit. Never delete a remote branch as implicit cleanup.
