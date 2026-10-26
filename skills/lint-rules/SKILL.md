---
name: lint-rules
description: >
  Use when fixing a bug, the user corrects the same mistake twice or says
  "you missed some", a review comment or gotcha recurs, or you are about to
  write a "don't do X"/"always do Y" into AGENTS.md, docs, a skill or a
  comment: the default answer is a repo-local lint rule. Also use when setting
  up or tightening linting, or contributing a rule back. Detects languages,
  wires shared rules into the linter and CI, and turns patterns into custom
  rules with failing examples.
---

# Lint Rules

Stock linters are the baseline. Custom rules catch the patterns that keep
coming back. Every rule ships with a failing example and a catalog entry.

If a linter can check it, it is a lint rule, not prose. Prose (AGENTS.md,
docs, skills, comments) holds only the why and what no linter can check.
Rules break the build; prose gets skipped.

## Rules Explain Themselves

Every message states the problem, why it matters (one line, or a link to the
catalog row or decision note) and the fix, so an agent hitting it can act
without reading anything else.

```yaml
message: >-
  Do not branch on error text. Messages change and wrapping breaks the match.
  Compare with errors.Is against a sentinel, or errors.As into a typed error.
```

```toml
{ path = "std::thread::sleep", reason = "fixed sleeps waste time when fast and flake when slow; wait on an event or a bounded deadline" },
```

## Apply The Rules To A Repo

1. Detect languages: `go.mod` (Go), `Cargo.toml` (Rust), Kubernetes YAML
   (`kind:` + `apiVersion:`, kustomize or Helm output), CSS (`*.css` or a
   Stylelint config).
2. Go: enable golangci-lint with the settings in
   [rules/go/golangci.yml](rules/go/golangci.yml). Copy the `*.yml` files in
   [rules/go/semgrep/](rules/go/semgrep/) to `lint/semgrep/` and run
   `semgrep --config lint/semgrep --exclude vendor`.
3. Rust: merge [rules/rust/Cargo-lints.toml](rules/rust/Cargo-lints.toml) into
   `[workspace.lints]` and opt in library crates only, as its header says, and
   merge [rules/rust/clippy.toml](rules/rust/clippy.toml) into `clippy.toml`.
   Run `cargo clippy --all-targets -- -D warnings`. On an existing crate, gate
   new code first: the cast lints alone hit hundreds of existing `as` casts.
   Land the config with the current hits allowed at crate level, block new
   ones in CI, then remove the allows as you burn them down.
4. Kubernetes: run kubeconform as in [rules/k8s/README.md](rules/k8s/README.md)
   and copy the `*.yml` files in [rules/k8s/semgrep/](rules/k8s/semgrep/).
   Run them on rendered output (`kustomize build`), so placeholder image names
   are already replaced.
5. CSS: copy [rules/css/](rules/css/) (or merge
   [rules/css/create-config.cjs](rules/css/create-config.cjs)). Set `themeFile`,
   `prefix`, and `deprecated` for your palette. Install `stylelint@16.23.1`
   and `postcss-html`, and run Stylelint as in
   [rules/css/README.md](rules/css/README.md).
6. Wire each linter into CI so a failure blocks merge. On an existing codebase,
   gate only new code first (golangci-lint `new-from-rev`, semgrep
   `--baseline-commit`), then burn down the backlog.
7. Run [rules/test.sh](rules/test.sh) here (or `semgrep --test` on the repo's
   rule directory) so every rule's failing example still fails.

## Write A New Rule

1. A problem happens: a bug, a gotcha, a review comment, or a "don't do X"
   you were about to write down. One occurrence is enough for a repo-local
   rule when a linter can catch it.
2. Pick the cheapest tool that can express it: a linter setting, then semgrep
   or ruleguard for Go, Clippy `disallowed_methods`/`disallowed_types` or a
   dylint lint for Rust, Stylelint for CSS, semgrep for YAML.
3. Write the failing example first (`// ruleid: <id>` lines for semgrep, a
   fixture for others) and one passing example (`// ok: <id>`). One rule per
   file: `<id>.yml` plus `<id>.go`, `<id>.test.yaml`, or `<id>.gomod`. Clippy
   rules get a line in [rules/rust/fixture](rules/rust/fixture/src/lib.rs)
   and in the `want` list in `rules/test.sh`.
4. Measure: run it over the repo (every repo you can reach for a global
   rule), count hits and read a sample. Keep it only if most hits are real.
5. Add the rule, its failing example and a catalog row (rule, why, bad
   example, date, origin, hits) in the same PR as the fix. If the repo has no
   `lint/` yet, scaffold the minimum: `lint/semgrep/` plus one CI step running
   `semgrep --test lint/semgrep && semgrep --error --config lint/semgrep`, or
   a `[lints.clippy]` table plus `clippy.toml`.
6. Wire it into CI so a violation fails the build.
7. Delete any AGENTS.md, doc, skill or comment text the rule now enforces.
8. When you later hit a rule, read its message and fix the code. Propose
   loosening the rule only if it is wrong, with the hits that show it.
9. When the same rule recurs across repos, promote it here (see below).

## Global Versus Repo-Local

- **Global** (lives here): true in any codebase, names no product, employer,
  internal package or domain type. Example: no bare `TODO`.
- **Repo-local** (lives in the repo's `lint/` and `docs/lint-rules.md`):
  mentions the repo's own types, packages, APIs or business rules. Example:
  "use `billing.Cents`, not `int64`, for money".

## Contribute A Rule Back

When a repo-local rule turns out to be general, open a PR to agent-skills that
adds the rule file, its failing example and a [CATALOG.md](CATALOG.md) row.
Strip work or employer-specific code, names, paths and identifiers. Rewrite the
example from scratch rather than copying work code. Set origin to a generic
description, never the employer or internal repo name.
