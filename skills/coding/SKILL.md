---
name: coding
description: >
  ALWAYS activate for programming-language work and code-adjacent project work:
  implementation, refactoring, debugging, code review, tests, CI fixes,
  package/build files, scripts, migrations, generated-code inputs,
  infrastructure-as-code, PR/diff discussion, and software architecture or
  design planning. Use this skill as the coding entrypoint router, then load
  principles.md and only the task-specific reference files needed for the
  current request.
---

# Coding

Use this skill as the entrypoint for software work. Keep this file in context,
then load the smallest relevant reference set for the task.

## Routing

- For any implementation, refactor, review, debug, or design task, read
  `./references/principles.md`. It owns the shared decision rules. The other
  references add task steps and point back there.
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
- For a consequential new data model, shared scaffold, or concurrent state,
  also read `./references/foundations.md`.
- For a new application or subsystem, source file placement, package or module
  placement, file moves, file splits, or layout discussion, read
  `./references/file-organization.md`.
- For Go code, application layout, package placement, or Go file organization,
  also read `./references/go.md`.
- For choosing tests or verification, especially for persistence, framework,
  service, CLI, or UI behavior, read `./references/verification.md`.
- For shared code-quality rules, read `./references/standards.md` whenever
  writing or changing code, and as needed during review/design.
- For assertion strategy (invariant checks, debug-vs-production assert
  decisions, and fail-fast contracts), read `./references/assertions.md` when
  writing or reviewing code that validates internal assumptions.

If a loaded reference describes ordered phases, run every remaining phase until
its completion check. Load the next reference when that phase starts. Do not
stop after the first phase, including subtraction.

Do not commit between phases unless the user asked to land the work as
separate commits. When they did, commit each coherent unit with
`git-commit-workflow`, then continue. Otherwise finish the remaining phases in
the working tree.

## Preserve The User's Mental Model

When code supports a user-facing workflow, model the implementation around the
user's intent rather than mirroring every field, option, or state in the backing
data model. Prefer one cohesive operation that expresses the user's decision
over several controls or helpers that expose internal structure. Keep advanced
state available behind the interaction that needs it, and make visible state
reversible and consistent with the user's expectations.

Before adding a control, state variable, or abstraction, identify the user
decision it owns. If an existing operation can express that decision directly,
extend it instead of creating a parallel path. Verify the resulting workflow at
its public boundary, including realistic content and the states that change the
user's next action.

## Example Lookup

Use examples only when a rule decision is unclear.

- `./examples/README.md` indexes focused examples.
- Open one matching example file at a time.
