# Refactoring Workflow

Use this reference for cleanup, simplification, extraction, restructuring, and
rewrite work. Decision rules live in `./principles.md`. Code-shape rules live
in `./standards.md`. This file is the refactor sequence.

## Refactor Contract

- Preserve externally observable behavior unless the user requests a behavior
  change. See Contracts in `./principles.md`.
- Establish a known-good baseline before refactoring when practical: run the
  project's relevant tests or checks first so later failures can be attributed
  to the change. If baseline tests already fail or cannot be run, record that
  before editing.
- Use `./standards.md` as the rewrite target, constrained by project tooling
  and systemic local patterns.
- Keep the refactor inside Scope in `./principles.md`.
- Prefer a sequence of understandable mechanical changes over a clever rewrite
  that is hard to verify.

## Sequence

Run these phases in order until the completion check. Subtraction is the first
edit phase, not the end of the refactor.

1. Establish the test or check baseline.
2. Subtract dead and accidental complexity in the selected area.
3. State the target design for the remaining shape.
4. Rewrite toward that target through one coherent boundary at a time.
5. Run the completion check.

If the selected area is large, finish one coherent boundary through rewrite,
then take the next boundary. Do not leave the request on a subtract-only tree.

Commit between phases only when the user asked to land the refactor as
separate commits. Those commits use `git-commit-workflow`. Otherwise leave the
working tree uncommitted until the user asks.

## Subtract Before You Rebuild

Apply Subtract before you add in `./principles.md` to the selected area, then
continue this sequence. Do not stop after subtraction.

## Design The Target First

For a meaningful refactor, design the selected area as though its current
requirements had existed from the beginning. Use the existing code to learn
about callers, contracts, and migration risk. Do not let its accidental shape
dictate the target design.

- State the smallest coherent target before making mechanical edits. Prefer the
  shape that makes the affected callers and responsibilities simplest.
- Carry that target through the selected boundary: directly affected callers,
  types, tests, examples, documentation, and obsolete paths.
- Apply Fix duplication while the change is open and Boundaries earn their cost
  from `./principles.md`. Do not retain an awkward intermediate API or
  compatibility wrapper to make the refactor smaller.
- Keep the work inside the area the user asked to refactor. Report a wider
  redesign opportunity instead of turning the work into a repository rewrite.

## Rewrite Rules

Apply these homes while rewriting. Do not restate them:

- `./principles.md`: Boundaries earn their cost, Fix duplication while the
  change is open, and A write that stuck stays stuck.
- `./standards.md`: control flow, functions, object state, dependencies, magic
  values, and comments.

## Refactor Boundaries

- Apply Scope from `./principles.md`. If a local legacy pattern conflicts with
  standards, update the path needed for the requested change.
- When a touched test fails or becomes awkward, apply the Tests principle in
  `./principles.md`.
- When moving or splitting code, use `./file-organization.md`.
- Prefer tests or characterization checks before changing complex behavior.

## Completion Check

- Behavior is preserved, or intentional changes are explicit.
- The touched code is closer to `./standards.md` without broad unrelated churn.
- Baseline test status is known from before the refactor, or the reason it
  could not be established is documented.
- Apply Review before done and Prove the outcome from `./principles.md`.
- Tests cover the risk introduced by the rewrite and pass after the change,
  except for pre-existing failures documented before editing. New and removed
  tests follow the Tests principle in `./principles.md`.
