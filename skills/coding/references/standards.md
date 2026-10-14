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
- Extraction and new boundaries follow Boundaries earn their cost and Fix
  duplication while the change is open in `./principles.md`.

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

Apply Boundaries earn their cost in `./principles.md`. At an integration edge,
call the SDK or API directly when a wrapper adds no domain semantics, policy,
translation, lifecycle, or other meaningful compression.

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

- No comments by default. Two kinds stay:
  - Doc comments on public or exported APIs, stating the contract callers rely
    on (Go doc comments, Rust `///`, JSDoc on exports). Keep them short and in
    sync with the code.
  - A constraint we can't change (an external API quirk, a spec reference),
    with a link to the source.
- Express contracts with types, names and tests. See the
  `type-system-discipline` and `test-behavior-not-implementation` skills.
- A `HACK:` or `FIXME:` needs a linked issue, and lint rejects the bare ones.
- When you'd write a comment, turn it into a type, test or lint rule instead.
  Use the `encode-lessons` skill to pick the smallest guard, and the
  `no-comments` skill to audit a diff.
- Stock linters are only the baseline. Turn recurring patterns into custom
  lint rules with the `lint-rules` skill.
- [Examples](../examples/comments.md)
