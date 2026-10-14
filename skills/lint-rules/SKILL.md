---
name: lint-rules
description: >
  Use when setting up or tightening linting in any repo, when a review comment
  or bug pattern recurs, or when contributing a lint rule back. Detects the
  repo's languages, wires the shared rule files into the linter config and CI,
  and turns recurring patterns into custom rules with failing examples.
---

# Lint Rules

Stock linters are the baseline. Custom rules catch the patterns that keep
coming back. Every rule ships with a failing example and a catalog entry.

## Apply The Rules To A Repo

1. Detect languages: `go.mod` (Go), `Cargo.toml` (Rust), Kubernetes YAML
   (`kind:` + `apiVersion:`, kustomize or Helm output).
2. Go: enable golangci-lint with the settings in
   [rules/go/golangci.yml](rules/go/golangci.yml). Copy
   [rules/go/semgrep.yml](rules/go/semgrep.yml) to `lint/semgrep/go.yml`.
3. Rust: merge [rules/rust/Cargo-lints.toml](rules/rust/Cargo-lints.toml) into
   `Cargo.toml` (or `[workspace.lints]`) and
   [rules/rust/clippy.toml](rules/rust/clippy.toml) into `clippy.toml`. Run
   `cargo clippy --all-targets -- -D warnings`.
4. Kubernetes: run kubeconform as in [rules/k8s/README.md](rules/k8s/README.md)
   and copy [rules/k8s/semgrep.yml](rules/k8s/semgrep.yml).
5. Wire each linter into CI so a failure blocks merge. On an existing codebase,
   gate only new code first (golangci-lint `new-from-rev`, semgrep
   `--baseline-commit`), then burn down the backlog.
6. Run `semgrep --test` on the rule directory so every rule's failing example
   still fails.

## Write A New Rule

1. Spot the pattern: the same review comment twice, a bug class that shipped
   twice, or a comment that says "don't do X".
2. Pick the cheapest tool that can express it: a linter setting, then semgrep
   or ruleguard for Go, Clippy `disallowed_methods`/`disallowed_types` or a
   dylint lint for Rust, semgrep for YAML.
3. Write the failing example first (`// ruleid: <id>` lines for semgrep, a
   fixture for others) and one passing example (`// ok: <id>`).
4. Add the rule and a row in the repo's rules catalog: rule, why, bad example,
   date added, origin.
5. Delete any comment or skill text the rule now enforces.

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
