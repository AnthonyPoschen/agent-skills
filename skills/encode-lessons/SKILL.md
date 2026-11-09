---
name: encode-lessons
description: >
  Use when a correction, failure, workaround, or instruction recurs, or when a
  costly mistake could be prevented structurally. Encode the lesson in the
  smallest durable guard instead of relying on someone to remember it. Capture
  one lesson, or refresh lessons for staleness.
---

# Encode Lessons in Structure

A repeated correction is evidence that attention is the wrong control. When a
small durable mechanism can prevent the same mistake, prefer it to another
reminder.

## Choose Capture Or Refresh

Do one job per run.

- **Capture** records the lesson of one solved, verified problem. This is the
  default after a fix, review finding, or costly miss. A session with several
  lessons gets several sequential captures. Do not batch them into one write.
- **Refresh** audits lessons that already exist (written learnings, standards,
  lint rules, and tests that exist only to hold a lesson) against the current
  tree. Use it when asked to clean up, dedupe, or check lessons for drift, or
  when a named set of lessons may be stale. Refresh does not hunt for new
  lessons. It updates or removes what no longer earns its place.

An existing lesson that would now mislead is a capture of a correction: update
that lesson instead of adding a second one.

## Decide Whether The Lesson Is Durable

Treat a human correction, preventable failure, repeated workaround, or missed
instruction as a learning signal. Classify it before changing the system:

- **One-off:** fix the immediate problem. Do not create a permanent rule for a
  minor isolated event.
- **Recurring or costly:** find the smallest guard that prevents the pattern.
- **Requires judgment:** keep a clear instruction. Make the failure mode and
  the decision it requires obvious.

Do not confuse a vague preference with a proven pattern. A single high-risk or
high-cost mistake can still justify a durable guard.

Run this check before writing anything, and again after a stronger form
lands. Keep the lesson only when the finished change still hides the reason:
a reader of the code, tests, and docs would miss why the wrong path was taken,
and would be likely to make the same mistake or repeat the investigation. A
large diff, a long session, or a finished task is not enough. If the check
fails, write nothing and say why.

## Pick The Cheapest Durable Form

Prefer a mechanism that prevents the wrong path over text that asks someone not
to take it. Take the first form that holds:

1. **Test.** A test can fail on the defect's shape, runs where this repo
   already runs tests, and its failure tells the reader what to fix.
2. **Lint rule.** A lint, type, or CI check the repo already runs can reject
   the pattern. Tighten an existing rule before adding one. A new custom rule
   counts only where the repo already has a place to write rules.
3. **Standard.** No check can hold it, but an author must avoid it while
   writing and a reviewer can see it in a diff. Add one condition to the
   nearest existing standards or instruction file, in that file's style. Do
   not create a standards file that the repo does not have.
4. **Written learning.** Use this only for reasoning a test, lint rule, or
   standard cannot carry, such as a measured fact or a rejected alternative.
   Point at the stronger form when one exists.

When the current change can make the invalid state impossible through data
shape, types, ownership, or one canonical helper, do that in the layer that
owns the failure. That outranks a written learning, and it outranks a test
that only restates an impossible state. Choose the smallest mechanism that is
reliable for the actual risk. Do not build a framework, tool, or policy layer
for a minor one-off preference.

If a test, lint rule, or standard already states the lesson, the counterfactual
has failed. Do not add a written learning that repeats it.

## Close The Loop

When a durable guard is justified:

1. Put it in the layer that owns the failure.
2. Prove the guard catches the bad path or makes it impossible.
3. Remove or shorten any redundant reminder, while keeping useful rationale,
   exceptions, and user-facing documentation.
4. Leave the surrounding code, tooling, or skill easier to follow than before.

On refresh, re-apply the counterfactual to each lesson in scope. Delete or
shorten a lesson the current code, tests, or docs already make obvious. Update
a lesson that contradicts the tree. Leave a lesson that still passes the bar.
Report each keep, update, and removal in a few lines. Do not open unrelated
rewrites.

When the guard reaches beyond the current task's natural scope, report the
specific pattern and proposed guard. Do not start an unrelated framework or
repository-wide migration without approval.

## Avoid Weak Learning

- Do not promise to "remember" a recurring correction without recording or
  enforcing it.
- Do not fix one instance while leaving a known repeatable cause untouched.
- Do not record a lesson without giving it a concrete owner and next action.
- Do not preserve explanatory text that only repeats what a reliable mechanism
  now guarantees.
- Do not capture a lesson the final code, tests, and docs already explain.
- Do not refresh by appending a new summary of lessons that are still accurate.

## Source

Adapted from Pstack's
[Encode Lessons in Structure principle](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-encode-lessons-in-structure/SKILL.md).

The durable-bar counterfactual, capture versus refresh, and cheapest-form
order are adapted from EveryInc/compound-engineering-plugin
`skills/ce-compound@8e5d1c87613f484a5bdfeacedf86a7aaaccfa264`, MIT.

## Related

- A comment that states a contract is a lesson waiting for a guard. The
  `coding` skill's comment rules and the [no-comments reference](../coding/references/no-comments.md) send those here.
- `scripted-verification` when the repeated work is judging results.
