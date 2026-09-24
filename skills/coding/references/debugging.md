# Debugging Workflow

Use this reference for failing tests, runtime errors, regressions, CI failures,
and bug fixes.

## Debugging Process

- Reproduce the failure through the real affected path, or identify the closest
  reliable signal when that is not possible.
- Read the failing output, stack trace, logs, or assertion before editing.
- Trace from symptom to boundary: input, state, dependency, transformation,
  output, and side effects.
- Make the smallest fix that addresses the root cause.
- Rerun the same scenario after the fix and observe the expected outcome.
- If the bug can recur through a public seam, admit one regression at that
  seam through the `test-audit` skill. Prove it as Tests in `./principles.md`
  describes: watch it fail, then confirm the fix makes it pass.
- For persistence, framework wiring, or service behavior, prefer a real local
  dependency and the normal application path over a stub that cannot prove the
  failure or fix. Load `./verification.md` when choosing the evidence.

## Fix Rules

- Do not silence errors or weaken assertions to make a failure disappear.
  Deleting a test for that reason fails the `test-audit` skill.
- Do not paper over races, nil/null cases, parse failures, or missing data with
  broad catch-all behavior unless that is the correct product behavior.
- If the bug is the public contract, change it only with the user's approval.
  Otherwise apply Contracts in `./principles.md`.
- When the failure happens after a successful write, apply A write that stuck
  stays stuck in `./principles.md`.
- Keep diagnostics useful: errors should include enough context to troubleshoot
  without leaking secrets.

## Completion Check

- The observed failure is explained.
- The fix targets the root cause, not just the symptom.
- The original scenario was rerun and the fixed outcome was observed directly.
- A bug that can recur through a public seam has one regression at that seam,
  admitted by the `test-audit` skill. A test that skill rejects was not added.
- Supporting verification was run or any blocker is reported.
