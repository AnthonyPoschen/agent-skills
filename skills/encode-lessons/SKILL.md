---
name: encode-lessons
description: >
  Use when a correction, failure, workaround, or instruction recurs, or when a
  costly mistake could be prevented structurally. Encode the lesson in the
  smallest durable guard instead of relying on someone to remember it.
---

# Encode Lessons in Structure

A repeated correction is evidence that attention is the wrong control. When a
small durable mechanism can prevent the same mistake, prefer it to another
reminder.

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

## Choose The Strongest Suitable Guard

Prefer a mechanism that prevents the wrong path over text that asks someone not
to take it.

1. Make an invalid state impossible through data shape, types, or ownership.
2. Add a static check, lint rule, CI validation, or banned API.
3. Provide a canonical helper, template, generated structure, or safe default.
4. Validate at runtime when input is external or the condition is dynamic.
5. Use a concise rule and a concrete failure example when the decision cannot be
   enforced without judgment.

Choose the smallest mechanism that is reliable for the actual risk. Do not
build a framework, tool, or policy layer for a minor one-off preference.

## Close The Loop

When a durable guard is justified:

1. Put it in the layer that owns the failure.
2. Prove the guard catches the bad path or makes it impossible.
3. Remove or shorten any redundant reminder, while keeping useful rationale,
   exceptions, and user-facing documentation.
4. Leave the surrounding code, tooling, or skill easier to follow than before.

When the guard reaches beyond the current task's natural scope, report the
specific pattern and proposed guard. Do not start an unrelated framework or
repository-wide migration without approval.

## Avoid Weak Learning

- Do not promise to “remember” a recurring correction without recording or
  enforcing it.
- Do not fix one instance while leaving a known repeatable cause untouched.
- Do not record a lesson without giving it a concrete owner and next action.
- Do not preserve explanatory text that only repeats what a reliable mechanism
  now guarantees.

## Source

Adapted from Pstack's
[Encode Lessons in Structure principle](https://github.com/cursor/plugins/blob/main/pstack/skills/principle-encode-lessons-in-structure/SKILL.md).
