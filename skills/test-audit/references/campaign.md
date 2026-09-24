# Campaign

Use this when the audit covers a subsystem, a package, or a budget the user
set. Retention, the evidence record, validation, and the budget rule in
[SKILL.md](../SKILL.md) apply to every step. This file is the order of work.
Finish a step before starting the next.

## 1. Baseline

Record the in-scope test and test-support line counts, and the pass or fail
state of every in-scope test file, at a pinned revision. Keep failures in
their own list. A baseline failure may be a real bug, not a stale test.

Done when every in-scope test file has a recorded result.

## 2. Lanes

Split the surface into lanes along production owners, not file-name prefixes.
Include the subsystem's cases at shared boundaries, and any QA or end-to-end
harness it owns.

Done when every in-scope test file belongs to exactly one lane.

## 3. Ledger

Read every assigned test in full, including table rows. Read the production
owners, their entry points, callers, history, and CI routing. Mark each test
declaration. Mark a table once unless its rows need different marks.

- `R`: retain. Name the contract and the bug it catches. A move to a clearer
  file stays `R`.
- `F`: retain the contract, but repair the assertion. A negative that passes
  when only one of several conditions is missing is this mark.
- `C`: consolidate. Name the owner that absorbs the assertion: a sibling table
  row, a stronger boundary suite, or a shared owner elsewhere.
- `D`: delete. Name the proof that remains, or why no contract exists.

Judge a test by its assertions, not its name.

Done when every declaration has a mark and an evidence line.

## 4. Layer plan

Treat the ledger as input, not as the edit list. Look for a redundant layer:
several suites replaying one owner through a mock, beside a stronger suite at
the real boundary. Name the keeper for each contract. Prefer the real
boundary with a fake external edge over a mocked collaborator. Correct ledger
mistakes this pass finds.

Done when each lane names its retired files, its keeper per contract, the
assertions to carry into those keepers, and the test-only production seams the
cut unlocks.

## 5. Cutover

Edit one lane at a time. Serialize edits to shared harnesses and support files
through one owner. With each lane, remove the test-only seams it unlocks.
Register a moved suite where CI or a test inventory lists files.

Done when every lane plan is applied and each lane's keepers pass.

## 6. Preservation review

Before calling the campaign done, compare deleted coverage with the keepers.
Look for a contract that lost its only proof, and for a new assertion that
cannot fail.

For each restored contract, mutate the production owner so the keeper must
fail, confirm that it fails, then restore the owner byte for byte.

Done when every reported gap is restored or rejected with source evidence, and
every restored contract has a caught mutation.

## 7. Product defects

A baseline failure that survives into a keeper is a bug. Fix it at its owner
as separate work from the deletion, and prove it through the real flow. Run a
control that reverts the fix and shows the old behavior. Record an unrelated
product discrepancy as a follow-up.

Done when each repaired defect has a failing control and a passing candidate
on the same harness.

## 8. Hand off

When the campaign has been open across upstream changes, bring those changes
in and port any new contract into the keeper. Confirm a regression added
upstream still has a home. Rerun the subsystem's checks on the merged head.

Report with the handoff in [SKILL.md](../SKILL.md), plus:

- Baseline and final test and test-support line counts, with production counted
  separately.
- Lanes, retired layers, and keepers.
- Preservation gaps and the mutations that caught them.
- Product defects, with control and candidate proof.
