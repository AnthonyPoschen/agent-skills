---
name: coding
description: >
  ALWAYS activate for coding, refactor, review, debugging, and design tasks
  across programming languages. High-priority guardrails for clean, idiomatic,
  maintainable code while respecting project tools and language conventions.
---

## Activation

- MUST activate for programming-language tasks in any mode: implementation,
  coding, debugging, refactor/restructure, code review, PR/diff discussion,
  and planning/chat on architecture or design decisions.
- Hard triggers:
  - target files are source code, tests, build files, package manifests,
    scripts, migrations, generated-code inputs, or infrastructure-as-code
  - PR/diff includes code or code-adjacent project files
  - explicit request for coding, debugging, refactor, review, implementation,
    architecture, or design of software behavior
- Intent triggers (including plan/chat mode): refactor, cleanup, simplify,
  split, extract, modularize, review, critique, improve, code quality,
  debugging, tests, CI, PR discussion, diff discussion, architecture/design
  discussion.
- Do not activate for generic discussion with no code/design/review intent,
  incidental word matches, prose-only docs, data-only files, or unrelated
  non-code scope.
- Activation priority: direct file/diff evidence -> explicit intent
  (refactor/review/plan/chat) -> project/repo context signals.

## Scope

- Applies to new code, modified regions, and agent-added helpers.
- No mandatory full-file rewrites.
- Preserve architecture, API behavior, and public contracts unless user requests
  a change.

## Priority And Conflict Rules

- Conflict order: correctness/safety -> project settings/conventions -> this
  skill.
- Rule strength: `MUST`/`MUST NOT` -> `SHOULD` -> `MAY`.
- If formatter/linter/build/tooling conflicts with style preference, follow
  project tools while preserving correct behavior.

## Example Lookup Protocol

- Match task to rule category first.
- If unclear, open only matching example file.
- Do not load all example files by default.
- Examples index: `./examples/README.md`

## Rules

### Control Flow

- `CF-1 MUST`: Put simple fail-state checks first and return immediately.
  Check: guard clauses appear before main logic.
- `CF-2 MUST`: Avoid `if/else` when `if` branch already exits (`return`,
  `continue`, `break`). Check: prefer guard clauses + flat main path.
- `CF-3 SHOULD`: Keep independent fail checks as separate guards. Check:
  independent conditions are not collapsed unless correctness needs it.
- `CF-4 MUST`: Conditional nesting guidance applies to conditional chains only.
  Check: loop/switch nesting does not count.

Detailed examples: `./examples/control-flow.md`

### Formatting

- `FMT-1 MUST`: Project formatter, linter, and language conventions win over
  manual style preferences.
  Check: do not fight automatic formatting or established local style.
- `FMT-2 SHOULD`: Prefer function-body lines `<= 80` chars when project style
  and formatter allow it.
  Check: long statements/expressions are simplified or wrapped by formatter.
- `FMT-3 SHOULD`: Prefer comment lines `<= 80` chars when project style allows
  it. Check: docs, phase comments, and tagged comments wrap cleanly.
- `FMT-4 SHOULD`: Strongly prefer single-line function signatures and
  declarations when formatter and language conventions allow it.
  Check: avoid manual multiline signatures unless formatter, language syntax,
  or readability clearly requires them.
- `FMT-5 MUST NOT`: Split signatures only to satisfy manual width preferences.
  Check: if a signature is too large, first consider parameter cohesion,
  naming, or function shape.

### Functions

- `FN-1 SHOULD`: Function definitions/declarations stay single-line when
  formatter and language conventions allow it.
  Check: signature line not split without a formatter, syntax, or readability
  reason.
- `FN-1b SHOULD`: Avoid multiline function parameter lists when the project
  formatter does not require them.
  Check: prefer improving parameter cohesion over manual signature wrapping.
- `FN-2 MUST`: If parameter count is `>= 5`, justify cohesion or refactor
  before accepting the function shape.
  Check: either (a) cohesive inputs are encapsulated, or (b) function is split
  to reduce unrelated inputs.
- `FN-2b SHOULD`: Record brief rationale when keeping `>= 5` parameters.
  Check: reason explains why inputs are independent and cannot be compacted.
