---
name: coding
description: >
  ALWAYS activate for programming-language work and code-adjacent project work:
  implementation, refactoring, debugging, code review, tests, CI fixes,
  package/build files, scripts, migrations, generated-code inputs,
  infrastructure-as-code, PR/diff discussion, and software architecture or
  design planning. Use this skill as the coding entrypoint router, then load
  only the task-specific reference files needed for the current request.
---

# Coding

Use this skill as the entrypoint for software work. Keep this file in context,
then load the smallest relevant reference set for the task.

## Routing

- For writing new code or extending behavior, read
  `./references/implementation.md`.
- For refactoring, cleanup, simplification, extraction, restructuring, or
  rewrite work, read `./references/refactoring.md`.
- For code review, PR review, diff critique, or risk assessment, read
  `./references/review.md`.
- For debugging, failing tests, regressions, runtime errors, CI failures, or
  bug fixes, read `./references/debugging.md`.
- For architecture, API design, module boundaries, dependency direction, or
  design-only discussion, read `./references/design.md`.
- For a new application or subsystem, source file placement, package or module
  placement, file moves, file splits, or layout discussion, read
  `./references/file-organization.md`.
- For Go application layout, package placement, or Go file organization, also
  read `./references/go.md`.
- For choosing tests or verification, especially for persistence, framework,
  service, CLI, or UI behavior, read `./references/verification.md`.
- For shared code-quality rules, read `./references/standards.md` whenever
  writing or changing code, and as needed during review/design.
- For assertion strategy (invariant checks, debug-vs-production assert
  decisions, and fail-fast contracts), read `./references/assertions.md` when
  writing or reviewing code that validates internal assumptions.

Load only the references needed for the user's current request. If the task has
multiple phases, load the next reference when that phase starts.

## Baseline Rules

- Correctness, safety, and user intent take priority over style preferences.
- Project formatter, linter, tests, build system, framework conventions, and
  systemic local patterns take priority over this skill's default standards.
- A local pattern is systemic when it appears across multiple nearby files, is
  enforced by tooling, or is part of a clear framework/application convention.
  For file and folder placement, it must also be coherent, discoverable, and
  compatible with the language and framework. Repetition alone does not make a
  dumping ground authoritative. Do not copy one-off weak code just because it is
  adjacent.
- When no systemic project pattern exists, use `./references/standards.md`.
- Preserve public behavior, APIs, data contracts, migrations, and operational
  semantics unless the user explicitly asks to change them.
- Keep edits scoped to the task. Improve touched code enough for a coherent
  result, but do not normalize unrelated legacy code.

## Example Lookup

Use examples only when a rule decision is unclear.

- `./examples/README.md` indexes focused examples.
- Open one matching example file at a time.
