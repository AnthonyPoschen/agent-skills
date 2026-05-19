# Debugging Workflow

Use this reference for failing tests, runtime errors, regressions, CI failures,
and bug fixes.

## Debugging Process

- Reproduce the failure or identify the closest available signal first.
- Read the failing output, stack trace, logs, or assertion before editing.
- Trace from symptom to boundary: input, state, dependency, transformation,
  output, and side effects.
- Make the smallest fix that addresses the root cause.
- Add or update a regression test when the bug could plausibly return.
- Run the narrowest relevant verification first, then broader checks if risk
  justifies them.

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
- Regression coverage or a clear reason for no test is documented.
- Relevant verification was run or any blocker is reported.
