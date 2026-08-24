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
- Start by looking for subtraction: dead paths, pass-throughs, duplicated
  decisions, unnecessary conversions, and state that can be derived instead of
  synchronized.

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
