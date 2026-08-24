# Refactoring Workflow

Use this reference for cleanup, simplification, extraction, restructuring, and
rewrite work.

## Refactor Contract

- Preserve externally observable behavior unless the user requests behavior
  change.
- Establish a known-good baseline before refactoring when practical: run the
  project's relevant tests/checks first so later failures can be attributed to
  the change. If baseline tests already fail or cannot be run, record that
  before editing.
- Use `./standards.md` as the rewrite target, constrained by project tooling and
  systemic local patterns.
- Keep the refactor scoped to touched code and directly related collaborators.
- Prefer a sequence of understandable mechanical changes over a clever rewrite
  that is hard to verify.

## Subtract Before You Rebuild

Begin by removing what the selected area does not need. Deletion exposes the
essential shape, reduces migration work, and keeps the redesign from preserving
complexity just because it was already present.

- Remove dead paths, pass-throughs, duplicated decisions, unnecessary
  conversions, stale compatibility code, and state that can be derived instead
  of synchronized.
- Cut unused scope before improving names, structure, or polish. Build the
  target on the minimum behavior the selected area must retain.
- Do not retain unused options, validators, parsers, guards, or extension
  points for hypothetical cases.
- Keep subtraction within the selected boundary and verify its effect. Report
  broader cleanup opportunities rather than widening the refactor without
  approval.

## Design The Target First

For a meaningful refactor, design the selected area as though its current
requirements had existed from the beginning. Use the existing code to learn
about callers, contracts, and migration risk. Do not let its accidental shape
dictate the target design.

- State the smallest coherent target before making mechanical edits. Prefer the
  shape that makes the affected callers and responsibilities simplest.
- Carry that target through the selected boundary: directly affected callers,
  types, tests, examples, documentation, and obsolete paths.
- Do not retain an awkward intermediate API, compatibility wrapper, or
  special-case adapter merely to make the refactor smaller.
- Keep the work inside the area the user asked to refactor. Report a wider
  redesign opportunity instead of silently turning the work into a repository
  rewrite.

## Rewrite Rules

- Flatten conditional chains with guard clauses when failure states are simple.
- Remove avoidable `else` branches after `return`, `continue`, `break`, or
  equivalent exits.
- Separate independent validation/failure checks unless combining them improves
  correctness or error reporting.
- Keep short linear functions and cohesive flows inline when extraction would
  only relocate the same thinking into another name or file.
- Introduce a boundary when it gives callers a clearer operation or centralizes
  an invariant, lifecycle, policy, representation conversion, or integration
  detail that belongs elsewhere.
- Do not extract solely to remove duplicate syntax. A small amount of explicit
  repetition is preferable to a generic helper that obscures the work.
- Consolidate repeated decisions or mutations when a named operation can own
  their shared semantics and make callers simpler. Keep direct field access or
  inline code when there is no invariant or meaningful operation to own.
- Keep orchestration and caller-owned policy close to the caller; move only the
  coherent work the new boundary owns.
- Replace mutable global runtime state with startup construction and dependency
  injection when the touched area allows it.
- Consolidate repeated literals at the narrowest useful scope.
- Remove pass-through wrappers that add no policy or domain value.

## Refactor Boundaries

- Do not rewrite unrelated legacy code merely because it violates current
  standards.
- If a local legacy pattern conflicts with standards, update only the path
  needed for the requested change unless the user asked for broader cleanup.
- When a touched test fails or becomes awkward, evaluate the stable contract it
  protects before updating it. Keep a meaningful contract check; remove or
  replace a test that only pins incidental diagnostics, private mechanics, or
  stale mock behavior. Never delete a test merely to conceal a real failure.
- When a touched boundary creates the same conversion, sequencing, or workaround
  burden for several direct consumers, assess those consumers together. Redesign
  and migrate the bounded set when it lowers total reader work; otherwise keep
  the current change direct and leave the larger redesign out of scope.
- Do not solve one caller's friction by adding a special-case escape hatch,
  compatibility layer, or new adapter when the underlying boundary is the real
  problem.
- When moving or splitting code, use `./file-organization.md`. Split only when
  a part has a clear name and is a place a reader would look for separately.
  Do not split by line count, one type per file, or one method per file.
- Do not move unrelated code to normalize the repository layout. If the current
  task needs a clearer home, move the smallest coherent unit. Otherwise report
  the broader layout issue without starting a migration.
- Preserve names and shapes that are part of public API, serialization,
  migrations, config, CLI flags, env vars, URLs, or test fixtures unless the
  change explicitly includes them.
- Prefer tests or characterization checks before changing complex behavior.

## Completion Check

- Behavior is preserved or intentional changes are explicit.
- The touched code is closer to `./standards.md` without broad unrelated churn.
- Baseline test status is known from before the refactor, or the reason it could
  not be established is documented.
- Tests or equivalent verification cover the risk introduced by the rewrite and
  pass after the change, except for pre-existing failures that were documented
  before editing.
