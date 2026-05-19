---
name: calver-release
description: >
  Use when planning, naming, tagging, or documenting releases with Calendar
  Versioning (CalVer), especially SemVer-compatible date-based Git tags such as
  vYYYY.M.D.
  Also use when the user asks for a Git version tag, release tag, version tag,
  tagging the current commit, or publishing a tagged release. Helps choose the
  next release tag, handle multiple releases on the same day, and keep release
  tags and package versions SemVer-compatible. Coordinates with changelog or
  release-note conventions before tagging.
---

# CalVer Release

Use this skill when the user wants date-based release tags instead of semantic
versions, or asks to prepare, tag, publish, or create a release.

## Version Format

- Use Git tags in this SemVer-compatible CalVer format: `vYYYY.M.D`.
- Do not zero-pad month or day. SemVer numeric identifiers cannot contain
  leading zeroes: use `v2026.5.20`, not `v2026.05.20`.
- Use the local calendar date relevant to the release process unless the user or
  project specifies UTC.
- Treat the base date tag as the first release of that day.
- For additional releases on the same day, keep the tag SemVer-compatible. Do
  not append a fourth numeric component.

Examples:

```text
v2026.5.20       # first release on 2026-05-20
v2026.5.20+2     # second same-day release when build metadata is acceptable
v2026.5.20+3     # third same-day release when build metadata is acceptable
```

Build metadata is SemVer-valid but does not affect SemVer precedence. If a
project needs multiple same-day releases with strict SemVer ordering, stop and
ask the user to choose a different scheme before tagging.

## Choosing The Next Tag

1. Get the release date from the local system unless the user gives a specific
   date or the project requires UTC:

```sh
date +%Y.%-m.%-d
```

2. Build the tag prefix: `vYYYY.M.D`.
3. Inspect existing tags with that prefix:

```sh
git tag --list "vYYYY.M.D*"
```

4. Match only exact date tags and same-day SemVer build metadata suffixes:
   - `vYYYY.M.D`
   - `vYYYY.M.D+N`
5. For the release date:
   - if no matching tag exists, use `vYYYY.M.D`
   - if `vYYYY.M.D` exists and no suffix exists, use `vYYYY.M.D+2`
   - if build metadata suffixes exist, use the next integer suffix
6. Do not use `+1`; the unsuffixed date tag is release one.
7. If the user gives a specific date, use that date instead of today's date.

## Git Workflow

- When the user asks to do a release, tag the latest local commit (`HEAD`).
- Do not proceed if there are uncommitted, staged, or untracked files. The
  release must point at a committed state.
- Inspect repository state before tagging:

```sh
git status --short
```

- If `git status --short` prints anything, stop and ask the user to commit,
  discard, or intentionally ignore those files before releasing.
- Inspect the commit being released:

```sh
git log -1 --oneline
```

- Check for existing changelog or release-note conventions before tagging:

```sh
find . -maxdepth 3 \( -iname 'CHANGELOG*' -o -path './docs/releases' -o -path './docs/changelog' -o -path './releases' -o -path './.github/releases' \)
```

- If changelog/release-note evidence exists, require the `changelog` workflow
  before tagging unless the user explicitly says to skip changelog/release-note
  updates for this release. In-repo changelog/release-note updates must be
  committed before tagging.
- If no changelog convention exists, do not invent one during release unless the
  user asks. Draft release notes text only when requested.
- Create an annotated tag:

```sh
git tag -a vYYYY.M.D -m "Release vYYYY.M.D"
```

- Never overwrite, move, or delete an existing tag unless the user explicitly
  asks and understands the impact.
- Push the tag only when the user asks to publish, push, or complete the
  release:

```sh
git push origin vYYYY.M.D
```

- If the user only asks to choose or prepare a release tag, do not create or
  push anything.

## Release Command Checklist

When asked to perform a release:

1. Run `git status --short`.
2. Stop if any tracked or untracked changes are present.
3. Run `date +%Y.%-m.%-d` unless a date was provided.
4. Run `git tag --list "vYYYY.M.D*"`.
5. Choose the next tag using the suffix rules.
6. Check for changelog/release-note evidence.
7. If evidence exists, ensure changelog/release notes are updated or explicitly
   skipped by the user before tagging. If files were updated, ensure those
   updates are committed.
8. Run `git log -1 --oneline` and report the commit being tagged.
9. Create the annotated tag on `HEAD`.
10. Push the tag if the user asked to publish the release.
11. Report the tag and commit hash.

## Package Version Notes

- `vYYYY.M.D` is good for Git tags and GitHub releases, and the version portion
  is SemVer-compatible.
- If package metadata also stores a version, use the same SemVer-compatible
  CalVer value without the leading `v`.
- For Zig packages, update `build.zig.zon` `version` with dot-separated numeric
  components without leading zeroes:

```zig
.version = "2026.5.20",
```

- Do not use a fourth numeric component such as `2026.5.20.2`; it is not
  SemVer.
- For same-day package releases, use SemVer build metadata only if the package
  tooling accepts it and the lack of precedence difference is acceptable:

```text
tag:      v2026.5.20+2
version:  2026.5.20+2
```

- If package tooling rejects build metadata or requires strictly increasing
  SemVer precedence for same-day releases, stop and ask the user to choose a
  different version scheme before releasing.

## Release Notes

- Use the `changelog` skill when release notes or changelog updates are needed.
- Use the tag as the release title unless the project has a stronger convention.
- Group notable changes by user impact:
  - features
  - fixes
  - breaking changes
  - documentation
  - maintenance
- Mention the previous tag used as the comparison base when summarizing changes.
