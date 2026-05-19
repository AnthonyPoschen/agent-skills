---
name: calver-release
description: >
  Use when planning, naming, tagging, or documenting releases with Calendar
  Versioning (CalVer), especially date-based Git tags such as vYYYY.MM.DD.
  Also use when the user asks for a Git version tag, release tag, version tag,
  tagging the current commit, or publishing a tagged release. Helps choose the
  next release tag, handle multiple releases on the same day, and avoid
  accidental SemVer assumptions.
---

# CalVer Release

Use this skill when the user wants date-based release tags instead of semantic
versions, or asks to prepare, tag, publish, or create a release.

## Version Format

- Use Git tags in this format: `vYYYY.MM.DD`.
- Always zero-pad month and day: `v2026.05.20`, not `v2026.5.20`.
- Use the local calendar date relevant to the release process unless the user or
  project specifies UTC.
- Treat the base date tag as the first release of that day.
- For additional releases on the same day, append a numeric suffix starting at
  `.2`.

Examples:

```text
v2026.05.20      # first release on 2026-05-20
v2026.05.20.2    # second release on 2026-05-20
v2026.05.20.3    # third release on 2026-05-20
```

## Choosing The Next Tag

1. Get the release date from the local system unless the user gives a specific
   date or the project requires UTC:

```sh
date +%Y.%m.%d
```

2. Build the tag prefix: `vYYYY.MM.DD`.
3. Inspect existing tags with that prefix:

```sh
git tag --list "vYYYY.MM.DD*"
```

4. Match only exact date tags and same-day numeric suffixes:
   - `vYYYY.MM.DD`
   - `vYYYY.MM.DD.N`
5. For the release date:
   - if no matching tag exists, use `vYYYY.MM.DD`
   - if `vYYYY.MM.DD` exists and no suffix exists, use `vYYYY.MM.DD.2`
   - if suffixes exist, use the next integer suffix
6. Do not use `.1`; the unsuffixed date tag is release one.
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

- Create an annotated tag:

```sh
git tag -a vYYYY.MM.DD -m "Release vYYYY.MM.DD"
```

- Never overwrite, move, or delete an existing tag unless the user explicitly
  asks and understands the impact.
- Push the tag only when the user asks to publish, push, or complete the
  release:

```sh
git push origin vYYYY.MM.DD
```

- If the user only asks to choose or prepare a release tag, do not create or
  push anything.

## Release Command Checklist

When asked to perform a release:

1. Run `git status --short`.
2. Stop if any tracked or untracked changes are present.
3. Run `date +%Y.%m.%d` unless a date was provided.
4. Run `git tag --list "vYYYY.MM.DD*"`.
5. Choose the next tag using the suffix rules.
6. Run `git log -1 --oneline` and report the commit being tagged.
7. Create the annotated tag on `HEAD`.
8. Push the tag if the user asked to publish the release.
9. Report the tag and commit hash.

## Package Version Notes

- `vYYYY.MM.DD` is good for Git tags and GitHub releases.
- If a package manager requires SemVer, do not put `YYYY-MM-DD` or
  `YYYY.MM.DD` into the package version unless that tool explicitly supports it.
- For SemVer-only package metadata, keep a compatible package version and use
  the CalVer tag/release name separately.
- For Zig packages, Git tags can use the zero-padded release tag
  `vYYYY.MM.DD`, but `build.zig.zon` `version` must be SemVer-compatible. Use
  dot-separated numeric components without leading zeroes:

```zig
.version = "2026.5.20",
```

- Do not mirror a same-day tag suffix as a fourth numeric component in
  `build.zig.zon`; `2026.5.20.2` is not SemVer.
- For a second same-day Zig package release, keep the Git tag suffix and choose
  one of these package-version approaches deliberately:
  - leave `build.zig.zon` at `2026.5.20` when the package version does not need
    to distinguish same-day tag rebuilds
  - use SemVer build metadata such as `2026.5.20+2` when the tooling accepts it
    and version precedence is not important for that suffix

```text
tag:                 v2026.05.20.2
version option 1:    2026.5.20
version option 2:    2026.5.20+2
```

- When updating package metadata, keep the Git tag and package version derived
  from the same date/suffix where the package manager supports it. If package
  metadata cannot represent the suffix safely, document the tag as the canonical
  release identifier.

## Release Notes

- Use the tag as the release title unless the project has a stronger convention.
- Group notable changes by user impact:
  - features
  - fixes
  - breaking changes
  - documentation
  - maintenance
- Mention the previous tag used as the comparison base when summarizing changes.
