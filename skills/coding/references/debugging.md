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
- Add an automated regression check only when it protects a stable, meaningful
  contract and has a trustworthy oracle. Do not add a test merely because the
  bug happened or the code could change again.
- For persistence, framework wiring, or service behavior, prefer a real local
  dependency and the normal application path over a stub that cannot prove the
  failure or fix. Load `./verification.md` when choosing the evidence.

## Fix Rules

- Do not silence errors, weaken assertions, or delete tests to make failures
  disappear.
- Do not paper over races, nil/null cases, parse failures, or missing data with
  broad catch-all behavior unless that is the correct product behavior.
- Preserve public contracts unless the bug is the contract itself and the user
  approves changing it.
- Keep diagnostics useful: errors should include enough context to troubleshoot
  without leaking secrets.

## Completion Check

- The observed failure is explained.
- The fix targets the root cause, not just the symptom.
- The original scenario was rerun and the fixed outcome was observed directly.
- Any automated regression check protects a stable contract rather than an
  incidental detail. If no such check is warranted, do not add one merely for
  coverage.
- Supporting verification was run or any blocker is reported.
