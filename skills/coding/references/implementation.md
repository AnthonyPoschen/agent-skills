# Implementation Workflow

Use this reference when writing new code or extending behavior.

## Before Editing

- Identify the user-visible behavior, API contract, and success criteria.
- Inspect nearby files for systemic project patterns before choosing a design.
- Check project tooling and existing tests so the implementation fits the repo.
- Prefer existing helpers, types, framework conventions, and composition roots
  over new abstractions.

## Pattern Priority

Follow this order:

1. Correct behavior, safety, and explicit user requirements.
2. Project tooling and enforced conventions.
3. Systemic local architecture and naming patterns.
4. Shared standards from `./standards.md`.

When local code is inconsistent, follow the pattern with the clearest tool,
framework, or repeated usage support. If no pattern is strong, use the shared
standards.

## Implementation Rules

- Make the smallest coherent change that fully handles the request.
- Design from existing boundaries: keep domain logic, adapters, UI state,
  persistence, and transport concerns where the project already places them.
- Add new modules only when they match the repo's organization or remove real
  complexity.
- Avoid broad wrappers around SDKs or APIs unless the wrapper adds policy,
  validation, retries, translation, composition, or testability.
- Reuse existing cohesive structs, options, contexts, or config objects before
  creating new parameter containers.
- Add tests in proportion to risk: broader tests for shared behavior, narrow
  tests for isolated changes.
- Update docs, comments, fixtures, generated inputs, or examples only when they
  are part of the changed behavior.

## Completion Check

- Before finishing, review the completed diff as if it were a code review:
  check correctness, public behavior, shared abstractions, duplicated internals,
  object ownership, tests, docs, and whether the change stayed scoped.
- The code follows systemic local patterns where they exist.
- The code falls back to `./standards.md` where local guidance is absent.
- Review the completed change for shared-behavior opportunities discovered
  during implementation. If multiple new or touched functions now perform the
  same traversal, lookup, parsing, validation, dispatch, or edge handling,
  extract one named helper when the behavior can be combined cleanly.
- Review object state changes before finishing. If multiple callers configure or
  mutate the same object fields in the same way, add a method or focused helper
  for that named state transition unless direct field access is clearer and has
  no shared invariant, defaulting, validation, or future-change risk.
- Public contracts and operational behavior changed only where intended.
- Relevant tests, formatters, linters, or builds were run when practical.
