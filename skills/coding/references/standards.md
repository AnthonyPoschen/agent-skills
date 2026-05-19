# Coding Standards

Use these standards when project tooling and systemic local patterns do not
already settle the decision.

## Control Flow

- Put simple fail-state checks first and return immediately.
- Avoid `if/else` when the `if` branch exits.
- Keep independent fail checks as separate guards unless combining them improves
  correctness, diagnostics, or readability.
- Do not count loop or switch nesting as conditional-chain nesting.
- Prefer explicit boolean comparisons when returning or assigning negated
  boolean expressions, especially when the expression is a method call or
  property access. Use `value == false` over `!value` when it makes the
  evaluation harder to miss. Follow stronger local language/project conventions
  if they consistently prefer unary negation.

## Functions

- Keep functions focused on one primary responsibility.
- Prefer single-line signatures when the formatter and language conventions
  allow it.
- Do not split signatures only to satisfy manual width preferences.
- Treat five or more parameters as a design smell: reuse existing cohesive
  types first, split unrelated responsibilities, or introduce a focused input
  type for a cohesive subset.
- Do not bundle unrelated inputs into a catch-all wrapper just to reduce a
  parameter count.
- Keep short linear functions inline unless extraction is needed for reuse,
  substantial duplication, API boundaries, or testing.
- Extract helpers only when they represent meaningful work: cohesive behavior,
  reusable logic, independently testable detail, or a substantial phase.
- Extract shared behavior when multiple functions reveal the same traversal,
  lookup, parsing, validation, dispatch, or edge-handling policy. The decision
  is not based on duplicate syntax alone; extract when the callers can be
  combined around one named behavior and future drift would create bugs or
  inconsistent semantics.
- When sibling functions have near-identical bodies with small substitutions
  such as expected states, fields, operators, enum cases, or callback behavior,
  keep the public functions when they are the right consumer-facing API, but
  combine their internals behind a helper or parameterized operation. This is
  especially important when one small function could drift and become subtly
  wrong.
- When combining sibling functions that differ by current/previous state,
  expected enum value, field, slice, or predicate input, prefer passing the
  concrete varying data and expected value into the helper. Do not introduce a
  mode enum, wrapper type, or helper name that bakes in one specific case when
  callers can pass the resolved state directly. The helper should name the
  shared relation, not one caller's outcome.
- Hoist shared preparation out of branches when many switch/conditional cases
  repeat the same sanitization, normalization, parsing, validation, lookup, or
  derived-value calculation. Name the prepared values so the branch table is
  about selection and policy, not repeated setup. Do not introduce temporaries
  for one-off expressions that are clearer inline.
- Parent functions should retain orchestration and policy decisions; child
  helpers should execute resolved steps.

## Object State And Ownership

- When multiple callers heavily modify the same object fields in the same way,
  prefer a method or focused helper that applies the named state transition or
  configuration pattern.
- Treat repeated external mutation as a sign that ownership may belong on the
  object. This is especially important when future changes would require
  updating several callers to keep invariants, defaults, derived fields, or
  validation consistent.
- Do not hide fields reflexively. Public fields may be part of the intended API,
  useful for literals, serialization, tests, or low-level data access. In those
  cases, keep the fields public if appropriate and still offer helpers for
  common state changes.
- Avoid helpers that only assign one field without adding a named behavior,
  invariant, defaulting rule, validation rule, or repeated configuration
  pattern.

## Abstraction Boundaries

- Do not introduce pass-through wrappers that merely mirror SDK or API methods.
- Add wrapper/helper layers only when they add policy, retries, validation,
  translation, composition, or a meaningful test boundary.
- Prefer direct SDK/API calls at integration edges when a wrapper adds no domain
  semantics.
- Do not create internal objects solely to hide provider names.

## State And Dependencies

- Avoid mutable global runtime state when values can be created at startup and
  injected.
- Prefer dependency injection over package-level singletons for runtime
  collaborators.
- Use the smallest practical variable scope.
- Declare immutable values with the language's const/final/readonly equivalent
  when the value never changes.

## Magic Values

- One-off obvious literals are acceptable in narrow scope.
- Reused literals within a scope or module should become named constants or
  variables.
- Reused literals across files should move to the narrowest shared boundary.
- Shared constants may be global only when immutable domain constants with broad
  reuse.

## Comments And Docs

- Public/exported APIs, non-obvious behavior, domain policy, edge cases, and
  complex helpers should have short comments that explain purpose, contract, or
  consequence.
- Private obvious helpers and small local functions usually do not need
  comments.
- Do not narrate obvious syntax or line-by-line operations.
- Use phase comments only when they materially improve skimming of a non-trivial
  block.
- Keep comments synchronized with changed code.
- Use exact uppercase tagged prefixes: `TODO:`, `FIXME:`, `NOTE:`, `HACK:`,
  `WARNING:`.
- `FIXME:`, `HACK:`, and `WARNING:` should include the consequence if ignored.
- `HACK:` should include a removal condition when one is known.
