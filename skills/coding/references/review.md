# Review Workflow

Use this reference for code review, PR review, diff critique, and risk
assessment.

## Review Stance

Prioritize findings that could cause bugs, regressions, security issues,
operational failures, data loss, broken contracts, or missing critical tests.
Style-only comments are secondary unless they hide a real maintainability risk.

## Review Process

- Understand the intended behavior before judging the implementation.
- Inspect changed code, relevant callers/callees, tests, migrations, config, and
  generated artifacts when they affect risk.
- Compare the change to systemic local patterns and `./standards.md`.
- Look for shared-behavior opportunities that only became visible after the
  implementation: multiple functions with the same traversal, lookup, parsing,
  validation, dispatch, or edge-handling policy. Recommend extraction when those
  functions can be combined around one named behavior, not merely because code
  looks textually similar.
- Look for repeated external object mutation. When several callers set the same
  fields or apply the same configuration sequence, recommend a method or focused
  helper if it would centralize a named state transition, invariant, default,
  validation rule, or derived-field update. Do not require private fields when
  public fields are intentional; helpers can coexist with public state.
- Check edge cases, error paths, concurrency, data boundaries, and compatibility.
- Verify tests cover the changed behavior and likely failure modes.

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
