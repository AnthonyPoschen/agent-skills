---
name: implement-tickets
description: Orchestrate a selected set, range, query, or dependency-linked backlog of work items through isolated workers, worktrees, review feedback, and human-controlled merges. Use when the user asks to implement tickets, implement all tickets, run selected issues in parallel, coordinate ticket agents, monitor worker progress, or process review feedback across GitLab, GitHub, or local work-item files.
---

# Implement Tickets

Use this skill to deliver a selected dependency-linked work-item queue safely. It is
an orchestration workflow, not a shortcut around the project implementation
workflow: every worker implements one vertical slice only.

Read [project setup and discovery](references/project-setup.md), then read
[the core workflow](references/core-workflow.md) before scheduling work.
Select and read one work-source adapter and one harness adapter before the first
tracker or worker action:

- GitHub: [GitHub work source](references/work-sources/github.md)
- GitLab: [GitLab work source](references/work-sources/gitlab.md)
- Jira beta: [Jira work source](references/work-sources/jira.md)
- Codex: [Codex harness](references/harnesses/codex.md)
- OpenCode: [OpenCode harness](references/harnesses/opencode.md)

For harnesses that already provide agents, sandboxes, or progress UI, read the
[capability contract](references/harnesses/capabilities.md). Use only the
supervisor capabilities that add durable or deterministic behavior.

Codex can run the manager loop directly and use `codex exec` as the worker
isolation boundary. OpenCode needs the bundled Go supervisor to provide durable
state and process management. Build `scripts/orchestrate.go` for OpenCode and
for any environment where the chat session cannot safely own a long-running
loop. The supervisor uses only the Go standard library.
If no adapter fits, inspect the repository and write the missing adapter
contract before starting work.

## Start Safely

1. Read project instructions, domain documentation, applicable decisions, and
   the project work-item conventions.
2. Follow the project's `AGENTS.md` reference chain. Load
   `.github/implement-tickets.json` when present. If required settings remain
   ambiguous, stop preflight and ask the bundled setup questions together.
3. Inspect the target branch, primary checkout, existing isolated checkouts, active
   workers, open reviews, and candidate work items.
4. Refuse to start when the primary worktree has unrelated uncommitted changes.
   Never reset, clean, or overwrite another person's work.
5. Confirm the manager may push branches and create draft reviews, but only a human may
   merge.
6. Select an isolated checkout location and concurrency cap. Default to three;
   a smaller ready set naturally lowers concurrency. Do not launch an unbounded
   backlog merely because many items appear ready.
7. Select the provider/model for the chat manager and the provider/model for
   background workers. Preserve existing harness provider settings. Do not add
   provider restrictions unless the user explicitly requests them.
8. For OpenCode, initialize and start the bundled supervisor before launching
   workers. For Codex, either run the manager loop directly or use the same
   supervisor when durable unattended operation is requested.

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
- Use a deterministic branch and one isolated checkout per work item.
- Create worker branches without an upstream. Only the manager pushes the
  assigned source branch through an explicit refspec.
- The supervisor owns routine polling, dependency scheduling, worker health,
  rebase requests, and unambiguous repair dispatch. The chat manager stays
  responsive and writes direct user requests to the supervisor inbox.
- A worker creates and reviews a tested checkout diff, then exits. The manager
  or supervisor independently verifies it, creates the checkpoint commit,
  pushes through an explicit refspec, and
  creates or updates the draft review. Keeping publication outside the worker
  makes retries idempotent and keeps tracker/network authority centralized.
- The implementing AI proposes a concise Conventional Commit subject from the
  actual diff in its final handoff. The supervisor validates and uses it with
  the required assistance trailer. Never use a fixed issue-number subject.
- Update active branches after target-branch changes before publication. Prefer
  a normal merge when preserving independently reviewed work; use rebase only
  when the project requires linear history. Use force-with-lease only after a
  verified rebase. Pause for conflicts that need product or architecture
  judgment.
- Detect dead workers and stalled logs. Preserve their checkout and diff before
  one replacement worker validates the existing state and continues.
- Treat review comments as work inputs. Auto-handle only specific, in-scope,
  unambiguous feedback. Pause for scope changes or conflicting direction.
- Poll feedback while reviews wait for a human merge. A user must not need to
  prompt the manager after leaving a review comment. Feedback discovered while
  a worker owns the branch is queued for a later follow-up worker.
- In Codex, keep the manager turn open while the run is nonterminal and relay
  important supervisor events through commentary. Detached execution cannot
  write into a chat turn after a final answer; never promise that it can.
- Every agent-written review comment begins with `AI-generated:`. Use concise
  ASD-STE100-style technical English: changed location first, then short action
  bullets. Do not write commit-message paragraphs or routine test lists.
- Never merge, enable auto-merge, close a work item manually, approve the
  manager's own review, or mark human feedback resolved on the reviewer's
  behalf.

## Completion

Finish only when every selected work item is merged, explicitly skipped by the
user, closed by a human without merge, or reported blocked or failed. Once the
whole run is terminal, stop workers, remove verified clean run-owned checkouts,
retain audit logs outside those checkouts, emit `run_completed`, and exit.
Return a compact table with item, state, review, verification, blockers, and
outstanding feedback. Do not delete remote branches without separate authority.
