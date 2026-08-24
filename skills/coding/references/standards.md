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

## Assertions

- Use assertions for internal invariants and impossible states, not for
  expected runtime failures such as user input errors or I/O problems.
- Prefer fail-fast behavior for invariants where continuing could cause
  corruption, unsafe output, or silent correctness drift.
- Default to debug assertions for expensive checks or checks whose production
  crash cost is higher than the risk of continued execution.
- Promote debug assertions to production assertions when violation means the
  process cannot safely continue and the check is cheap.
- Keep assertion expressions side-effect free and messages actionable.
- For detailed assertion policy, load `./assertions.md`.

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
- Keep short linear functions and cohesive flows inline when extraction would
  only add a name, file, or call to learn.
- Extract a helper or create a module only when it gives callers a clearer
  operation or owns coherent work they should not need to understand: an
  invariant, lifecycle, policy, representation conversion, or integration
  detail.
- Do not extract solely to remove repeated syntax. A few similar statements can
  be clearer than a generic helper; consolidate only when a named operation can
  own their shared semantics and reduce future caller burden.
- Keep orchestration and caller-owned policy close to the caller. Helpers should
  execute the coherent work the boundary owns.

## Object State And Ownership

- When callers repeatedly assemble the same meaningful state transition or
  configuration, give the owning object a method or focused helper. Do not add
  one for simple assignment with no named behavior, invariant, default, or
  validation to own.
- Do not hide fields reflexively. Public fields may be part of the intended API,
  useful for literals, serialization, tests, or low-level data access. In those
  cases, keep the fields public if appropriate and still offer helpers for
  common state changes.
- Avoid helpers that only assign one field without adding a named behavior,
  invariant, defaulting rule, validation rule, or repeated configuration
  pattern.

## Abstraction Boundaries

- Keep cohesive work direct when the caller already owns the decisions and can
  understand the flow in place.
- Add a boundary only when it reduces total reader work by giving callers a
  simpler operation or hiding coherent knowledge they should not carry.
- Do not introduce pass-through wrappers that merely mirror SDK or API methods,
  or interfaces created only for a hypothetical second implementation.
- Prefer direct SDK/API calls at integration edges when a wrapper adds no domain
  semantics, policy, translation, lifecycle, or other meaningful compression.

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
