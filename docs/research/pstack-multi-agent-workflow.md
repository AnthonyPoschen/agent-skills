# Pstack multi-agent workflow research

## Purpose

This note is a handoff for redesigning `skills/implement-tickets`. It compares
the current skill with Pstack's primary-source skills. It recommends a small,
provider-neutral orchestration model rather than importing Pstack wholesale.
It does not change the skill.

## The short conclusion

Keep `implement-tickets` as the user-facing workflow, but rebuild it around
one idea: a ticket is a **verifiable delivery unit**, not merely work assigned
to an agent.

Use parallel workers only for independent units. Give every unit one owner,
one isolated write area, a clear done predicate, and a direct check. Verify
and publish a completed unit immediately, then recompute what is ready. Treat
restart and retry as reconciliation with the existing tracker, Git, review,
and worker state.

Do not model the persistent ticket pipeline after Pstack's `swarm`. `swarm` is
a finite fan-out tool: it frames a task, launches a fixed number of workers,
drains them, and returns one report. It is useful within ticket orchestration
for a short discovery, review, or design race, not for ownership, review
feedback, dependency gating, or an ongoing work queue.

## Current skill: useful rules and structural problem

The existing skill already states many good rules:

- Workers own one vertical slice in an isolated checkout.
- Dependencies come only from explicit blockers and unblock only after Git
  proves integration into the target branch.
- A worker leaves a tested diff. A manager independently verifies, commits,
  publishes, and opens or updates the review.
- Completion is processed per item rather than at an artificial batch barrier.
- Retries preserve the existing checkout and inspect it before a replacement
  worker continues.
- A human owns merge authority.

Those rules are in `skills/implement-tickets/SKILL.md` and
`references/core-workflow.md`. They are compatible with Pstack's principles.

The structural problem is that the skill also embeds a full, multi-provider
orchestration product. The bundle is 5,380 lines: the `SKILL.md`, references,
assets, a 4,098-line Go supervisor, and 693 lines of supervisor tests. The
supervisor knows GitHub, GitLab, Jira, Codex, Grok, OpenCode process handling,
review feedback authorization, persistence, polling, recovery, and
publication. Its OpenCode-specific paths are spread through the shared model,
launching, artifact parsing, recovery, commit-subject fallback, and health
repair. This makes a provider limitation shape the general workflow.

The redesign should preserve the safety invariants above while making the
normal path understandable without learning the supervisor's internals.

## What to adopt

### Make the unit the unit of delivery

Pstack's [Sequence work into verifiable
units](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-sequence-verifiable-units/SKILL.md)
says each unit should end in a check, be verified before the next unit, and be
ordered so the resulting commits or PRs make the argument visible to a
reviewer.

For tickets, define one unit as:

1. Read the item and record its acceptance condition, baseline, dependencies,
   owner, and direct verification method.
2. Give one worker the bounded slice and one isolated writable branch or
   checkout.
3. Have the worker implement and run the closest appropriate check.
4. Independently inspect and verify the actual result.
5. Commit, publish, and create or update the review as one delivery event.
6. Treat the unit as dependency-satisfying only after its review is merged and
   Git proves the target contains that merge.

This rule applies at two levels. Within an item, do not build a sequence of
unverified changes. Across items, do not build dependent work atop an
unpublished or unverified predecessor. A ready queue is not a reason to defer
the completed item's verification or publication until a whole wave ends.

The existing per-item publication rule already moves in this direction. Make
it the first architectural concept and reduce the surrounding mechanics to
support it.

### Reconcile instead of retrying blindly

Pstack's [Make operations
idempotent](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-make-operations-idempotent/SKILL.md)
requires operations to converge after repeated execution and partial crashes.
It calls for scanning existing state, adopting live sessions, content-based
cleanup, stale-lock detection, and a reconciliation step whenever the answer
to a restart question depends on leftover state.

Adopt this as the lifecycle rule for every state-changing orchestration step:

- `start` reads existing tracker assignment, branch, checkout, worker/session,
  review, target head, and verification evidence before creating anything.
- `dispatch` refuses a second active owner for a unit. It resumes or replaces
  only after inspecting the prior worker's preserved work.
- `publish` first discovers a matching commit and review, then creates only
  missing state.
- `integrate` proves that the review's merge is an ancestor of the fetched
  target before changing ticket readiness.
- `cleanup` removes only content-equivalent, clean, run-owned artifacts.
- A durable lock has an owner identity and stale-owner recovery, not merely an
  instruction that agents should not overlap.

The important result is not a large state machine. It is a small, explicit
reconciliation boundary around each external write.

### Separate writes before adding locks

Pstack's [Separate before serializing shared
state](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-separate-before-serializing-shared-state/SKILL.md)
is directly relevant even though it was not originally in the shortlist. It
says to eliminate shared mutable targets first; serialize only when one shared
writer is a real invariant.

Adopt these ownership boundaries:

