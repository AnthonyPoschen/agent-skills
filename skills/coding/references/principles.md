# Principles

Shared decision rules for implementation, refactoring, review, debugging, and
design. Task references add their own steps and point here. Do not restate
these rules in those files.

## Priority

When rules conflict, use this order:

1. Correct behavior, safety, and the user's explicit requirements.
2. The project formatter, linter, tests, build system, framework conventions,
   and systemic local patterns.
3. This file, then `./standards.md`.

A local pattern is systemic when it appears across multiple nearby files, is
enforced by tooling, or is part of a clear framework or application convention.
For file and folder placement it must also be coherent, discoverable, and
compatible with the language and framework. Repetition alone does not make a
dumping ground authoritative. Do not copy one-off weak code just because it is
adjacent.

When no systemic project pattern exists, use `./standards.md`.

## Contracts

Preserve public behavior, APIs, data contracts, migrations, and operational
semantics unless the user explicitly asks to change them. That includes public
names, serialization, configuration, CLI flags, environment variables, URLs,
and fixture contracts.

## Scope

Keep edits inside the task. Do not normalize unrelated legacy code, and do not
rewrite a local legacy pattern outside the path this change needs.

Inside that path, apply Fix duplication while the change is open. Cleanup of
code this change already reaches is part of the task. A separate cleanup
project is not.

## Subtract before you add

Before adding non-trivial structure, simplify the touched area to the smallest
shape that meets the request.

- Remove dead paths, pass-throughs, duplicated decisions, stale compatibility
  code, misleading state, unnecessary conversions, and redundant validation
  when their absence can be verified within the task.
- Do not add an option, parser, validator, guard, persistence path, retry,
  abstraction, or extension point for a case the request and observed usage do
  not require.
- Cut unused scope before polishing names or structure.
- Stay inside the touched boundary. Report wider dead weight instead of turning
  the task into an unasked cleanup.

## Boundaries earn their cost

Add or reshape a module, type, or helper when it gives callers a simpler,
self-contained operation or owns knowledge they should not carry: an invariant,
lifecycle, policy, representation conversion, or integration detail.

A boundary must make the real caller's job shorter, clearer, or safer. One
caller is enough when the boundary owns meaningful complexity. Several callers
are not enough when the helper only moves the same reasoning.

Keep a cohesive local flow direct when the caller already owns the decisions
and can understand the work in place. Do not add a wrapper, interface, or
module that only renames an underlying API, forwards the same types and
arguments, or exists for a hypothetical second implementation.

Do not extract a helper only to remove repeated syntax. Similar statements can
be clearer than a generic helper. Consolidate when a named operation owns the
shared semantics and makes callers simpler. Keep orchestration and
caller-owned policy with the caller.

## Fix duplication while the change is open

When the work shows that a decision, conversion, sequencing rule, or state
transition would be copied, give that work one owner and migrate the direct
consumers this change reaches. Do that in the same change. The context is
already loaded. A later refactor pass has to relearn the area.

Use the shape you would have chosen if the new requirement had existed from
the start, when that shape is clear and those direct consumers can be migrated
and verified within this task. Delete the pass-through, duplicated decision,
awkward conversion, or needless state that the new path would otherwise
preserve.

If the same fix reaches code this task does not touch, report the opportunity.
Do not start a partial migration, and do not add a special-case escape hatch or
compatibility wrapper that entrenches the old boundary.

## A write that stuck stays stuck

If a function has already committed work, a later failure must not tell the
caller that nothing happened. Return the committed result with a distinct
error, or commit the follow-up in the same transaction. Retrying the caller's
original request must not duplicate the committed work.

## Prove the outcome

Before declaring a change complete, state its expected observable outcome and
verify that outcome directly.

- When the changed source file is itself the delivered artifact, inspect the
  resulting file or diff.
- For a deterministic transformation, run it with representative input and
  inspect the returned or rendered value.
- For dynamic behavior, run the service, command, application, or UI and
  inspect the result through its normal path.
- For persistence or side effects, inspect the actual row, file, message, or
  other state the operation produced.

A formatter, linter, build, or passing test is supporting evidence unless it
observes that outcome. Report the direct proof separately from supporting
checks and any path that could not run.

## Tests

A worthwhile test protects a stable contract through a public seam: a product
rule, public API, data or serialization guarantee, security property, or
correctness invariant. It needs a trustworthy oracle and should survive a
reasonable refactor.

Prove it at that seam. For a bug, watch the test fail on the unfixed behavior
for the intended reason, then confirm the fix makes it pass. Exercise the
contract with an executable check: a rendered flow, an HTTP request, a command,
or a deterministic exported operation.

Reading a source file to assert a string, selector, function name, route
literal, CSS declaration, or markup block, including a regex over HTML, CSS,
JavaScript, Go, or configuration, is a source-text check. That proves behavior
only when the text itself is the contract, such as a generated artifact or
another deliberately stable text contract.

Whether to add, keep, or remove a test is the `test-audit` skill.

## Review before done

Before declaring a code change complete, judge the diff with the review stance
in `./review.md`. Bugs, regressions, security issues, operational failures,
data loss, broken contracts, and unproven changed behavior outrank style. Fix
those findings before finishing. Use the output format in `./review.md` only
when the user asked for a review.
