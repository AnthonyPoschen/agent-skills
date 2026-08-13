---
name: implement-all-tickets
description: Orchestrate a dependency-linked backlog of work items through isolated workers, worktrees, review feedback, and human-controlled merges. Use when the user asks to implement all tickets, execute a backlog in parallel, coordinate ticket agents, monitor worker progress, or process review feedback across GitLab, GitHub, or local work-item files.
---

# Implement All Tickets

Use this skill to deliver a dependency-linked work-item backlog safely. It is
an orchestration workflow, not a shortcut around the project implementation
workflow: every worker implements one vertical slice only.

Read `references/core-workflow.md` before scheduling work. Select and read one
work-source adapter from `references/work-sources/` and one harness adapter from
`references/harnesses/` before the first tracker or worker action. If no adapter
fits the environment, inspect the repository and write the missing adapter
contract before starting work.

## Start Safely

1. Read project instructions, domain documentation, applicable decisions, and
   the project work-item conventions.
2. Inspect the target branch, primary worktree, existing worktrees, active
   workers, open reviews, and candidate work items.
3. Refuse to start when the primary worktree has unrelated uncommitted changes.
   Never reset, clean, or overwrite another person's work.
4. Confirm workers may push branches and create reviews, but only a human may
   merge.
5. Select an isolated worktree location and optional concurrency cap. Without a
   cap, schedule every ready item concurrently.
6. Select the provider/model for the chat manager and the provider/model for
   background workers. Preserve existing harness provider settings. Do not add
   provider restrictions unless the user explicitly requests them.
7. Start the configured supervisor before launching workers.

## Adapter Contract

The selected work-source adapter must provide a way to:

- List open work items and read their full body and comments.
- Identify explicit blockers, review feedback, and acceptance criteria.
- Read and update work-item and review state.
- Discover review source/target branches and integration status.
- Identify whether a closed work item is actually integrated into the target
  branch.

The selected harness adapter must provide a way to:

- Launch one isolated worker with a known model and prompt template.
- Record the worker task/session ID, PID when available, worktree, branch, and
  unique stdout/stderr log.
- Read worker health and request a checkpoint, stop, or replacement.
- Keep only one worker responsible for a branch at a time.

Record the selected chat-manager and worker provider/model in the supervisor
state. The chat manager can use a higher-capability model while workers use the
user-selected background model. Do not assume a provider, model, or cost tier.

## Manager Rules

- Build the dependency graph only from explicit blocker references. Do not infer
  dependencies from ticket order.
- A work item is ready only when every declared blocker is integrated into the
  target branch.
- Use a deterministic branch and one isolated worktree per work item.
- Create worker branches without an upstream. Workers push only their assigned
  source branch through an explicit refspec.
- The supervisor owns routine polling, dependency scheduling, worker health,
  rebase requests, and unambiguous repair dispatch. The chat manager stays
  responsive and writes direct user requests to the supervisor inbox.
- A worker creates an early tested checkpoint and review, then exits. Restart it
  only for feedback, CI, rebase, or an inbox request.
- Rebase active branches after target-branch changes. Use force-with-lease only.
  Pause for conflicts that need product or architecture judgment.
- Detect dead workers and stalled logs. Preserve their worktree and diff before
  one replacement worker validates the existing state and continues.
- Treat review comments as work inputs. Auto-handle only specific, in-scope,
  unambiguous feedback. Pause for scope changes or conflicting direction.
- Every agent-written review comment begins with `AI-generated:`. Use concise
  ASD-STE100-style technical English: changed location first, then short action
  bullets. Do not write commit-message paragraphs or routine test lists.
- Never merge, close a work item manually, or mark feedback resolved without
  verification.

## Completion

Finish only when every selected work item is merged, explicitly skipped by the
user, or reported blocked or failed. Return a compact table with item, state,
review, verification, blockers, and outstanding feedback. Retain run-owned
worktrees until their reviews merge and the user authorizes cleanup.