| State | Owner | Rule |
| --- | --- | --- |
| Ticket implementation diff and branch | One worker | No other worker writes it. |
| Verification, commit, push, review, and tracker transition | One manager/publisher | Workers have no publication authority. |
| Tracker, Git remote, and review system | External sources of truth | The manager re-reads them after writes and on restart. |
| Optional run journal | Coordinator | It records recovery data. It does not replace external truth. |

Workers should not share a branch, checkout, prompt file, or mutable run-state
object. A single publisher is a justified serialized boundary because a review
and target-branch transition have one canonical meaning. Do not introduce a
global lock for work that could instead have separate branches and result
records.

### Use a lever, but keep it narrow

Pstack's [Build the
lever](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-build-the-lever/SKILL.md)
prefers a repeatable script, generator, codemod, query, or delegate contract
when work is non-trivial. The lever improves both repeatability and review.
It also says a deterministic lever beats fan-out when it can process the work
directly, and that a delegate skill can be the shared recipe, verification
contract, and do-not-touch fence.

Adopt the principle in two bounded forms:

- Put the worker brief, ownership fence, acceptance format, and verification
  expectation in a stable contract that every worker receives. Do not keep
  rebuilding this policy in prompt strings.
- Automate only deterministic recurring facts: isolated-checkout setup,
  branch naming, status collection, or a single provider's worker launch. A
  tool should make one real operation repeatable and observable.

Do not use this principle to justify another all-purpose supervisor. The
current 4,098-line program is a broad lifecycle product, not the smallest
lever for a narrow task. Each future adapter needs a small capability contract
and a reason to exist. If OpenCode needs durable process supervision, make it
a separately versioned tool or experimental adapter with its own tests and
lifecycle, rather than letting it dictate the core skill's model.

### Keep adapters at boundaries

Pstack's [Boundary
discipline](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-boundary-discipline/SKILL.md)
puts validation and defensive logic at external boundaries, while keeping the
inside typed and direct. Applied here:

- Work-source adapters translate GitHub, GitLab, Jira, or local-ticket facts
  into a small normalized item, dependency, review, and integration result.
- Harness adapters launch, inspect, stop, and resume a worker. They do not own
  ticket policy, review authorization, or Git publication policy.
- The core scheduling policy consumes normalized facts and has no provider
  command strings, provider-specific JSON parsing, or permission-model
  branches.
- The worker contract is the only boundary where free-form ticket text becomes
  an implementation request. The worker returns a structured handoff plus
  real evidence.

This is a reason to keep a provider-neutral `orchestration-model.md` reference
and small adapters, not an argument to create a top-level skill for every
provider.

### Optimize the rewrite for its target state

Pstack's [Outcome-oriented
execution](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-outcome-oriented-execution/SKILL.md)
is for planned rewrites and migrations with explicit phase boundaries. It
says to optimize for the verifiable end state rather than preserving every
intermediate compatibility layer.

Use it to plan the *rewrite of `implement-tickets`*, not as the normal ticket
execution rule. State the target architecture first, declare which temporary
breakages are acceptable, preserve high-signal checks in affected paths, and
run full validation at explicit migration boundaries. Do not retain a bloated
OpenCode-compatible core solely to make every intermediate commit support the
old lifecycle. If an adapter cannot fit the small core, park it as experimental
or remove it until it can.

### Preserve the cheapest shape

Two other Pstack principles strengthen the same redesign:

- [Laziness protocol](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-laziness-protocol/SKILL.md): delete before adding, keep call paths flat, consolidate decisions, and question broad signal threading. This supports a small coordinator that owns scheduling rather than a multi-layer manager/supervisor/adapter policy maze.
- [Minimize reader load](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-minimize-reader-load/SKILL.md): collapse one-caller wrappers and mutable scope that do not reduce the reader's work. The core execution path should let a maintainer answer “what owns this ticket?” and “what happens next?” without tracing provider code.

These principles complement the library's existing `subtract-before-you-add`,
first-principles refactoring, and foundation guidance.

### Prove the one safety fact that matters

Pstack's [Blast radius](https://github.com/cursor/plugins/blob/main/pstack/skills/blast-radius/SKILL.md)
asks the reviewer to find the one or two facts a change is safe because of,
then prove those facts with real code rather than a persuasive write-up. It
also calls out risks that symbol search misses: wire formats, database columns,
external library behaviour, other languages, feature flags, and teardown or
ordering behaviour.

Use this as a review gate for a wide or high-risk ticket unit, its integration,
or a shared change that could affect other queued work. The manager should
name the safety fact, run the cheapest relevant direct proof, and label a fact
unproven if it cannot reach executable evidence. Do not make a generic
blast-radius report a scheduling prerequisite for every small ticket.

## What not to adopt

