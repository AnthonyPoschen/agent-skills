# Implementation Workflow

Use this reference when writing new code or extending behavior.

## Before Editing

- Identify the user-visible behavior, API contract, and success criteria.
- Inspect nearby files for systemic project patterns before choosing a design.
- Check project tooling and existing tests so the implementation fits the repo.
- Identify the real caller and the knowledge it should not need to carry.
- Before creating a source file, package, or module, identify the language and
  framework convention and choose the path a reader is most likely to search.
  Load `./file-organization.md` when placement is a real decision.
- Prefer existing helpers, types, framework conventions, and composition roots
  when they make the caller's path clearer. Do not preserve an awkward local
  shape when the task is to improve a meaningful boundary.

## Pattern Priority

Follow this order:

1. Correct behavior, safety, and explicit user requirements.
2. Project tooling and enforced conventions.
3. Systemic local architecture and naming patterns.
4. Shared standards from `./standards.md`.

When local code is inconsistent, follow the pattern with the clearest tool,
framework, or repeated usage support. If no pattern is strong, use the shared
standards.

## Subtract Before You Add

Before adding non-trivial structure, simplify the touched path to the smallest
shape that can meet the request. A smaller base reveals the real responsibility
and keeps the next addition from inheriting accidental complexity.

- Remove dead paths, pass-throughs, stale compatibility code, misleading state,
  and redundant validation when their absence can be verified within the task.
- Do not add an option, parser, validator, guard, persistence path, retry, or
  abstraction for a speculative case the request and observed usage do not
  require.
- Cut unused scope before polishing the remaining design. Prefer a complete
  small workflow over a broader design that needs supporting machinery.
- Keep the cleanup inside the touched boundary. Report wider dead weight rather
  than turning an implementation task into an unasked-for cleanup project.

## Implementation Rules

- Make the smallest coherent change that fully handles the request.
- Design from existing boundaries: keep domain logic, adapters, UI state,
  persistence, and transport concerns where the project already places them.
- Keep process entrypoints for startup, dependency wiring, shutdown, and top
  level errors. Do not make them the default home for feature code.
- Do not extend an incoherent local layout when new code has a clean home under
  the language and framework convention. Do not move unrelated code merely to
  make the tree look tidy.
- Keep a cohesive local flow direct when its caller already owns the decisions
  and can understand the work in place.
- Add or reshape a module, type, or helper when it gives its real callers a
  simpler, self-contained way to do their job and owns knowledge they should
  not carry: an invariant, lifecycle, policy, representation conversion, or
  integration detail.
- Do not add a boundary that only renames an underlying API, forwards the same
  types and arguments, or moves the same reasoning into another file.
- For a library or important module API, write or inspect a realistic usage
  example before settling the shape. Let consumer workflows drive the internal
  data structures when needed.
- Prefer deletion, direct composition, and existing structures over speculative
  layers, extension points, validators, or configuration.
- Reuse existing cohesive structs, options, contexts, or config objects before
  creating new parameter containers.
- Choose verification based on the changed behavior; load
  `./verification.md` when deciding between unit tests, real dependencies, and
  mocks.
- Update docs, comments, fixtures, generated inputs, or examples only when they
  are part of the changed behavior.

## Improve The Touched Path

When a change touches an existing test, API, module boundary, or abstraction,
assess the directly involved path rather than preserving it by default.

- For a touched test, identify the stable contract it protects. Keep or update
  it when that contract remains meaningful; remove or replace it when it only
  asserts incidental logs, private call order, stale implementation shape, or
  mock choreography. Do not change or delete a test merely to hide a legitimate
  failure.
- When a caller needs conversions, sequencing knowledge, or a workaround just
  to use an existing boundary, inspect that boundary's direct consumers. Decide
  whether the friction is isolated or systemic before adding another escape
  hatch.
- When a meaningful new requirement exposes that friction, compare the direct
  change with the smallest shape that would make sense if the requirement had
  existed from the beginning. Use that shape when it clearly simplifies the
  affected callers and can be completed within the task's natural scope.
- If a revised boundary clearly makes its direct consumers simpler and they can
  be migrated and verified within the task's natural scope, improve the boundary
  and migrate them together. Do not leave a new hot-path workaround that
  entrenches a known bad interface.
- If the needed redesign reaches beyond the task's natural scope, keep the
  current change direct and report the broader opportunity. Do not begin a
  partial migration or add speculative compatibility machinery.

## Completion Check

- Before finishing, review the completed diff as if it were a code review:
  check correctness, public behavior, shared abstractions, duplicated internals,
  object ownership, tests, docs, and whether the change stayed scoped.
- The code follows systemic local patterns where they exist.
- The code falls back to `./standards.md` where local guidance is absent.
- Review every new boundary: it should make the caller's job clearer or own a
  coherent responsibility. Inline or remove any layer that does neither.
- Make an opportunistic refactor only when it simplifies the touched path by
  deleting a pass-through, duplicated decision, awkward conversion, or needless
  state. Do not extract merely because code is textually similar.
- Public contracts and operational behavior changed only where intended.
- Identify the expected observable outcome and verify it directly before
  declaring the work complete. A formatter, linter, build, or passing test is
  supporting evidence unless it observes that outcome itself.
- Add automated coverage only when it protects a stable, meaningful contract;
  do not pin incidental implementation details or diagnostics. Report the
  direct proof, supporting validation, and any real path that could not run.
