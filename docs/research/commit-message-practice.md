# Commit Message Practice

## Question

Should `git-commit-workflow` require Conventional Commit prefixes when this
repository already prepares useful changelogs?

## Evidence

### Omarchy

The [latest 100 Omarchy commits](https://github.com/basecamp/omarchy/commits/master/)
sampled on 2026-08-25 had no subjects matching the Conventional Commit form
`type(scope): summary`. David Heinemeier Hansson's subjects state the change
in ordinary English, normally with an imperative verb:

- [Regenerate mise wrappers that still print mise's output to stdout](https://github.com/basecamp/omarchy/commit/535d8f3485581b35ab0d368d4179b362eb6dea32)
- [Fall back to polkit when the DNS sudoers grant is missing](https://github.com/basecamp/omarchy/commit/1e70cca144eb54a998217cb61fe92e1a6ca51ed2)
- [Decode mixed UTF-16 clipboard text](https://github.com/basecamp/omarchy/commit/238021cd679b9ec213e4117721288628279ff930)

The first commit uses its body for information the diff alone does not make
clear: why rewriting existing wrappers is needed, the safety condition for
recognising an existing wrapper, and that a second run is a no-op. Many smaller
commits have no explanatory body. This supports a short, concrete subject by
default, with a body only for consequential reasoning, risk, migration, or
validation.

Some Omarchy subjects are intentionally terse, such as `Words` and `Better
without gaps`. They show that the project does not impose a format, but they
are not a useful default for a reusable workflow.

### Rails

Rails offers the same pattern in DHH-authored commits. The subjects are direct
descriptions such as [Add `ActionDispatch::Request#bearer_token` to extract the
bearer token from the Authorization header](https://github.com/rails/rails/commit/2533c938acb97c2b44e6600fc0d35962e7ba9c7d)
and [Allow Rails.app.creds to access .env values in dev](https://github.com/rails/rails/commit/80570a6dd73d632738ed86561b10d33217c6fa7e).
They do not use Conventional Commit prefixes. The latter commit records an
explicit changelog update in its development history, separating release-note
work from the subject's job of describing the change.

## Implications for `git-commit-workflow`

- Replace the required Conventional Commit form with a plain-English subject
  that states the concrete behaviour or structural change.
- Retain imperative mood, concise wording, logical commit boundaries, and the
  AI assistance trailer.
- Add a body only when the reason, tradeoff, risk, migration, or validation is
  not apparent from the diff. Do not add a body merely to restate the subject.
- Treat changelogs as their own reader-facing release artifact. Do not force a
  commit type solely to make automated changelog grouping possible.
- Respect an existing repository convention when it explicitly requires a
  different format. In the absence of one, do not invent a type or scope.

Suggested subjects:

```text
Add file organization guidance
Keep generated wrappers quiet
Split notification parsing from delivery
```