| Pstack idea | Do not adopt it as | Why |
| --- | --- | --- |
| `swarm` | The persistent ticket scheduler | It deliberately drains a fixed worker set and returns one report. It has no long-lived ownership, dependency, review, or recovery model. |
| `swarm` | The default for every ticket set | Parallelism only helps independent ready units. One coherent change should stay with one owner; dependent work must wait for verified integration. |
| `build-the-lever` | Permission to build a generic orchestration framework | The lever must be the smallest artifact that does or proves the current repeated operation. |
| `outcome-oriented-execution` | Permission to leave ordinary tickets broken until the end | It is for declared migration phases. A ticket unit should still finish in a checkable state. |
| `boundary-discipline` | A mandate for adapters and interfaces everywhere | Use an adapter only at a real provider or authority boundary. Direct core code is better than a one-caller abstraction. |
| `exhaust-the-design-space` | Required parallel prototyping of established ticket work | Use a bounded design race only when the orchestration architecture has several viable, unconstrained shapes. |
| `blast-radius` | A long speculative risk report for each ticket | Use it only where the unit or integration has real breadth or risk, and prove the key safety fact rather than listing callers. |

Pstack's skill is named `blast-radius`, not `principle-blast-radius`, in the
current `main` tree. Its concern also reinforces the library's existing scope,
verification, and isolated-checkout rules: keep each unit's files, authority,
provider writes, and recovery action as small as practical.

## A proposed small core model

This is a design target, not an implementation prescription.

```text
discover selected items and external truth
        |
        v
normalize explicit dependencies and acceptance conditions
        |
        v
for each ready independent unit: claim -> isolate -> dispatch one owner
        |
        v
worker handoff -> independent proof -> publish review
        |
        v
human merge -> Git integration proof -> recompute ready set
```

The normal result table needs only the information a user and coordinator act
on: item, owner, state, direct verification, review, blocker, and next action.
Persist additional data only when a restart genuinely needs it.

Use a finite Pstack-style swarm inside this loop only for bounded work such as:

- parallel read-only discovery before splitting a large backlog;
- independent review or test-coverage passes of one completed unit; or
- a declared best-of architecture race before committing to an unfamiliar
  orchestration adapter.

For every swarm, declare the artifact, done predicate, slice or race rule,
isolated output, and evidence format, exactly as Pstack's
[Swarm](https://github.com/cursor/plugins/blob/main/pstack/skills/swarm/SKILL.md)
requires. Aggregate it into a decision. Do not leave it responsible for the
ongoing queue.

## Suggested rewrite phases

1. **Write the model first.** Replace the long entrypoint with a short
   provider-neutral `SKILL.md` and one `orchestration-model.md` that states
   unit, ownership, dependency, verification, publication, authority, and
   recovery rules.
2. **Choose a minimal initial capability set.** Support one work source and
   one harness that can be exercised end-to-end. A manager may run the loop in
   the interactive session before a durable service exists.
3. **Define adapter contracts from actual facts.** Add a work-source adapter
   or harness adapter only after it can demonstrate discovery, dispatch,
   worker health, evidence capture, and recovery without changing core policy.
4. **Move or remove the OpenCode supervisor.** Preserve it temporarily only
   as an explicitly experimental, separately owned tool. Do not make its
   process model or permission quirks a requirement of the core workflow.
5. **Add restart cases before automation grows.** For each external write,
   test a second invocation and each meaningful partial-state boundary.
6. **Evaluate with real scenarios.** Use one single-ticket change, independent
   tickets, a dependency chain, a worker crash with an existing diff, review
   feedback, and a target-branch update. Judge outcome and evidence, not just
   whether the coordinator stayed busy.

## Decisions still needed

- Should the initial rewrite support only an interactive Codex-style manager,
  or must it include an unattended durable loop on day one?
- Which single tracker and review combination is the first supported source of
  truth?
- Does OpenCode have enough stable process and sandbox semantics to justify a
  separately maintained adapter, or should it be removed until proven?
- Which manager writes are pre-authorized for a normal run: assignment, branch
  push, draft review, ticket comment, and cleanup? Human merge remains a
  deliberate authority boundary.
- Is the desired delivery unit one ticket/PR by default, or may a selected
  ticket be explicitly divided into ordered sub-units with their own review
  boundaries?

## Sources

- [Pstack: Swarm](https://github.com/cursor/plugins/blob/main/pstack/skills/swarm/SKILL.md)
- [Pstack: Sequence work into verifiable units](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-sequence-verifiable-units/SKILL.md)
- [Pstack: Make operations idempotent](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-make-operations-idempotent/SKILL.md)
- [Pstack: Separate before serializing shared state](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-separate-before-serializing-shared-state/SKILL.md)
- [Pstack: Build the lever](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-build-the-lever/SKILL.md)
- [Pstack: Boundary discipline](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-boundary-discipline/SKILL.md)
- [Pstack: Outcome-oriented execution](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-outcome-oriented-execution/SKILL.md)
- [Pstack: Laziness protocol](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-laziness-protocol/SKILL.md)
- [Pstack: Minimize reader load](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-minimize-reader-load/SKILL.md)
- [Pstack: Exhaust the design space](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-exhaust-the-design-space/SKILL.md)
- [Pstack: Blast radius](https://github.com/cursor/plugins/blob/main/pstack/skills/blast-radius/SKILL.md)
- `skills/implement-tickets/SKILL.md`
- `skills/implement-tickets/references/core-workflow.md`
- `skills/implement-tickets/references/harnesses/capabilities.md`
- `skills/implement-tickets/scripts/orchestrate.go`