- `FN-3 MUST`: Prefer reusing existing cohesive structs/types/contexts before
  introducing new input structs. Check: existing types are considered first.
- `FN-4 SHOULD`: If compaction is not obvious, evaluate whether function is too
  complex or truly needs many independent inputs.
  Check: decision recorded in code shape (split vs keep inputs).
- `FN-5 SHOULD`: If function is too complex, split responsibilities into smaller
  functions. Check: parent orchestrates, children execute.
- `FN-6 SHOULD`: If many inputs are still required, create a focused input
  struct for related inputs only. Check: struct groups one cohesive subset.
- `FN-7 MUST NOT`: Do not bundle all parameters into one catch-all wrapper
  struct. Check: unrelated inputs remain separate.
- `FN-8 MUST NOT`: Do not pass broad unrelated context/struct only to reduce
  parameter count. Check: passed fields remain cohesive to function purpose.
- `FN-9 MUST`: Function responsibility aligns with function name and one primary
  purpose. Check: no unrelated mixed responsibilities.
- `FN-10 MUST`: Do not extract helper functions from functions `<= 40 code
  lines` for decomposition, phase labeling, or visual organization alone.
  Check: short functions stay inline unless extraction is required by reuse,
  duplicated substantial logic, API boundaries, or clear correctness/testing
  need.
- `FN-11 MUST`: For functions `> 40 code lines`, extract only substantial
  phases or cohesive phase groups when doing so improves local reasoning,
  naming, reuse, or testability.
  Check: many short, linear phases may remain inline; avoid one helper per
  small phase.
- `FN-12 MUST`: Extracted helpers represent meaningful work, not just comments
  converted into function names.
  Check: helper name adds domain meaning and body is substantial, reusable,
  independently testable, or hides cohesive detail.
- `FN-13 SHOULD`: After extraction, parent function retains meaningful
  orchestration or policy.
  Check: if parent only calls helpers in sequence, reconsider whether parent
  should exist or caller should compose steps directly.
- `FN-14 MUST`: Parent keeps branching/decision logic; children execute steps.
  Check: high-level decision trees are not spread across children.
- `FN-15 SHOULD`: Pass resolved decision inputs to children (flags/enums/
  prepared values), not raw context requiring new policy decisions.
  Check: children do not re-derive high-level business decisions.
- `FN-16 MUST`: Children may validate local invariants but must not duplicate
  orchestration-level branching. Check: no duplicated high-level branching.
- `FN-17 SHOULD`: Extract substantial reusable logic into dedicated helper
  functions/types. Check: non-trivial repeated logic is not duplicated inline.

Detailed examples: `./examples/functions.md`

### Abstraction Boundaries

- `AB-1 MUST NOT`: Introduce pass-through wrappers that simply mirror external
  SDK/API methods with renamed plumbing.
  Check: wrapper adds clear behavior, not just argument forwarding.
- `AB-2 MUST`: Add wrapper/helper layers only when they add domain value
  (policy, retries, validation, translation, composition, or test seam).
  Check: abstraction has explicit purpose beyond hiding the underlying call.
- `AB-3 SHOULD`: Prefer calling vendor SDK/API directly at integration edges
  when wrapper adds no domain semantics.
  Check: simple one-to-one SDK calls stay direct in boundary adapters.
- `AB-4 MUST NOT`: Create internal objects solely to hide provider names or
  mimic provider operations one-for-one.
  Check: object methods are not thin aliases for provider methods.

### Variables And State

- `VS-1 MUST`: Do not introduce mutable global runtime state when value can be
  created at startup and injected. Check: services/clients/config created in
  `main` (or composition root) and passed in.
- `VS-2 MUST`: Prefer dependency injection over package-level singletons for
  runtime collaborators. Check: tests can substitute fakes/stubs without global
  patching.
- `VS-3 SHOULD`: Use smallest practical scope; widen only for shared read use.
  Check: variable scope is local unless broader scope clearly needed.
- `IMM-1 MUST`: If value never changes, declare immutable (`const` when
  available, else language-equivalent immutable pattern).
  Check: no reassignment and immutable form used where supported.

