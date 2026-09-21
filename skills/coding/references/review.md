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
- Compare the change to systemic local patterns, `./principles.md`, and
  `./standards.md`.
- Check new boundaries against Boundaries earn their cost in `./principles.md`.
- Check whether new paths and filenames make code easy to find. Flag feature or
  infrastructure code that has accumulated in an entrypoint, or a patch that
  extends a known dumping ground when a clean conventional path was available.
- Do not flag a file for size alone. Flag it when it mixes parts with clear,
  separate names that a reader would look for separately. Do not ask for one
  file per type or method.
- Flag repeated path names only when they add no meaning. A canonical entry
  file, matching type and implementation pair, or role suffix can justify a
  repeated name.
- When a touched boundary causes conversions, sequencing knowledge, or
  workarounds, inspect its direct consumers. Flag a miss of Fix duplication
  while the change is open in `./principles.md`.
- Check edge cases, error paths, concurrency, data boundaries, and compatibility.
- Apply Prove the outcome from `./principles.md`. A build, linter, or mocked
  test is insufficient when it does not observe the affected behavior.
- Raise a verification finding when a meaningful changed behavior remains
  unproven and there is a concrete practical proof path. A bug fix that can
  recur through a public seam without a regression test at that seam is a
  finding. Test admission otherwise follows Tests in `./principles.md`. Do not
  request coverage that the Tests principle excludes.
- Treat a removed or changed test as a finding only when it abandons a stable
  contract. Do not preserve a test that the Tests principle says to remove.

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
