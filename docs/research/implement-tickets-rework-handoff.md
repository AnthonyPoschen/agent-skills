# Implement tickets rework handoff

## Status

This is a planning handoff only. Do not treat it as approval to edit or delete
`skills/implement-tickets/`.

The repository was clean and pushed through `613d305f` before this discussion.
The notes created during this discussion are intentionally uncommitted so the
next session can revise, commit, or discard them as one decision.

## The problem to solve

The current `implement-tickets` skill is meant to coordinate selected tracker
items through isolated workers, review, feedback, and human merges. It has
grown into a provider-specific orchestration product instead of a focused skill.

The user wants a better multi-agent workflow. The desired result is not simply
more parallel agents. It is a clear way to decide:

- whether multiple agents help at all;
- when work can run in parallel and when it must be sequenced;
- who owns each piece of state and publication step;
- how each piece proves it worked before later work depends on it;
- how a failed or restarted run converges without duplicate ownership; and
- when a human should decide product direction versus when the manager should
  make an ordinary execution choice.

The user also called out that the current OpenCode path is poor and that the
skill contains far too much code merely to support OpenCode.

## Current inventory

`skills/implement-tickets/` currently contains:

| Part | Size or role |
| --- | --- |
| `SKILL.md` | 168 lines. User-facing router plus substantial manager policy. |
| `references/core-workflow.md` | 256 lines. Readiness, checkout lifecycle, worker contract, publication, feedback, recovery, completion. |
| `references/project-setup.md` | 160 lines. Repository contract and preflight. |
| Harness references | Codex: 112 lines. Grok: 207 lines. OpenCode: 165 lines. |
| Work-source references | GitHub, GitLab, and Jira adapters. |
| `scripts/orchestrate.go` | 4,098 lines. Durable supervisor, tracker clients, state, process handling, review publication, feedback, and harness behavior. |
| `scripts/orchestrate_test.go` | 693 lines. Tests for that supervisor. |

The total is about 6,090 lines. The Go supervisor is the dominant source of
complexity. It contains tracker adapters, Git behavior, queue scheduling,
worker launch and recovery, review and feedback policy, durable state, logs,
and OpenCode-specific parsing. It is not a small reliability helper.

Examples of mixed responsibilities in `orchestrate.go` include:

- GitHub, GitLab, and Jira discovery and state normalization;
- worker prompts, isolated checkout creation, claims, and recovery;
- Codex and OpenCode launch behavior plus Grok exceptions;
- publication, draft reviews, integration checks, and item closure;
- feedback authorization and comment handling;
- state locks, event logs, process health, and cleanup.

This is why the current skill is hard to reason about and hard to make robust
for OpenCode. A skill should describe how to make a decision and use a small
reliable tool where necessary. This package embeds a long-running controller
with a wide external API surface.

## Current strengths worth retaining

Do not throw away the underlying safety goals merely because the form is wrong.
The current workflow has several sound ideas:

- One worker owns one bounded item and branch at a time.
- Workers work in isolated checkouts.
- The manager, not the worker, owns publication and tracker writes.
- A dependent item waits until an explicit blocker is integrated, not merely
  tracker-closed.
- A completed item should be verified and published as soon as it is ready,
  rather than waiting for a whole batch.
- Review feedback should be classified and authorized before dispatching more
  work.
- A human controls merging.
- Recovery should preserve existing work and inspect it before replacement.

These are candidate rules for a new, simpler workflow. They must be rejustified
against the actual host harness rather than preserved because they exist today.

## Pstack principles that matter

Pstack does not provide a full persistent ticket manager. Its value here is a
set of smaller operating principles.

### Sequence work into verifiable units

This is the central principle for the rework. A unit has a known starting
state, one bounded change, a direct check, and a result that can be delivered
or used by the next unit. Do not run a broad batch and discover failures only
at the end.

For ticket orchestration, a ready item should become:

1. a clearly stated outcome and acceptance condition;
2. one worker with one isolated writable area;
3. one direct verification of the item;
4. manager review and publication;
5. integration proof; then
6. recomputation of work newly unblocked by that integration.

Commits and reviews should also be independently understandable. Their order
should help a reviewer follow the proof rather than hide many unrelated edits.

Source: [Pstack sequence verifiable units](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-sequence-verifiable-units/SKILL.md).

### Make operations idempotent

Every manager action needs a restart answer: what happens if it runs twice, or
if it stopped partway through?

