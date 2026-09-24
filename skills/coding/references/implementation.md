# Implementation Workflow

Use this reference when writing new code or extending behavior.
Decision rules live in `./principles.md`. This file is the implementation steps.

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

Apply Priority in `./principles.md`.

## Implementation Steps

- Make the smallest coherent change that fully handles the request.
- Apply Subtract before you add, then shape what remains with Boundaries earn
  their cost, both in `./principles.md`.
- Design from existing boundaries: keep domain logic, adapters, UI state,
  persistence, and transport concerns where the project already places them.
- Keep process entrypoints for startup, dependency wiring, shutdown, and top
  level errors. Do not make them the default home for feature code.
- Do not extend an incoherent local layout when new code has a clean home under
  the language and framework convention. Do not move unrelated code merely to
  make the tree look tidy. Placement details are in `./file-organization.md`.
- Reuse an existing cohesive struct, options, context, or config object before
  creating a new parameter container.
- For a library or important module API, write or inspect a realistic usage
  example before settling the shape. Let consumer workflows drive the internal
  data structures when needed.
- When the change would duplicate a decision, apply Fix duplication while the
  change is open in `./principles.md`. Combine it in this change. Do not leave
  it for a later refactor pass.
- When a write can succeed and a later step can still fail, apply A write that
  stuck stays stuck in `./principles.md`.
- Choose the kind of check from `./verification.md`. What a test protects and
  how it proves that is Tests in `./principles.md`. Whether to add or remove
  it is the `test-audit` skill.
- Update docs, comments, fixtures, generated inputs, or examples only when they
  are part of the changed behavior.

## Completion Check

- Apply Review before done in `./principles.md`.
- The code follows systemic local patterns where they exist, and `./standards.md`
  where they do not.
- Public contracts and operational behavior changed only where intended.
- The expected observable outcome was verified directly. See Prove the outcome
  in `./principles.md`.
- What automated checks protect follows Tests in `./principles.md`. Whether
  they are added follows the `test-audit` skill.
- Report the direct proof, supporting validation, and any real path that could
  not run.
