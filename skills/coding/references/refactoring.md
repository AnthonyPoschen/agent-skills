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

## Rewrite Rules

- Flatten conditional chains with guard clauses when failure states are simple.
- Remove avoidable `else` branches after `return`, `continue`, `break`, or
  equivalent exits.
- Separate independent validation/failure checks unless combining them improves
  correctness or error reporting.
- Split responsibilities when a function mixes unrelated decisions, effects, or
  data transformations.
- Keep short linear functions inline when extraction would only rename comments
  or create one-helper-per-step choreography.
- Extract substantial cohesive phases when they improve reasoning, reuse,
  testability, or hide meaningful detail.
- Extract small repeated loops or condition groups when multiple functions can
  be combined around shared behavior: a common traversal protocol, lookup
  policy, parsing rule, validation rule, dispatch path, or edge-handling
  semantic. This is a behavioral decision, not a line-count decision.
- Combine near-identical sibling function bodies when they differ only by small
  substitutions such as expected states, fields, operators, enum cases, or
  callbacks. Preserve the public method/function names when they are the right
  API; refactor the repeated body behind them.
- Hoist repeated sanitization, normalization, parsing, validation, lookup, or
  derived-value calculation out of switch/conditional branches when many cases
  need the same prepared values. The refactor should make each branch express
  its choice, not repeat setup mechanics.
- Move repeated external object mutation into methods or focused helpers when
  multiple callers configure the same object state in the same way. Centralize
  the named state transition so invariants, defaults, validation, and derived
  fields change in one place. Keep fields public when that is part of the API,
  but do not force common configuration patterns to stay open-coded.
- Keep orchestration decisions in the parent; pass resolved values to helpers.
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
