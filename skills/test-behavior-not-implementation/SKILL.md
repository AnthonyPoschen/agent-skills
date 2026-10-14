---
name: test-behavior-not-implementation
description: "Apply when you write, change, or keep a test. Call the code as its users do and assert the exact result they see. Rewrite or delete a test that would pass if every imported function returned undefined."
---

# Test Behavior, Not Implementation

This skill covers how a test asserts. Whether a test should exist at all, and
tautology tests that copy the code under test, belong to the `test-audit` skill
and the coding skill's [Tests principle](../coding/references/principles.md#tests).
Apply those; this skill does not restate them.

A test calls the code as its users do and asserts the exact result they see. Rewrite or delete any test that would still pass if every imported function returned `undefined` or a wrong value.

These shapes fail that check:

- Loose or no assertion: `toBeDefined`, `toBeTruthy`, `not.toThrow`, `toBeGreaterThan(0)`.
- Call or absence only: `toHaveBeenCalled`, `not.toHaveBeenCalled`, `toBeUndefined`, `toEqual([])`.
- Restated constant, config default, table row, or prompt string.
- Assertions on data the test built, with no call to the subject.

Fix: call the subject in the test body with one concrete input and assert the literal output or effect. For a mock, assert the payload it received. For an absence, also assert the presence on another input. For a constant, test the mechanism that reads it with one input instead of restating the value.

Keep relation checks across table rows (a key in two tables, a parent that exists) and `*.test-d.ts` type checks.

## Source

Adapted rewrite of Pstack's
[`principle-test-behavior-not-implementation` skill](https://github.com/cursor/plugins/blob/1e5c90a4b6f249917bd21fe361b3a2a7deb741b8/pstack/skills/principle-test-behavior-not-implementation/SKILL.md)
(cursor/plugins commit `1e5c90a`), by poteto. The upstream licence is
unconfirmed. Anthony approved adopting this adaptation on 2026-10-14.
