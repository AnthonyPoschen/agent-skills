---
name: changelog
description: >
  Generate, update, or review changelogs and release notes. Use when the user
  asks for a changelog, release notes, GitHub release text, changes since a tag,
  what changed between versions, upgrade notes, release summaries, or release PR
  descriptions. Converts commits and diffs into user-facing change entries and
  follows existing repo release-note conventions.
---

# Changelog

Use this skill to turn Git history and diffs into user-facing release notes.
Write for the person deciding whether to upgrade, install, review, or publish a
release.

## Activation

Use this skill for:

- changelogs
- release notes
- GitHub release text
- changes since a tag
- changes between versions
- upgrade notes
- release summaries
- release PR descriptions

## Convention Detection

Before writing files, inspect existing conventions:

- `CHANGELOG.md`
- `CHANGELOG*`
- `docs/releases/`
- `docs/changelog/`
- `releases/`
- `.github/releases/`
- previous GitHub release text when available

Follow the existing convention unless the user asks to change it. If no
convention exists, prefer producing release notes text instead of creating files
automatically. Ask before introducing a new changelog storage convention.

## Commit State

- Generate release changelogs and release notes from committed changes by
  default.
- Before generating release-ready notes, require a clean worktree:

```sh
git status --short
```

- If `git status --short` prints anything, stop and ask the user to commit,
  discard, or intentionally ignore those changes before generating release
  notes.
- Do not include staged, unstaged, or untracked work in release notes unless the
  user explicitly asks for draft notes.
- If updating in-repo changelog files, make that update as a separate
  release-prep commit before tagging.
- Tag the commit that includes any required changelog/release-note file updates.
- If producing GitHub Release text only, no in-repo changelog commit is
  required.

## Preferred Release Notes Model

When introducing or recommending a convention, prefer each location for a
different purpose:

- Git tag: canonical version pointer and immutable release anchor.
- GitHub Release: primary public release page for polished notes, downloads,
  links, and upgrade guidance.
- `docs/releases/<tag>.md`: durable in-repo release record, reviewable in PRs
  and preserved with source history.
- `CHANGELOG.md`: concise human index with recent summaries and links to
  detailed release files. Avoid making it a forever-growing dump.

Example index:

```markdown
# Changelog

## Recent Releases

- `v2026.5.20` - Added frontend design skill and SemVer-compatible CalVer
  release workflow. See `docs/releases/v2026.5.20.md`.
```

## Range Selection

- If the user gives a tag, version, branch, or range, use it.
- Otherwise, compare the previous release tag to `HEAD`.
- Inspect tags before choosing a range:

```sh
git tag --sort=-creatordate
```

- Inspect commits and changed files:

```sh
git log --oneline <previous-tag>..HEAD
git diff --stat <previous-tag>..HEAD
```

- Read detailed diffs only for changes that need clarification.

## Writing Rules

- Group by user impact, not Git implementation detail.
- Prefer these headings, omitting empty sections:
  - Breaking
  - Added
  - Changed
  - Fixed
  - Deprecated
  - Removed
  - Security
  - Notes
- Put breaking changes first and include migration guidance.
- Write entries in plain language.
- Mention package names, skill names, command names, config names, or APIs when
  they help users understand impact.
- Exclude noisy internal-only commits unless they affect users, maintainers,
  release process, or upgrade behavior.
- Include validation notes only when they are useful to release reviewers.

Bad:

```markdown
- Refactored release skill internals.
```

Good:

```markdown
- Release tags now use SemVer-compatible CalVer format: `vYYYY.M.D`.
```

## Output Modes

- "What changed?" means summarize in the response only.
- "Release notes" means draft publishable GitHub Release text.
- "Update changelog" means edit the repo's existing changelog convention.
- "Prepare release" means produce release notes before tagging or publishing.

## File Update Rules

- If `docs/releases/` exists, create or update one release file per tag.
- If `CHANGELOG.md` is a concise index, add a short entry linking to the release
  file.
- If `CHANGELOG.md` is the existing full changelog convention, update it in the
  same style.
- If no changelog convention exists, do not create files unless the user asks or
  confirms a proposed convention.
- Do not duplicate full release notes in multiple places unless the repo already
  does that.

## Review Checklist

- The range is correct.
- The notes are user-facing, not just commit messages.
- Breaking changes and migration notes are obvious.
- Internal-only noise is removed.
- The output follows existing repo convention.
- Links, tag names, dates, and version names are consistent.