Detailed examples: `./examples/variables-magic-values.md`

### Magic Values

- `MV-1 MAY`: One-off literal allowed when obvious and used once in narrow
  scope. Check: literal appears once and meaning is clear.
- `MV-2 MUST`: Reused literal within scope/module must become named constant/
  variable. Check: repeated literals consolidated to one definition.
- `MV-3 MUST`: Reused literal across files/areas must be promoted to shared
  named constant at narrowest shared boundary. Check: cross-area duplicates use
  shared constant.
- `MV-4 MUST`: Shared constants may be global only when immutable domain
  constants with broad reuse. Check: global shared values are immutable and not
  runtime mutable config/state.

### Comments And Docs

- `CD-1 MUST`: Public/exported APIs, non-obvious behavior, domain policy,
  edge-case handling, and complex helpers have short doc comments.
  Check: comments explain purpose, contract, or consequence, not line-by-line
  implementation.
- `CD-1b SHOULD`: Private obvious helpers and small local functions usually do
  not need doc comments.
  Check: name, types, and local context make purpose clear without a comment.
- `CD-2 MUST`: Do not narrate obvious syntax/operations. Check: no trivial
  per-line commentary.
- `CD-3 MUST`: In non-trivial functions, add one short comment at start of each
  major phase block. Check: one comment covers block intent until next phase
  comment.
- `CD-4 SHOULD`: Phase comments read as coherent story on skim. Check:
  high-level flow clear from comments only.
- `CD-5 MUST`: Phase comments concise and action-oriented. Check: usually one
  sentence with verb-first phrasing.
- `CD-6 SHOULD`: Skip phase comments in tiny single-phase functions
  (`<= 10 code lines`, excluding blank/comment lines). Check: no forced comment
  noise in tiny functions.
- `CD-7 MUST`: Keep comments synchronized with code changes. Check: stale docs/
  phase comments updated in same change.
- `TAG-1 MUST`: Use exact uppercase prefixes: `TODO:`, `FIXME:`, `NOTE:`,
  `HACK:`, `WARNING:`. Check: exact prefix spelling/casing.
- `TAG-2 MUST`: Tagged comments may be multi-line but stay concise and obey
  comment width limits. Check: one logical note, wrapped to `<= 80`.
- `TAG-3 MUST`: `FIXME:`, `HACK:`, and `WARNING:` include consequence if
  ignored. Check: risk/failure mode/impact stated.
- `TAG-4 SHOULD`: `TODO:`, `FIXME:`, and `HACK:` include issue reference when
  available (`#123` or URL). Check: reference present when known.
- `TAG-5 SHOULD`: `HACK:` includes removal condition. Check: clear removal
  trigger documented.
- `TAG-6 MUST NOT`: Tagged comments must not narrate obvious code. Check: tags
  capture deferred work, risk, or non-obvious constraints.
- `TAG-7 MUST`: Remove or update stale tagged comments in touched regions.
  Check: tags in edited areas remain accurate.
- `TAG-8 SHOULD`: Keep tagged comments near relevant code location.
  Check: tags not dumped into unrelated header blocks.
- `TAG-9 MUST`: Put one blank line between tagged-comment blocks and generic
  comments in both directions. Check: no direct adjacency.

Detailed examples: `./examples/comments.md`

## Quick Checklist

- Activation matches programming-language work or code-adjacent project files.
- Project tools/conventions take precedence on conflict.
- Guard clauses first; avoid avoidable nested conditional chains.
- Prefer `<= 80` body/comment lines and single-line signatures when formatter
  and language conventions allow it.
- Do not split signatures only for manual width preferences.
- For `>= 5` params: justify cohesion or refactor before accepting shape; reuse
  existing structs first; split function if complex; use focused input struct
  only for related subsets.
- Avoid pass-through provider wrappers; keep direct SDK calls unless abstraction
  adds clear domain behavior.
- Keep short linear functions inline; decompose long functions by substantial
  phases or cohesive phase groups; parent retains meaningful orchestration.
- No mutable global runtime state; repeated literals consolidated; purposeful
  docs and tag-comment rules satisfied.
