# Review Workflow

Use this reference for code review, PR review, diff critique, and risk
assessment.

## Review Stance

Prioritize findings that could cause bugs, regressions, security issues,
operational failures, data loss, broken contracts, or unproven changed behavior.
Style-only comments are secondary unless they hide a real maintainability risk.

## Review Process

- Understand the intended behavior before judging the implementation.
- Inspect changed code, relevant callers/callees, tests, migrations, config, and
  generated artifacts when they affect risk.
- Compare the change to systemic local patterns and `./standards.md`.
- Check whether new boundaries reduce caller burden or merely add indirection.
  Keep cohesive local work direct; recommend a boundary only when it owns a
  coherent responsibility callers should not carry.
- Check edge cases, error paths, concurrency, data boundaries, and compatibility.
- Identify the changed behavior's expected observable outcome and inspect the
  direct proof that it occurred. A build, linter, or mocked test is insufficient
  when it does not observe the affected behavior.
- Raise a verification finding only when a meaningful changed behavior remains
  unproven and there is a concrete practical proof path. Do not request generic
  test coverage or tests that would pin incidental logs, private call sequences,
  or implementation details.
- When an automated check is appropriate, it should protect a stable product,
  API, data, security, or correctness contract with a trustworthy oracle.

## Output Format

Lead with findings, ordered by severity. Each finding should include file/line
evidence and the concrete impact. Keep summaries brief and after findings.

Use this structure:

```markdown
**Findings**

1. Severity: concise issue title.
   Evidence and impact.

**Open Questions**

- Question or assumption, if any.

**Summary**

Brief context only.
```

If there are no findings, say so clearly and mention residual risk or test gaps.