The desired property is convergence, not a fragile memory of the previous run.
Before claiming, launching, publishing, commenting, closing, or cleaning up,
inspect actual branch, checkout, review, tracker, and process state. Adopt a
valid existing result or repair it. Do not create a second owner or duplicate
publication.

Source: [Pstack make operations idempotent](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-make-operations-idempotent/SKILL.md).

### Swarm is a bounded tool, not the ticket manager

Pstack's `swarm` is useful for a finite fan-out that drains and returns one
report. It fits discovery, a review panel, test investigation, or a short race
between alternative approaches. It does not fit a persistent queue that waits
for review, human merge, or later dependencies.

Use a swarm inside an orchestration run only when the work is finite and the
manager can aggregate it immediately. Do not model the full ticket lifecycle
as a swarm.

Source: [Pstack swarm](https://github.com/cursor/plugins/blob/main/pstack/skills/swarm/SKILL.md).

### Outcome-oriented execution is a migration rule

For a planned rework or repository migration, prefer the verified target state
over compatibility layers that only serve an intermediate step. This supports a
clean replacement of the current workflow. It is not a license to leave normal
ticket work broken between steps.

Source: [Pstack outcome-oriented execution](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-outcome-oriented-execution/SKILL.md).

### Existing local principles already apply

- `act-within-scope`: collect material questions early. Once execution starts,
  make reasonable reversible choices instead of stopping for ordinary questions.
- `verification`: observe real outcomes. A passing proxy or worker self-report
  does not prove the requested change landed.
- `experience-first`: manager and worker contracts should make the human
  reviewer and ticket author successful, not merely make automation convenient.
- `encode-lessons`: repeated orchestration failures should become the smallest
  durable guard, not another paragraph of policy.
- `context-management`: retain compact state and evidence rather than replaying
  raw worker logs into every manager turn.

## Proposed architecture to discuss

This is a recommendation, not a settled design.

### Keep one user-facing entrypoint

Keep `implement-tickets` as the public skill name. A user asking to implement a
set of tracker items should not need to discover a second generic skill before
the workflow can begin.

Rewrite its body as a concise entrypoint. It should select the work source,
inspect the repository's local contract, decide whether the selected work is
parallelizable, and load the smallest relevant reference.

Do not add a second top-level "multi-agent logic" skill unless it develops a
separate user request and output. A second skill that always accompanies
`implement-tickets` adds another activation and routing problem without
reducing complexity.

### Make the core provider-neutral

Create one concise reference such as `references/orchestration-model.md`. It
should contain only durable reasoning:

- choose single-worker, sequential, parallel, or bounded-race execution;
- make dependencies explicit;
- assign one owner per unit, branch, and mutable resource;
- define a completion and verification condition before dispatch;
- publish verified units independently;
- recompute readiness after integration;
- retry by reconciling real state; and
- stop for genuine product, authority, or destructive-action boundaries.

This reference must not contain OpenCode flags, Codex CLI syntax, tracker JSON,
or provider/model defaults.

### Keep adapters narrow

Tracker adapters should answer only source-specific questions: list selected
items, read full bodies and comments, locate explicit blockers, publish a
review or comment when authorized, and prove integration.

Harness adapters should answer only host-specific questions: create an isolated
worker, give it the worker contract, read its handoff, and request or stop work
when the host genuinely supports it.

An adapter must describe an actual tested capability. It must not claim that a
host can persist a chat turn, wake after a final answer, own a process safely,
or provide isolation that it does not provide.

### Separate the OpenCode decision

The OpenCode supervisor should not remain bundled merely because it exists.
There are three plausible paths:

1. **Remove OpenCode support for now.** Keep the rework small and support only
   harnesses whose workflow can be directly exercised. Reintroduce OpenCode
   after a narrow adapter can be proven.
2. **Keep an experimental OpenCode adapter.** Document only a bounded,
   foreground, one-run flow. Do not promise durable unattended operation.
3. **Move durable OpenCode orchestration into a separate tool.** Give it its own
   repository or lifecycle, explicit CLI contract, release process, tests, and
   operator documentation. The skill would call that tool rather than carrying
   thousands of lines of implementation.

The default recommendation is option 1 unless there is an active real use case
that proves options 2 or 3 are worth their cost.

### Avoid a universal daemon

Do not recreate the current supervisor under a new name. A persistent manager
has real operational requirements: state durability, authentication, process
ownership, cleanup, locking, retries, and external write authority. If those
requirements are essential, they deserve a deliberately designed tool. If they
are not essential, the active chat manager should own a bounded run and report
where human action or a later invocation is required.

## Worker contract draft

Each worker should receive a short, standalone brief that names:

- the outcome and acceptance criteria;
- the exact bounded item and repository path or checkout it owns;
- relevant repository instructions and the implementation skill to use;
- expected direct verification;
- its handoff format: outcome, evidence, remaining risk, and changed paths;
- boundaries it does not own, such as merging, production writes, or unrelated
  cleanup.

The manager should not require a worker to invent a commit subject convention
that conflicts with the repository's own Git workflow. The current skill still
requires Conventional Commit subjects, which conflicts with the repository's
new descriptive commit-message policy. Resolve that during the rewrite.

## Planning and scheduling rules to preserve

- Parallelize only independent units with separate writable areas.
- Sequence work when items share mutable state, one changes the API used by
  another, or the first item's observed result decides the second item's shape.
- Do not infer dependencies from ticket order. Use explicit blockers or write a
  visible planning decision before scheduling.
- Keep concurrency modest. More workers do not make an unclear plan clearer.
- Finish verification and publication of a completed unit while other genuinely
  independent workers continue.
- Do not make an artificial wave completion barrier.
- Treat human review and merging as external state. The workflow may poll when
  its host can do so honestly, but it must never promise a wake-up mechanism it
  does not have.

## What not to carry forward unchanged

- A 4,000-line supervisor hidden in a skill.
- A mandatory provider, model, reasoning level, or default worker identity.
- OpenCode-specific behavior in the generic entrypoint.
- A claim that all harnesses can run unattended, persist state, or resume chat
  work in the same way.
- Conventional Commit requirements. Follow the repository's Git convention.
- Automation that merges, self-approves, or makes product decisions for a human.
- Unbounded backlog fan-out.
- Tests that merely lock in incidental controller internals without proving a
  meaningful orchestration guarantee.

## Questions to settle before editing

1. Is OpenCode needed now? If yes, which concrete run must it support:
   foreground selected tickets, persistent unattended work, or both?
2. Which harnesses matter enough to support first? Codex is the current host.
   Grok and OpenCode should not remain equal-priority by default.
3. Is the desired normal unit a tracker ticket, a user-selected task, or either?
4. What manager writes are standing authority in the intended use: branch push,
   draft review creation, ticket comment, ticket closure after verified merge?
5. Does the user want durable runs that survive a chat/session ending? If yes,
   decide whether that is a separate tool rather than a skill feature.
6. What is the smallest representative workflow to prove before adding more:
   two independent GitHub issues, one dependency chain, review feedback, or
   restart recovery?
7. Which outcomes deserve permanent automated tests? The answer should be
   stable orchestration guarantees, not every text or internal sequence.

## Suggested rework sequence

1. Answer the questions above and choose the first supported harness and
   tracker source.
2. Write a short target contract for one bounded run: inputs, ownership,
   terminal states, direct evidence, and explicit external boundaries.
3. Replace the entrypoint and create the provider-neutral orchestration model.
4. Retain or rewrite one tracker adapter and one harness adapter only.
5. Exercise a real local or disposable repository run with two verifiable
   units. Observe branch isolation, worker handoff, direct verification,
   publication, and dependency release.
6. Add only the small durable checks that protect the chosen contract.
7. Add another provider or tracker only after its adapter can be demonstrated
   without expanding the generic core.
8. Decide explicitly whether to delete, archive, or extract the OpenCode Go
   supervisor. Do not leave it half-supported.

## Related repository material

- `skills/implement-tickets/SKILL.md`
- `skills/implement-tickets/references/core-workflow.md`
- `skills/implement-tickets/references/project-setup.md`
- `skills/implement-tickets/references/harnesses/`
- `skills/implement-tickets/references/work-sources/`
- `skills/implement-tickets/scripts/orchestrate.go`
- `skills/implement-tickets/scripts/orchestrate_test.go`
- `docs/research/pstack-multi-agent-workflow.md` once the parallel source review
  completes.

## Next session

Start by reading this note and `docs/research/pstack-multi-agent-workflow.md`.
Do not begin by editing the existing skill. First choose the initial supported
runtime and the OpenCode disposition. Then build and test one narrow real run
before generalizing.
