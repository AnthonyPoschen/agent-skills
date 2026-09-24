---
name: test-audit
description: >
  Decide whether a proposed or existing test should exist. Activate alongside
  the coding skill whenever writing, changing, reviewing, or removing tests,
  including a test an agent is about to add for a small edit. Gate new tests
  before they land. Audit low-value, duplicative, implementation-coupled, or
  assertion-free tests and the test-only production seams they require. Also
  use when asked to prune, sweep, delete, or cut tests, or to remove useless
  tests while keeping real behavior covered.
---

# Test Audit

Decide whether a test should exist. Run this check before adding a test, when
changing or reviewing one, and when sweeping tests that are already there. It
runs with the coding skill.

The coding skill owns what a worthwhile test protects and how to prove it: the
Tests section of its `references/principles.md`, the kind of check in
`references/verification.md`, and invariant checks in `references/assertions.md`.
Read those when they are not already loaded. Do not restate them here.

Use the authoring gate for a test this task is about to add or rewrite. Use a
focused audit for a file, a diff, or a handful of suspects. For a subsystem,
package, or a budget the user set, read [references/campaign.md](references/campaign.md)
before editing.

## Authoring gate

Before adding a test, or landing a change that rewrites what a test claims,
answer four questions. A missing answer means do not add it yet.

1. What observable behavior, invariant, or independent contract does it protect?
2. What credible regression makes it fail?
3. Why does an existing test not already catch that failure? One contract has
   one primary test at the strongest boundary. Another layer needs a distinct
   risk the owner cannot reach, such as transport or lifecycle. Extend a table
   or shared fixture instead of adding a near-duplicate.
4. Does it need a production seam that no production caller needs, such as an
   export, flag, wrapper, or injection hook? If it does, test the real boundary
   instead.

Then match the test against every junk pattern. A match fails the gate unless
retention names the contract it guards on its own.

A test that breaks under a behavior-preserving refactor asserts implementation.
Rewrite it at the owning boundary, or drop it.

A bug that can recur through a public seam gets one regression at that seam.
Prove it the way the coding skill's Tests section describes. Reject a
regression that never failed. It proves the mock. Do not replay the same
scenario at every layer it crosses.

When rejecting a proposed test, name the failed question or junk pattern and
the existing proof, if any.

## Junk patterns

Reject a new test that matches one of these. On an audit, these are the
candidates to inspect.

- An assertion-free coverage probe.
- A self-comparison, or a copier that checks a value against itself.
- A copied fixture, inventory, manifest, or export list.
- A source-text check, unless the coding skill's Tests section treats that text
  as the contract.
- A private predicate or call-shape test duplicated at the real boundary.
- A second invocation of the same contract.
- A provider-local replay of a shared helper.
- A test whose only job is to preserve a test-only export, global, or wrapper.
- Production code whose only callers are tests.
- An expected value produced by the helper or renderer under test.
- A mock that implements the asserted behavior, or one identical mock standing
  in for different APIs.
- A fixture that supplies the receipt, admission, or callback order the owner
  should produce. A persistence assertion against a store the path never writes.
- A capability test that restates a declared flag instead of exercising the
  delivery or acknowledgement that flag promises.
- A negative control that passes for an unrelated reason, such as a denial from
  a different guard, or a rejection the production path never reaches.
- A name or fixture that promises more than the input exercises.

## Retention

Keep the primary proof of a contract from the coding skill's Tests section.
Also keep:

- Ordering, when a caller can observe that the order is wrong.
- A regression with a credible failure mode that the owner-boundary proof does
  not already catch.
- A test that fails on the current baseline. Treat that as a possible product
  bug. Reproduce it and repair the owner. Do not delete it to make the run
  green.

Keep a static or slow test when it is still the proof of a contract. A
coverage percentage does not decide admission. A test that resembles
implementation may still be the only proof of the contract. Show the stronger
remaining proof before removing it.

A touched test that still protects its contract stays. Update it when that
contract still holds. Remove or replace a test that protects no contract.

## Budget

When the user sets a budget, such as a fraction of a subsystem, continue in
coherent batches until the budget is met or the remaining candidates fail the
evidence bar. Stopping after the first obvious file leaves the rest of the
budget undone. Use only a budget the user set. Report coverage when the
project already measures it. The evidence record decides whether a contract
survived.

## Focused audit

Stay read-only until each candidate has the evidence below. Prefer a few
high-confidence candidates over a speculative inventory. Hunt the junk patterns.

Read the whole test and its production owner before judging it: entry point,
callers, callees, siblings, overlapping tests, how CI runs them, and the
history that explains why the test or seam exists. When the test claims
dependency-backed behavior, read that dependency's source or types.

Record every field before editing. A missing field means the candidate is not
ready to delete.

- Exact test name and location.
- The failure it can actually detect.
- Non-test callers of the production or support seam it covers.
- The stronger owner-boundary proof that remains, or why no proof is needed.
- The history and the reason the test or seam exists.
- The production or test-support deletion this unlocks.
- The risk, and the focused command that checks it.

## Edit

Change one coherent owner-boundary batch. Delete the test-only exports,
globals, wrappers, and dead production paths that batch unlocks. Move a
retained regression to its canonical owner. Do not add a replacement that
restates the same implementation. Do not turn an uncertain candidate into a
deletion to raise the count.

Commit, push, or open a PR only when the user asked. Use
`git-commit-workflow` when they did.

## Validation

Use the project's own check for the owner and its siblings. When a source-text
assertion goes away, run the executable check that owns the real contract.
Inspect the diff and report production and tooling changes separately from
tests and test support.

Do not edit files a test run is already writing.

## Handoff

Report:

- The low-value categories removed, and why they were low value.
- Production seams simplified.
- Candidates kept, and why they still earn their cost.
- The checks actually run, and any check that could not run.
- Production lines changed versus test and test-support lines.
- Commit or PR state, when the user asked to land the work.
- Named follow-ups.

## Source

Adapted from OpenClaw's
[test-audit skill](https://github.com/openclaw/openclaw/blob/main/.agents/skills/test-audit/SKILL.md)
(Copyright (c) 2026 OpenClaw Foundation), used under the MIT License.
Keep this repository's admission rules. Use the project's own test command,
layout, and git workflow.
